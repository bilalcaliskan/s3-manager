package options

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetSearchOptions(t *testing.T) {
	opts := GetSearchOptions()
	assert.NotNil(t, opts)
}

func TestSearchOptions_InitFlags(t *testing.T) {
	cmd := cobra.Command{}
	opts := GetSearchOptions()
	opts.InitFlags(&cmd)
}
