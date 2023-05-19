package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	opts := GetTagOptions()
	assert.NotNil(t, opts)
}
