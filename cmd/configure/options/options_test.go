package options

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetConfigureOptions(t *testing.T) {
	opts := GetConfigureOptions()
	assert.NotNil(t, opts)
}

func TestGetConfigureOptions_InitFlags(t *testing.T) {
	cmd := cobra.Command{}
	opts := GetConfigureOptions()
	opts.InitFlags(&cmd)
	opts.SetZeroValues()
}
