package prompt

import "github.com/manifoldco/promptui"

// TODO: uncomment when interactivity enabled again
/*type SelectRunner interface {
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
}*/

// TODO: remove when interactivity enabled again

type PromptRunner interface {
	Run() (string, error)
}

func GetPromptRunner(label string, isConfirm bool, valFunc func(s string) error) *promptui.Prompt {
	return &promptui.Prompt{
		Label:     label,
		IsConfirm: isConfirm,
		Validate:  valFunc,
	}
}
