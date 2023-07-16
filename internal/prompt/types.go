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

type PromptMock struct {
	Msg string
	Err error
}

func (p PromptMock) Run() (string, error) {
	return p.Msg, p.Err
}

type PromptWrapper struct {
	Prompt    *promptui.Prompt
	UserInput string
}

func (p *PromptWrapper) Run() (string, error) {
	return p.UserInput, p.Prompt.Validate(p.UserInput)
}
