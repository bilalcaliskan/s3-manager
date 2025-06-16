//go:build unit

package prompt

import (
	"github.com/bilalcaliskan/s3-manager/internal/pkg/constants"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPromptRunner(t *testing.T) {
	runner := GetPromptRunner("dummy prompt", false, nil)
	assert.NotNil(t, runner)
}

func TestGetConfirmRunner(t *testing.T) {
	testCases := []struct {
		caseName  string
		input     string
		expectErr bool
	}{
		{
			caseName:  "Valid input",
			input:     "y",
			expectErr: false,
		},
		{
			caseName:  "Invalid input",
			input:     "yy",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.caseName, func(t *testing.T) {
			t.Logf("starting case %s", tc.caseName)

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

func TestAskForApproval(t *testing.T) {
	testCases := []struct {
		caseName  string
		mock      *PromptMock
		expectErr bool
	}{
		{
			caseName: "Approve",
			mock: &PromptMock{
				Msg: "y",
				Err: nil,
			},
			expectErr: false,
		},
		{
			caseName: "Terminate",
			mock: &PromptMock{
				Msg: "n",
				Err: constants.ErrUserTerminated,
			},
			expectErr: true,
		},
		{
			caseName: "Invalid input",
			mock: &PromptMock{
				Msg: "adsjlkfasd",
				Err: constants.ErrInvalidInput,
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.caseName, func(t *testing.T) {
			t.Logf("starting case %s", tc.caseName)

			err := AskForApproval(tc.mock)
			if tc.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
