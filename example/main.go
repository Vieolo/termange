// A tour of termange. Run it from the root of the repo with
//
//	go run ./example
//
// to pick the examples from a list, or with the name of one, such as
//
//	go run ./example live
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/vieolo/mansil"
	"github.com/vieolo/termange"
	"github.com/vieolo/termange/cursor"
	"github.com/vieolo/termange/live"
	"github.com/vieolo/termange/tui"
)

type example struct {
	name string
	run  func()
}

var examples = []example{
	{"logging", logging},
	{"input", input},
	{"confirm", confirm},
	{"table", tableExample},
	{"live", liveRegion},
	{"cursor", cursorExample},
	{"command", command},
}

func main() {
	if len(os.Args) > 1 {
		for _, e := range examples {
			if e.name == os.Args[1] {
				e.run()
				return
			}
		}
		termange.PrintErrorf("Unknown example %q. The examples are: %s\n", os.Args[1], strings.Join(names(), ", "))
		os.Exit(1)
	}

	// Select picks the example, until the user quits with q, esc or ctrl+c
	var items []tui.SelectItem
	for _, n := range names() {
		items = append(items, tui.SelectItem(n))
	}
	for {
		choice, err := tui.Select(tui.SelectOptions{
			Title:         "termange examples",
			Items:         items,
			SelectMessage: func(s string) string { return termange.PaintText("> "+s, termange.ColorCyan) },
			QuitMessage:   "Bye",
		})
		if err != nil {
			termange.PrintErrorln(err.Error())
			os.Exit(1)
		}
		if choice == "" {
			return
		}
		for _, e := range examples {
			if e.name == choice {
				e.run()
			}
		}
		fmt.Println()
	}
}

func names() []string {
	var n []string
	for _, e := range examples {
		n = append(n, e.name)
	}
	return n
}

// The Print functions, and PaintText for any text
func logging() {
	termange.PrintSuccessln("PrintSuccessln: green")
	termange.PrintErrorln("PrintErrorln: red, to stderr")
	termange.PrintWarningln("PrintWarningln: yellow")
	termange.PrintInfoln("PrintInfoln: bold white")
	termange.PrintSuccessf("PrintSuccessf: %d%% like fmt.Printf\n", 100)

	colors := []termange.Color{
		termange.ColorBlack, termange.ColorRed, termange.ColorGreen, termange.ColorYellow,
		termange.ColorBlue, termange.ColorMagenta, termange.ColorCyan, termange.ColorWhite,
	}
	var painted []string
	for _, c := range colors {
		painted = append(painted, termange.PaintText(string(c), c))
	}
	fmt.Println("PaintText:", strings.Join(painted, " "))

	// The other colors and styles come from mansil
	fmt.Println("mansil:", mansil.Underline+mansil.BlueFGBright+"underlined bright blue"+mansil.Reset)
}

// TextInput, replacing the prompt with the answer
func input() {
	name := tui.TextInput(tui.TextInputOptions{
		Prompt: "What is your name?",
		PostValuePrint: func(v string) string {
			return "Name: " + termange.PaintText(v, termange.ColorCyan)
		},
	})
	termange.PrintInfof("Hello, %s\n", name)
}

// Confirm, where Enter gives the default
func confirm() {
	like := tui.Confirm(tui.ConfirmOptions{
		Prompt:  "Do you like the examples?",
		Default: true,
		PostConfirmationPrint: func(v bool) string {
			if v {
				return termange.PaintText("Glad to hear it", termange.ColorGreen)
			}
			return termange.PaintText("Noted", termange.ColorYellow)
		},
	})

	// The default is false, so the prompt shows [y/N]
	deleteAll := tui.Confirm(tui.ConfirmOptions{Prompt: "Delete everything?"})

	fmt.Printf("You answered %v and %v\n", like, deleteAll)
}

// Table, printing the row picked with Enter, until q
func tableExample() {
	fmt.Println("Pick a row with Enter, and quit with q")
	termange.Table(termange.TableData{
		Columns: []table.Column{
			{Title: "Package", Width: 10},
			{Title: "What it does", Width: 44},
		},
		Rows: []table.Row{
			{"termange", "Printing in colour and running commands"},
			{"tui", "Select, TextInput and Confirm"},
			{"live", "Lines redrawn in place, for progress"},
			{"cursor", "Moving the cursor and editing lines"},
		},
		OnRowClick: func(r table.Row) string { return "You picked " + r[0] },
	})
}

// A live region showing a tree of tasks, the way a test runner would
func liveRegion() {
	type task struct {
		name     string
		duration time.Duration
		fails    bool

		state    string // waiting, running or done
		started  time.Time
		finished time.Time
	}
	tasks := []*task{
		{name: "app", duration: 1500 * time.Millisecond},
		{name: "core", duration: 2500 * time.Millisecond},
		{name: "cli", duration: 1000 * time.Millisecond, fails: true},
	}

	// The tasks are changed under mu, and the region renders under mu
	var mu sync.Mutex
	render := func(width, height int, final bool) []string {
		mu.Lock()
		defer mu.Unlock()

		lines := []string{"example test"}
		for i, t := range tasks {
			branch := "├─"
			if i == len(tasks)-1 {
				branch = "└─"
			}
			var state string
			switch {
			case t.state == "running":
				state = fmt.Sprintf("%s running  %.1fs", live.Spinner(time.Now()), time.Since(t.started).Seconds())
			case t.state == "done" && t.fails:
				state = termange.PaintText("✗ failed", termange.ColorRed) + fmt.Sprintf("  %.1fs", t.finished.Sub(t.started).Seconds())
			case t.state == "done":
				state = termange.PaintText("✓ passed", termange.ColorGreen) + fmt.Sprintf("  %.1fs", t.finished.Sub(t.started).Seconds())
			default:
				state = mansil.Dim + "· waiting" + mansil.Reset
			}
			lines = append(lines, fmt.Sprintf("%s %-5s %s", branch, t.name, state))
		}
		return lines
	}

	region := live.New(os.Stdout, render, live.Options{Interval: 100 * time.Millisecond})

	// mu is released before Refresh, since Refresh calls render, which takes mu
	run := func(t *task) {
		mu.Lock()
		t.state, t.started = "running", time.Now()
		mu.Unlock()
		region.Refresh()

		time.Sleep(t.duration)

		mu.Lock()
		t.state, t.finished = "done", time.Now()
		mu.Unlock()
		region.Refresh()
	}

	region.Start()
	var wg sync.WaitGroup
	wg.Go(func() { run(tasks[0]); run(tasks[2]) }) // cli waits for app
	wg.Go(func() { run(tasks[1]) })
	wg.Wait()
	region.Stop()

	termange.PrintErrorln("cli: 1 test failed")
}

// Replacing and clearing lines that are already printed
func cursorExample() {
	fmt.Println("Line 1")
	fmt.Println("Line 2: will be replaced")
	fmt.Println("Line 3: will be cleared")

	c := cursor.Cursor{}
	time.Sleep(time.Second)
	c.ReplaceLine(2, termange.PaintText("Line 2: replaced", termange.ColorGreen))
	time.Sleep(time.Second)
	c.ClearPreviousLine(1)
}

// RunCommand with a start text that is cleared once it is done, and a failing command
func command() {
	result, err := termange.RunCommand(termange.CommandConfig{
		Command:        "sh",
		Args:           []string{"-c", "sleep 1; go version"},
		StartText:      "Running go version...",
		ClearStartText: true,
	})
	if err != nil {
		termange.PrintErrorln(err.Error())
		return
	}
	termange.PrintSuccessf("%s", result.Stdout.String())
	termange.PrintInfof("Took %dms\n", result.RealTime/1000)

	_, err = termange.RunRawCommand("exit 3")
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		termange.PrintWarningf("The failing command exited with %d\n", exitErr.ExitCode())
	}
}
