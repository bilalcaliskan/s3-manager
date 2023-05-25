package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: uncomment when interactivity enabled again
/*func TestGetPromptRunner(t *testing.T) {
	runner := GetPromptRunner("dummy prompt", nil)
	assert.NotNil(t, runner)
}

func TestGetSelectRunner(t *testing.T) {
	runner := GetSelectRunner("dummy prompt", nil)
	assert.NotNil(t, runner)
}*/

func TestGetPromptRunner(t *testing.T) {
	runner := GetPromptRunner("dummy prompt", false, nil)
	assert.NotNil(t, runner)
}
