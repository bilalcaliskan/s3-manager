package prompt

import (
	"errors"

	"github.com/manifoldco/promptui"
)

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

func GetConfirmRunner() *promptui.Prompt {
	return GetPromptRunner("Confirm? (y/N)", true, func(s string) error {
		if len(s) == 1 {
			return nil
		}

		return errors.New("invalid input")
	})
}
