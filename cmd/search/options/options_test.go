//go:build unit

package options

import (
	"testing"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// TODO: uncomment when interactivity enabled again
/*type promptMock struct {
	msg string
	err error
}

func (p promptMock) Run() (string, error) {
	// return expected result
	return p.msg, p.err
}*/

func TestGetSearchOptions(t *testing.T) {
	opts := GetSearchOptions()
	assert.NotNil(t, opts)
}

func TestSearchOptions_InitFlags(t *testing.T) {
	cmd := cobra.Command{}
	cmd.Use = "text"

	opts := GetSearchOptions()
	rootOpts := options.GetRootOptions()
	opts.RootOptions = rootOpts
	opts.InitFlags(&cmd)
	opts.SetZeroValues()
}

// TODO: uncomment when interactivity enabled again
/*func TestSearchOptions_PromptInteractiveValuesSubstringErr(t *testing.T) {
	opts := GetSearchOptions()

	substringRunnerOrg := substringRunner
	substringRunner = promptMock{
		msg: "",
		err: errors.New("dummy error"),
	}
	err := opts.PromptInteractiveValues()
	assert.NotNil(t, err)
	assert.Equal(t, opts.Text, "")

	substringRunner = substringRunnerOrg
	opts.SetZeroValues()
}*/

// TODO: uncomment when interactivity enabled again
/*func TestSearchOptions_PromptInteractiveValuesSubstringSuccess(t *testing.T) {
	opts := GetSearchOptions()

	substringRunnerOrg := substringRunner
	substringRunner = promptMock{
		msg: "thisissubstring",
		err: nil,
	}
	err := opts.PromptInteractiveValues()
	assert.NotNil(t, err)
	assert.Equal(t, opts.Text, "thisissubstring")

	substringRunner = substringRunnerOrg
	opts.SetZeroValues()
}*/

// TODO: uncomment when interactivity enabled again
/*func TestSearchOptions_PromptInteractiveValuesExtensionsErr(t *testing.T) {
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
	assert.Equal(t, opts.Text, "thisissubstring")

	substringRunner = substringRunnerOrg
	extensionRunner = extensionRunnerOrg
	opts.SetZeroValues()
}*/

// TODO: uncomment when interactivity enabled again
/*func TestSearchOptions_PromptInteractiveValuesExtensionsSuccess(t *testing.T) {
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
	assert.Equal(t, opts.Text, "thisissubstring")
	assert.Equal(t, opts.FileExtensions, "thisisextensions")

	substringRunner = substringRunnerOrg
	extensionRunner = extensionRunnerOrg
	opts.SetZeroValues()
}*/
