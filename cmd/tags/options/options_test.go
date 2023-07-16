//go:build unit

package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTagOptions(t *testing.T) {
	opts := GetTagOptions()
	assert.NotNil(t, opts)
}
