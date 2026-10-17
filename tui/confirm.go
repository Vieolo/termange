package tui

import (
	"io"
	"strings"

	"github.com/vieolo/termange/cursor"
)

type ConfirmOptions struct {
	// The question that the user has to give a yes or no answer.
	// The prompt will be automatically appended with `[Y/n]` when
	// Default is true, and with `[y/N]` when it is false
	Prompt string
	// The answer when the user just presses Enter
	Default bool
	// If this function provided, the prompt will be deleted and the
	// return string will be printed instead
	PostConfirmationPrint func(bool) string
	// Where the answer is read from. When nil, it is os.Stdin
	In io.Reader
	// Where the prompt is written. When nil, it is os.Stdout
	Out io.Writer
}

// Shows the question and asks the user to enter yes or no.
//
// `y`, `yes`, `n` and `no` are accepted in any case, and an empty answer
// gives Default. Any other answer asks the question again.
// When the input ends before an answer, such as a closed stdin, it returns
// false whatever the default, since nobody answered
func Confirm(options ConfirmOptions) bool {
	in, out := streams(options.In, options.Out)
	prompt := confirmPrompt(options.Prompt, options.Default)

	value := false
	for {
		answer, err := ask(prompt, in, out)
		if err != nil {
			break
		}
		v, ok := parseConfirmation(answer, options.Default)
		if ok {
			value = v
			break
		}
	}

	if options.PostConfirmationPrint != nil {
		cursor.Cursor{W: out}.ReplaceLine(1, options.PostConfirmationPrint(value))
	}

	return value
}

// Returns the prompt ending with the hint of the default, replacing
// any hint that it already ends with
func confirmPrompt(prompt string, def bool) string {
	prompt = strings.TrimRight(prompt, " ")
	prompt = strings.TrimSuffix(prompt, "[Y/n]")
	prompt = strings.TrimSuffix(prompt, "[y/N]")
	prompt = strings.TrimRight(prompt, " ")
	if def {
		return prompt + " [Y/n]"
	}
	return prompt + " [y/N]"
}

// Returns the answer as a bool, and whether it is a valid answer
func parseConfirmation(answer string, def bool) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(answer)) {
	case "":
		return def, true
	case "y", "yes":
		return true, true
	case "n", "no":
		return false, true
	}
	return false, false
}
