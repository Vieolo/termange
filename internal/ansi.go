package internal

import "fmt"

// ANSI escape codes for cursor manipulation
const (
	// CursorUp moves the cursor up N lines. We use 1 to target the immediate previous line.
	cursorUp = "\033[1A"
	// ClearLine erases the entire line containing the cursor (ANSI code 2K).
	clearLine = "\033[2K"
)

// Internal terminal
type IT struct{}

func (i IT) CursorUp() IT {
	fmt.Print(cursorUp)
	return i
}

func (i IT) ClearLine() IT {
	fmt.Print(clearLine)
	return i
}
