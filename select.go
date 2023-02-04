package terminalutils

import "github.com/manifoldco/promptui"

func Select(options []string, label string) (int, string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: options,
	}

	return prompt.Run()
}
