//go:build unit

package options

import (
	"testing"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetSearchOptions(t *testing.T) {
	opts := GetSearchOptions()
	assert.NotNil(t, opts)
}

func TestSearchOptions_InitFlags(t *testing.T) {
	cmd := cobra.Command{}
	cmd.Use = "text"

	opts := GetSearchOptions()
	rootOpts := options.GetRootOptions()
	opts.RootOptions = rootOpts
	opts.InitFlags(&cmd)
}

func TestSearchOptions_SetZeroValues(t *testing.T) {
	opts := GetSearchOptions()
	opts.SetZeroValues()
}
