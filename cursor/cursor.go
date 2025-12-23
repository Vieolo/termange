package cursor

import (
	"fmt"

	"github.com/vieolo/termange/internal"
)

// Internal terminal
type Cursor struct{}

func (c Cursor) MoveUp(line int) Cursor {
	for range line {
		fmt.Print(internal.CursorUp)
	}
	return c
}

func (c Cursor) ClearLine() Cursor {
	fmt.Print(internal.ClearLine)
	return c
}
