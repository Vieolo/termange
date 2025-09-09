package termange

import (
	"fmt"
	"strings"
)

var (
	Red    = "red"
	Green  = "green"
	Yellow = "yellow"
)

// Adds color to the given text
func PaintText(text string, color string) string {
	reset := "\033[0m"
	redC := "\033[31m"
	greenC := "\033[32m"
	yellowC := "\033[33m"
	if strings.ToLower(color) == Red {
		return redC + text + reset
	} else if strings.ToLower(color) == Green {
		return greenC + text + reset
	} else if strings.ToLower(color) == Yellow {
		return yellowC + text + reset
	}

	return text
}

// Prints a new line with the given color
func PrintColorln(text string, color string) {
	fmt.Println(PaintText(text, color))
}

// Prints a new line with the color Green
func PrintSuccess(text string) {
	PrintColorln(text, "green")
}

// Prints a new line with the color Red
func PrintError(text string) {
	PrintColorln(text, "red")
}
