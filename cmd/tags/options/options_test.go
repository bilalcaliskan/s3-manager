//go:build unit

package options

import (
	"testing"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetTagOptions(t *testing.T) {
	opts := GetTagOptions()
	assert.NotNil(t, opts)
}

func TestTagOptions_SetZeroValues(t *testing.T) {
	cmd := cobra.Command{}
	opts := GetTagOptions()
	assert.NotNil(t, opts)

	rootOpts := options.GetRootOptions()
	opts.RootOptions = rootOpts

	opts.InitFlags(&cmd)
	opts.SetZeroValues()
}
