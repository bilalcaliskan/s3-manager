//go:build unit

package options

import (
	"testing"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"

	"github.com/stretchr/testify/assert"
)

func TestGetVersioningOptions(t *testing.T) {
	opts := GetVersioningOptions()
	assert.NotNil(t, opts)
}

func TestVersioningOptions_SetZeroValues(t *testing.T) {
	rootOpts := options.GetRootOptions()
	opts := GetVersioningOptions()
	opts.RootOptions = rootOpts
	assert.NotNil(t, opts)

	opts.SetZeroValues()
}
