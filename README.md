# Termange

[![Go Reference](https://pkg.go.dev/badge/github.com/vieolo/termange.svg)](https://pkg.go.dev/github.com/vieolo/termange)

Utility functions and tools for terminal interaction

## TUI Components

- `Select` -> Creates a list and the user can select an item by navigating with the arrow keys.
- `TextInput` -> Creates a prompt asking for user input
- `Confirm` -> Asks a yes or no question from the user


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