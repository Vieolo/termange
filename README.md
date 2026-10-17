# Termange

[![Go Reference](https://pkg.go.dev/badge/github.com/vieolo/termange.svg)](https://pkg.go.dev/github.com/vieolo/termange)

Utility functions and tools for terminal interaction

## TUI Components

- `Select` -> Creates a list and the user can select an item by navigating with the arrow keys.
- `MultiSelect` -> Creates a list with checkboxes, where the user selects or unselects items with space and confirms with enter. `Selected` sets the items that start selected
- `TextInput` -> Creates a prompt asking for user input
- `Confirm` -> Asks a yes or no question from the user. `Default` is the answer when the user just presses Enter

```go
func main() {
    write := tui.Confirm(tui.ConfirmOptions{
        Prompt:  "Write sirup.yaml?", // Shown as "Write sirup.yaml? [Y/n]"
        Default: true,
    })
}
```


## Logging
Each function has two variations: `*ln` which print a new line and `*f` which would similar to `fmt.Printf`. Example:

```go
func main() {
    termange.PrintSuccessln("Success message")
    termange.PrintSuccessf("You entered: %s", "name")
}
```

- `PrintSuccess` -> Prints the text in a new line in green
- `PrintError` -> Prints the text in a new line in red
- `PrintWarning` -> Prints the text in a new line in yellow
- `PrintInfo` -> Prints the text in a new line in bold white

`PaintText` colors any text. For other colors and styles, use the codes of [mansil](https://github.com/vieolo/mansil).


## Live Region
The `live` package draws a block of lines at the end of the output and redraws it in place, for progress, spinners and status trees.

```go
func main() {
    region := live.New(os.Stdout, func(width, height int, final bool) []string {
        if final {
            return []string{"✓ done"}
        }
        return []string{live.Spinner(time.Now()) + " working"}
    }, live.Options{Interval: 100 * time.Millisecond})

    region.Start()
    // ... do the work, and call region.Refresh() whenever the state changes
    region.Stop() // Leaves the final frame on screen
}
```

Lines are cut to the terminal's width, and a frame taller than the terminal ends with `… N more`. Nothing else should be written to the terminal between `Start` and `Stop`.


## Cursor
`cursor.Cursor` moves the cursor and edits the lines of the terminal. It writes to stdout, or to `W` when set.

```go
func main() {
    cursor.Cursor{W: os.Stderr}.ReplaceLine(1, "Replaced")
}
```


## Examples
The `example` folder runs each part of termange. Pick one from a list with:

```bash
go run ./example
```

Or run one by its name: `logging`, `input`, `confirm`, `multiselect`, `table`, `live`, `cursor` or `command`:

```bash
go run ./example live
```


## Development
Run the checks (format, build, vet and tests with the race detector) with:

```bash
go run ./internal/check
```
