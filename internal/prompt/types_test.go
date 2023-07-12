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
