package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersioningOptions(t *testing.T) {
	opts := GetVersioningOptions()
	assert.NotNil(t, opts)
}

func TestVersioningOptions_SetZeroValues(t *testing.T) {
	opts := GetVersioningOptions()
	assert.NotNil(t, opts)

	opts.SetZeroValues()
}
