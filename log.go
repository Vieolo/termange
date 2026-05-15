package termange

import (
	"fmt"
	"os"
	"strings"
)

type Color string

var (
	ColorRed    Color = "red"
	ColorGreen  Color = "green"
	ColorYellow Color = "yellow"
	ColorWhite  Color = "white"
)

// Adds color to the given text
func PaintText(text string, color Color) string {
	reset := "\033[0m"
	redC := "\033[31m"
	greenC := "\033[32m"
	yellowC := "\033[33m"
	whiteC := "\033[1;37m"
	if strings.ToLower(string(color)) == string(ColorRed) {
		return redC + text + reset
	} else if strings.ToLower(string(color)) == string(ColorGreen) {
		return greenC + text + reset
	} else if strings.ToLower(string(color)) == string(ColorYellow) {
		return yellowC + text + reset
	} else if strings.ToLower(string(color)) == string(ColorWhite) {
		return whiteC + text + reset
	}

	return text
}

// Prints a new line with the given color
func PrintColorln(text string, color Color) {
	fmt.Println(PaintText(text, color))
}

// Prints a new line with the color Green
func PrintSuccessln(text string) {
	PrintColorln(text, ColorGreen)
}

// Prints a new error to stderr with the color Red
//
// Notice that the function prints to stderr, if you
// prefer to print to stdout, use `PrintErrorlnStdout`
func PrintErrorln(text string) {
	fmt.Fprintln(os.Stderr, PaintText(text, ColorRed))
}

// Prints a new error to stdout with the color Red
//
// Notice that the function prints to stdout, if you
// prefer to print to stderr, use `PrintErrorln`
func PrintErrorlnStdout(text string) {
	PrintColorln(text, ColorRed)
}

// Prints a new line with the color yellow
func PrintWarningln(text string) {
	PrintColorln(text, ColorYellow)
}

// Prints a new line with the color white
func PrintInfoln(text string) {
	PrintColorln(text, ColorWhite)
}

// Prints with the given color
func PrintColorf(color Color, format string, a ...any) {
	v := fmt.Sprintf(format, a...)
	fmt.Print(PaintText(v, color))
}

// Prints the color Green
func PrintSuccessf(format string, a ...any) {
	PrintColorf(ColorGreen, format, a...)
}

// Prints to the stderr with the color Red
//
// Please note that this function writes to stderr, if you
// prefer to write to stdout, using `PrintErrorfStdout`
func PrintErrorf(format string, a ...any) {
	v := fmt.Sprintf(format, a...)
	fmt.Fprint(os.Stderr, PaintText(v, ColorRed))
}

// Prints to the stdout with the color Red
//
// Please note that this function writes to stdout, if you
// prefer to write to stderr, using `PrintErrorf`
func PrintErrorfStdout(format string, a ...any) {
	PrintColorf(ColorRed, format, a...)
}

// Prints with the color Yellow
func PrintWarningf(format string, a ...any) {
	PrintColorf(ColorYellow, format, a...)
}

// Prints with the color White
func PrintInfof(format string, a ...any) {
	PrintColorf(ColorWhite, format, a...)
}
