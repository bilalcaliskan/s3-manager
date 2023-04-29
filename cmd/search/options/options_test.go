package options

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type promptMock struct {
	msg string
	err error
}

func (p promptMock) Run() (string, error) {
	// return expected result
	return p.msg, p.err
}

func TestGetSearchOptions(t *testing.T) {
	opts := GetSearchOptions()
	assert.NotNil(t, opts)
}

func TestSearchOptions_InitFlags(t *testing.T) {
	cmd := cobra.Command{}
	opts := GetSearchOptions()
	opts.InitFlags(&cmd)
}

func TestSearchOptions_PromptInteractiveValuesSubstringErr(t *testing.T) {
	opts := GetSearchOptions()

	substringRunnerOrg := substringRunner
	substringRunner = promptMock{
		msg: "",
		err: errors.New("dummy error"),
	}
	err := opts.PromptInteractiveValues()
	assert.NotNil(t, err)
	assert.Equal(t, opts.Substring, "")

	substringRunner = substringRunnerOrg
}

func TestSearchOptions_PromptInteractiveValuesSubstringSuccess(t *testing.T) {
	opts := GetSearchOptions()

	substringRunnerOrg := substringRunner
	substringRunner = promptMock{
		msg: "thisissubstring",
		err: nil,
	}
	err := opts.PromptInteractiveValues()
	assert.NotNil(t, err)
	assert.Equal(t, opts.Substring, "thisissubstring")

	substringRunner = substringRunnerOrg
}

func TestSearchOptions_PromptInteractiveValuesExtensionsErr(t *testing.T) {
	opts := GetSearchOptions()

	substringRunnerOrg := substringRunner
	extensionRunnerOrg := extensionRunner

	substringRunner = promptMock{
		msg: "thisissubstring",
		err: nil,
	}

	extensionRunner = promptMock{
		msg: "",
		err: errors.New("dummy error"),
	}

	err := opts.PromptInteractiveValues()
	assert.NotNil(t, err)
	assert.Equal(t, opts.Substring, "thisissubstring")

	substringRunner = substringRunnerOrg
	extensionRunner = extensionRunnerOrg
}

func TestSearchOptions_PromptInteractiveValuesExtensionsSuccess(t *testing.T) {
	opts := GetSearchOptions()

	substringRunnerOrg := substringRunner
	extensionRunnerOrg := extensionRunner

	substringRunner = promptMock{
		msg: "thisissubstring",
		err: nil,
	}

	extensionRunner = promptMock{
		msg: "thisisextensions",
		err: nil,
	}

	err := opts.PromptInteractiveValues()
	assert.Nil(t, err)
	assert.Equal(t, opts.Substring, "thisissubstring")
	assert.Equal(t, opts.FileExtensions, "thisisextensions")

	substringRunner = substringRunnerOrg
	extensionRunner = extensionRunnerOrg
}
