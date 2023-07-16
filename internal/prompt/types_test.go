//go:build unit

package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPromptRunner(t *testing.T) {
	runner := GetPromptRunner("dummy prompt", false, nil)
	assert.NotNil(t, runner)
}

func TestGetConfirmRunner(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{
			name:      "Valid input",
			input:     "y",
			expectErr: false,
		},
		{
			name:      "Invalid input",
			input:     "yy",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prompt := GetConfirmRunner()

			// Wrap prompt into PromptWrapper
			wrapper := &PromptWrapper{
				Prompt:    prompt,
				UserInput: tc.input,
			}

			_, err := wrapper.Run()
			if tc.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestPromptMock_Run(t *testing.T) {
	runner := PromptMock{
		Msg: "y",
		Err: nil,
	}

	_, err := runner.Run()
	assert.Nil(t, err)
}
