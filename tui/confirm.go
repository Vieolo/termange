package tui

import (
	"fmt"
	"strings"

	"github.com/vieolo/termange/cursor"
)

type ConfirmOptions struct {
	// The question that the user has to give a yes or no answer.
	// The prompt will be automatically appended with `[Y/n]`
	Prompt string
	// If this function provided, the prompt will be deleted and the
	// return string will be printed instead
	PostConfirmationPrint func(bool) string
}

// Shows the question and asks the user to enter Y or n
func Confirm(options ConfirmOptions) bool {
	finalPrompt := options.Prompt
	if !strings.HasSuffix(finalPrompt, " [Y/n]") {
		finalPrompt += " [Y/n]"
	}
	answer := TextInput(TextInputOptions{
		Prompt: finalPrompt,
	})

	answer = strings.ToLower(answer)
	value := answer == "yes" || answer == "y"

	if options.PostConfirmationPrint != nil {
		cursor.Cursor{}.
			MoveUp(1).
			ClearLine()

		fmt.Println(options.PostConfirmationPrint(value))
	}

	return value
}
