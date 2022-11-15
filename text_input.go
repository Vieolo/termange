package terminalutils

import (
	"fmt"
	"strings"
)

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
