# Change Log

## v0.9.0
- Added the `live` package, to draw a block of lines that is redrawn in place, for progress, spinners and status trees
- Added `MultiSelect`, to select several items of a list, starting from an initial selection
- Added `Default` to `ConfirmOptions`, the answer when the user just presses Enter
- Added `In` and `Out` to `TextInputOptions` and `ConfirmOptions`, to read from and write to other streams than stdin and stdout
- Added `W` to `cursor.Cursor`, to write to other writers than stdout
- Added the colors black, blue, magenta and cyan
- Added the `example` program, which runs each part of termange
- termange now uses `github.com/vieolo/mansil` for its ANSI codes
- Fixed `TextInput` losing the lines after the first one when stdin is piped
- Fixed `TextInput` dropping a last line that has no line ending, such as from `printf y | app`
- Fixed `TextInput` keeping the `\r` of Windows line endings, which made `Confirm` always answer no on Windows
- Fixed `PostValuePrint` and `PostConfirmationPrint` replacing the line above the prompt when stdin is piped
- Fixed `Select` quitting when `q` is typed in the filter
- Fixed `Select` quitting on esc while a filter is applied. It clears the filter instead

#### Deprecated
- Deprecated the `Ansi*` constants of `cursor` in favour of the codes of mansil

#### Breaking Changes
- `Confirm`'s default value is now `false`. This is evident with the prompt (`[y/N]`)
- `Confirm` now asks again if the provided answer is other than yes, no, or empty. Before, an unknown answer would've been considered to be a no.
- In `Select`, enter while typing a filter accepts the filter, and a second enter picks the item. Before, the first enter picked the item under the cursor
- `ANSIToMoveUp`, `ANSIToMoveDown`, `MoveRight` and `MoveLeft` return an empty code for `n <= 0`, so `ClearPreviousLine(0)` and `ReplaceLine(0, ...)` act on the current line. Before, a count of 0 moved by one

## v0.8.0 (2026-05-15)

#### Breaking Changes
- Adjusted the `PrintError` functions to write to stderr instead of stdout. A new set of functions is added in case you would need to write to stdout

## v0.7.1 (2025-12-24)
- Added the `cursor` for movement of the cursor in terminal and manipulation of the lines in terminal

## v0.7.0 (2025-12-10)
- Added real, user, and system times to the `RunCommand` result

#### Breaking Changes
- Fixed the return error type of `RunCommand` functions

## v0.6.0 (2025-12-03)
- Added `go.yaml` for the project
- Added the ability to print start text before running a command

#### Breaking Changes
- Changed the return value of `RunCommand` functions. Instead of returning a generic error value and two buffers (which were confusing), `RunCommand` functions return `CommandResult` containing stdout and stderr, and `CommandError` which provides a typed struct for the errors

## v0.5.1 (2025-10-31)
- Added the `Confirm` TUI component
- Added post input callback for `TextInput`

## v0.5.0 (2025-10-01)
- Added support for `*Printf` functions
- Added print functions for warning and info

#### Breaking Changes
- Changed the type of the colors from normal string to `Color`
- Changed the `*Print` functions to `*Println` functions to better reflect the `fmt`

## v0.4.0 (2025-09-17)
- Added support for filtering in `Select`

#### Breaking Changes
- Moved the `Select` and `TextInput` to `tui` sub module
- Changed the function arguments of `TextInput` to instance of `TextInputOptions`

## v0.3.0 (2025-09-09)

#### Breaking Changes
- Renamed the package to `termange`

## v0.2.2 (2024-04-03)
- Added `Table`

## v0.2.1 (2023-09-05)
- Added the command utils

## v0.2.0 (2023-07-10)
- Updated the `Select`

#### Breaking Changes
- The arguments and return types of `Select` are changed

## v0.1.3 (2023-02-04)
- Added `Select`

## v0.1.2 (2022-11-17)
- Fixed the `TextInput` when the text was mistaken for terminal commands

## v0.1.1 (2022-11-17)
- Fixed the accessibility of the colors in `log`

## v0.1.0 (2022-11-15)
- Added `TextInput`
- Changed the Print functions from `log` to `fmt`