package cursor

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// The Cursor struct
// This struct has no field, can be used for its functions
type Cursor struct{}

// Receives a list of Ansi codes and print
// them to the terminal, effectively apply
// the ANSI codes to the terminal
func (c Cursor) Print(codes ...Ansi) {
	b := bytes.Buffer{}
	for _, v := range codes {
		b.Write([]byte(v))
	}
	fmt.Print(b.String())
}

// Moves the cursor up, as given in `lineUp`, clear
// it and restore the cursor to its original position
func (c Cursor) ClearPreviousLine(lineUp int) {
	c.Print(
		AnsiSaveCursor,
		c.ANSIToMoveUp(lineUp),
		AnsiClearLine,
		AnsiRestoreCursor,
	)
}

// Moves the cursor up, as given in `lineUp`, clear
// it and replace it with the given replacement
func (c Cursor) ReplaceLine(lineUp int, replacement string) {
	c.Print(
		AnsiSaveCursor,
		c.ANSIToMoveUp(lineUp),
		AnsiClearLine,
		AnsiCursorToColumn0,
		Ansi(replacement),
		AnsiRestoreCursor,
	)
}

// returns the ANSI code to move the cursor up n times
func (c Cursor) ANSIToMoveUp(n int) Ansi {
	return Ansi(strings.Replace(AnsiCursorUp, "1", strconv.Itoa(n), 1))
}

// returns the ANSI code to move the cursor down n times
func (c Cursor) ANSIToMoveDown(n int) Ansi {
	return Ansi(strings.Replace(AnsiCursorDown, "1", strconv.Itoa(n), 1))
}

// returns the ANSI code to move the cursor right n times
func (c Cursor) MoveRight(n int) Ansi {
	return Ansi(strings.Replace(AnsiCursorRight, "1", strconv.Itoa(n), 1))
}

// returns the ANSI code to move the cursor left n times
func (c Cursor) MoveLeft(n int) Ansi {
	return Ansi(strings.Replace(AnsiCursorLeft, "1", strconv.Itoa(n), 1))
}
