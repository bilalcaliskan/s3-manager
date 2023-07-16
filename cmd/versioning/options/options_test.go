//go:build unit

package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersioningOptions(t *testing.T) {
	opts := GetVersioningOptions()
	assert.NotNil(t, opts)
}
