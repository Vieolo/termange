package tui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/vieolo/termange/internal"
)

type TextInputOptions struct {
	// The question/info to the user on what they should enter
	Prompt string
	// If this function provided, the prompt will be deleted and the
	// return string will be printed instead
	PostValuePrint func(string) string
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

	finalValue := strings.TrimRight(text, "\n")

	if options.PostValuePrint != nil {
		internal.IT{}.
			CursorUp().
			ClearLine()

		fmt.Println(options.PostValuePrint(finalValue))
	}

	return finalValue
}
