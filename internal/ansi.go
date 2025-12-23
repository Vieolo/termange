package internal

// ANSI escape codes for cursor manipulation
const (
	// Movement
	CursorUp    = "\033[1A"
	CursorDown  = "\033[1B"
	CursorRight = "\033[1C"
	CursorLeft  = "\033[1D"

	// Line Control
	CursorNextLine  = "\033[1E" // Moves to beginning of next line
	CursorPrevLine  = "\033[1F" // Moves to beginning of previous line
	CursorToColumn0 = "\033[1G" // Moves cursor to the first column

	// Erasing
	ClearLine        = "\033[2K" // Erases the entire current line
	ClearLineToRight = "\033[0K" // Erases from cursor to end of line
	ClearScreen      = "\033[2J" // Clears entire screen and moves cursor to home

	// Position Management
	SaveCursor    = "\033[s" // Saves the current cursor position
	RestoreCursor = "\033[u" // Restores the last saved position

	// Visibility
	HideCursor = "\033[?25l"
	ShowCursor = "\033[?25h"
)
