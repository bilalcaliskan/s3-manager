//go:build unit

package options

import (
	"testing"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetCleanOptions(t *testing.T) {
	opts := GetCleanOptions()
	rootOpts := options.GetRootOptions()
	opts.RootOptions = rootOpts
	assert.NotNil(t, opts)
}

func TestCleanOptions_InitFlags(t *testing.T) {
	cmd := cobra.Command{}

	rootOpts := options.GetRootOptions()
	opts := GetCleanOptions()
	opts.RootOptions = rootOpts

	opts.InitFlags(&cmd)
}
