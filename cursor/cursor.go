package cursor

import (
	"io"
	"os"
	"strings"

	"github.com/vieolo/mansil"
)

// The Cursor struct
// Its functions write the ANSI codes to W, effectively
// moving the cursor and editing the lines of the terminal
type Cursor struct {
	// Where the codes are written. When nil, they are
	// written to os.Stdout
	W io.Writer
}

// Receives a list of Ansi codes and writes them
// to W in a single write, effectively apply
// the ANSI codes to the terminal
func (c Cursor) Print(codes ...Ansi) {
	var b strings.Builder
	for _, v := range codes {
		b.WriteString(string(v))
	}
	io.WriteString(c.writer(), b.String())
}

// Moves the cursor up, as given in `lineUp`, clear
// it and restore the cursor to its original position.
// A `lineUp` of 0 clears the current line
func (c Cursor) ClearPreviousLine(lineUp int) {
	c.Print(
		AnsiSaveCursor,
		c.ANSIToMoveUp(lineUp),
		mansil.ClearLine,
		AnsiRestoreCursor,
	)
}

// Moves the cursor up, as given in `lineUp`, clear
// it and replace it with the given replacement.
// A `lineUp` of 0 replaces the current line
func (c Cursor) ReplaceLine(lineUp int, replacement string) {
	c.Print(
		AnsiSaveCursor,
		c.ANSIToMoveUp(lineUp),
		mansil.ClearLine,
		mansil.CursorColumn1,
		Ansi(replacement),
		AnsiRestoreCursor,
	)
}

// returns the ANSI code to move the cursor up n times,
// and an empty code when n <= 0
func (c Cursor) ANSIToMoveUp(n int) Ansi {
	return Ansi(mansil.CursorUp(n))
}

// returns the ANSI code to move the cursor down n times,
// and an empty code when n <= 0
func (c Cursor) ANSIToMoveDown(n int) Ansi {
	return Ansi(mansil.CursorDown(n))
}

// returns the ANSI code to move the cursor right n times,
// and an empty code when n <= 0
func (c Cursor) MoveRight(n int) Ansi {
	return Ansi(mansil.CursorRight(n))
}

// returns the ANSI code to move the cursor left n times,
// and an empty code when n <= 0
func (c Cursor) MoveLeft(n int) Ansi {
	return Ansi(mansil.CursorLeft(n))
}

func (c Cursor) writer() io.Writer {
	if c.W == nil {
		return os.Stdout
	}
	return c.W
}
