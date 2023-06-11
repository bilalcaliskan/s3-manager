package options

import (
	"testing"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"

	"github.com/stretchr/testify/assert"
)

func TestGetTransferAccelerationOptions(t *testing.T) {
	opts := GetTransferAccelerationOptions()
	assert.NotNil(t, opts)
}

func TestTransferAccelerationOptions_SetZeroValues(t *testing.T) {
	rootOpts := options.GetRootOptions()
	opts := GetTransferAccelerationOptions()
	opts.RootOptions = rootOpts
	assert.NotNil(t, opts)

	opts.SetZeroValues()
}
