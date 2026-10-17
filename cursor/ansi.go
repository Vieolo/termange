package cursor

import "github.com/vieolo/mansil"

// ANSI escape codes for cursor manipulation
type Ansi string

// The codes come from mansil, and these names are kept for existing code.
// Use mansil directly in new code, such as `mansil.CursorUp(3)` instead of
// `AnsiCursorUp` and `Cursor.ANSIToMoveUp`
const (
	// Movement

	// Deprecated: use mansil.CursorUp1 or mansil.CursorUp
	AnsiCursorUp = mansil.CursorUp1
	// Deprecated: use mansil.CursorDown1 or mansil.CursorDown
	AnsiCursorDown = mansil.CursorDown1
	// Deprecated: use mansil.CursorRight1 or mansil.CursorRight
	AnsiCursorRight = mansil.CursorRight1
	// Deprecated: use mansil.CursorLeft1 or mansil.CursorLeft
	AnsiCursorLeft = mansil.CursorLeft1

	// Line Control

	// Moves to beginning of next line
	//
	// Deprecated: use mansil.CursorNextLine1 or mansil.CursorNextLine
	AnsiCursorNextLine = mansil.CursorNextLine1
	// Moves to beginning of previous line
	//
	// Deprecated: use mansil.CursorPrevLine1 or mansil.CursorPrevLine
	AnsiCursorPrevLine = mansil.CursorPrevLine1
	// Moves cursor to the first column
	//
	// Deprecated: use mansil.CursorColumn1
	AnsiCursorToColumn0 = mansil.CursorColumn1

	// Erasing

	// Erases the entire current line
	//
	// Deprecated: use mansil.ClearLine
	AnsiClearLine = mansil.ClearLine
	// Erases from cursor to end of line
	//
	// Deprecated: use mansil.ClearLineRight
	AnsiClearLineToRight = mansil.ClearLineRight
	// Clears entire screen
	//
	// Deprecated: use mansil.ClearScreen
	AnsiClearScreen = mansil.ClearScreen

	// Position Management
	// mansil has no codes to save and restore the cursor yet

	AnsiSaveCursor    = "\033[s" // Saves the current cursor position
	AnsiRestoreCursor = "\033[u" // Restores the last saved position

	// Visibility

	// Deprecated: use mansil.HideCursor
	AnsiHideCursor = mansil.HideCursor
	// Deprecated: use mansil.ShowCursor
	AnsiShowCursor = mansil.ShowCursor
)
