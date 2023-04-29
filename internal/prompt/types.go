package prompt

import "github.com/manifoldco/promptui"

type SelectRunner interface {
	Run() (int, string, error)
}

func GetSelectRunner(label string, items []string) *promptui.Select {
	return &promptui.Select{
		Label: label,
		Items: items,
	}
}

type PromptRunner interface {
	Run() (string, error)
}

func GetPromptRunner(label string, valFunc func(s string) error) *promptui.Prompt {
	return &promptui.Prompt{
		Label:    label,
		Validate: valFunc,
	}
}
