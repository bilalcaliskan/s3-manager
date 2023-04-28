package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPromptRunner(t *testing.T) {
	runner := GetPromptRunner("dummy prompt", nil)
	assert.NotNil(t, runner)
}

func TestGetSelectRunner(t *testing.T) {
	runner := GetSelectRunner("dummy prompt", nil)
	assert.NotNil(t, runner)
}

func TestPromptAccessCreds(t *testing.T) {

}
