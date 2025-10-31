# Change Log

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