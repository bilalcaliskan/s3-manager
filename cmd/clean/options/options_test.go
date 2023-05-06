package options

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetCleanOptions(t *testing.T) {
	opts := GetCleanOptions()
	assert.NotNil(t, opts)
	opts.SetZeroValues()
}

func TestCleanOptions_InitFlags(t *testing.T) {
	cmd := cobra.Command{}

	opts := GetCleanOptions()
	opts.InitFlags(&cmd)
	opts.SetZeroValues()
}
