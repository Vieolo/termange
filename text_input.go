package terminalutils

import (
	"fmt"
	"strings"
)

// Creates an text input and returns the input
func TextInput(prompt string) string {
	var value string
	p := prompt
	if !strings.Contains(p, ": ") {
		p = prompt + ": "
	}
	fmt.Print(p)
	fmt.Scan(&value)
	return value
}
