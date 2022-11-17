package terminalutils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Creates an text input and returns the input
func TextInput(prompt string) string {
	p := prompt
	if !strings.Contains(p, ": ") {
		p = prompt + ": "
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(p)
	text, te := reader.ReadString('\n')
	if te != nil {
		return ""
	}
	return strings.TrimRight(text, "\n")
}
