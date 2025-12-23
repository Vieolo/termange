package cursor

// ANSI escape codes for cursor manipulation
type Ansi string

const (
	// Movement
	AnsiCursorUp    = "\033[1A"
	AnsiCursorDown  = "\033[1B"
	AnsiCursorRight = "\033[1C"
	AnsiCursorLeft  = "\033[1D"

	// Line Control
	AnsiCursorNextLine  = "\033[1E" // Moves to beginning of next line
	AnsiCursorPrevLine  = "\033[1F" // Moves to beginning of previous line
	AnsiCursorToColumn0 = "\033[1G" // Moves cursor to the first column

	// Erasing
	AnsiClearLine        = "\033[2K" // Erases the entire current line
	AnsiClearLineToRight = "\033[0K" // Erases from cursor to end of line
	AnsiClearScreen      = "\033[2J" // Clears entire screen and moves cursor to home

	// Position Management
	AnsiSaveCursor    = "\033[s" // Saves the current cursor position
	AnsiRestoreCursor = "\033[u" // Restores the last saved position

	// Visibility
	AnsiHideCursor = "\033[?25l"
	AnsiShowCursor = "\033[?25h"
)
