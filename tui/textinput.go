package tui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/x/term"
	"github.com/vieolo/termange/cursor"
)

type TextInputOptions struct {
	// The question/info to the user on what they should enter
	Prompt string
	// If this function provided, the prompt will be deleted and the
	// return string will be printed instead
	PostValuePrint func(string) string
	// Where the answer is read from. When nil, it is os.Stdin
	In io.Reader
	// Where the prompt is written. When nil, it is os.Stdout
	Out io.Writer
}

// Creates an text input and returns the input.
//
// The input is read up to the end of the line, and nothing after it,
// so the next prompt, or any other code reading the same input, gets the
// rest. The line ending, `\n` or `\r\n`, is removed. When the input ends,
// such as a closed stdin, it returns what was entered before the end,
// which is an empty string when nothing was
func TextInput(options TextInputOptions) string {
	in, out := streams(options.In, options.Out)
	finalValue, err := ask(options.Prompt, in, out)
	if err != nil {
		return ""
	}

	if options.PostValuePrint != nil {
		cursor.Cursor{W: out}.ReplaceLine(1, options.PostValuePrint(finalValue))
	}

	return finalValue
}

// Returns the given input and output, or os.Stdin and os.Stdout in their place
func streams(in io.Reader, out io.Writer) (io.Reader, io.Writer) {
	if in == nil {
		in = os.Stdin
	}
	if out == nil {
		out = os.Stdout
	}
	return in, out
}

// Writes the prompt to out and reads the answer from in.
//
// The cursor ends on the line below the prompt, so that the prompt is
// always the previous line. A terminal moves it there itself when it
// echoes the user's Enter. Otherwise, such as for a pipe or when the
// input ends with Ctrl-D, a new line is written to out instead
func ask(prompt string, in io.Reader, out io.Writer) (string, error) {
	if !strings.HasSuffix(prompt, " ") {
		prompt += " "
	}
	fmt.Fprint(out, prompt)
	line, err := readLine(in)
	if !strings.HasSuffix(line, "\n") || !isTerminal(in) {
		fmt.Fprintln(out)
	}
	return strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r"), err
}

// Reads a line from r, with its `\n` when it has one.
//
// It reads one byte at a time, so it never consumes anything after the line.
// A last line that ends at the end of the input, without a `\n`, is
// returned with a nil error. io.EOF is returned only when nothing was read
func readLine(r io.Reader) (string, error) {
	var line []byte
	b := make([]byte, 1)
	for {
		n, err := r.Read(b)
		if n > 0 {
			line = append(line, b[0])
			if b[0] == '\n' {
				break
			}
		}
		if err != nil {
			if err == io.EOF && len(line) > 0 {
				break
			}
			return "", err
		}
	}
	return string(line), nil
}

func isTerminal(r io.Reader) bool {
	f, ok := r.(interface{ Fd() uintptr })
	return ok && term.IsTerminal(f.Fd())
}
