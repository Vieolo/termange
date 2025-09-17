package tui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type TextInputOptions struct {
	// The question/info to the user on what they should enter
	Prompt string
}

// Creates an text input and returns the input
func TextInput(options TextInputOptions) string {
	p := options.Prompt
	if !strings.HasSuffix(p, " ") {
		p = options.Prompt + " "
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(p)
	text, te := reader.ReadString('\n')
	if te != nil {
		return ""
	}
	return strings.TrimRight(text, "\n")
}
