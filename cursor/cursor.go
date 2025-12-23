package cursor

import (
	"fmt"

	"github.com/vieolo/termange/internal"
)

// The Cursor struct
// This struct has no field, can be used for its functions
type Cursor struct{}

// Moves the cursor up n times
func (c Cursor) MoveUp(n int) Cursor {
	for range n {
		fmt.Print(internal.CursorUp)
	}
	return c
}

// Moves the cursor down n times
func (c Cursor) MoveDown(n int) Cursor {
	for range n {
		fmt.Print(internal.CursorDown)
	}
	return c
}

// Moves the cursor right n times
func (c Cursor) MoveRight(n int) Cursor {
	for range n {
		fmt.Print(internal.CursorRight)
	}
	return c
}

// Moves the cursor left n times
func (c Cursor) MoveLeft(n int) Cursor {
	for range n {
		fmt.Print(internal.CursorLeft)
	}
	return c
}

// Moves the cursor to the start of the line
func (c Cursor) MoveToLineStart() Cursor {
	fmt.Print(internal.CursorToColumn0)
	return c
}

// Clears the current line
func (c Cursor) ClearLine() Cursor {
	fmt.Print(internal.ClearLine)
	return c
}
