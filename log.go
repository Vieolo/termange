package terminalutils

import (
	"log"
	"strings"
)

func PaintText(text string, color string) string {
	reset := "\033[0m"
	red := "\033[31m"
	green := "\033[32m"
	yellow := "\033[33m"
	if strings.ToLower(color) == "red" {
		return red + text + reset
	} else if strings.ToLower(color) == "green" {
		return green + text + reset
	} else if strings.ToLower(color) == "yellow" {
		return yellow + text + reset
	}

	return text
}

func PrintColorln(text string, color string) {
	log.Println(PaintText(text, color))
}

func PrintSuccess(text string) {
	PrintColorln(text, "green")
}

func PrintError(text string) {
	PrintColorln(text, "red")
}
