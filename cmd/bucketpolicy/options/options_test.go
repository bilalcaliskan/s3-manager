//go:build unit
// +build unit

package options

import (
	"testing"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"

	"github.com/stretchr/testify/assert"
)

func TestGetBucketPolicyOptions(t *testing.T) {
	opts := GetBucketPolicyOptions()
	assert.NotNil(t, opts)
}

func TestBucketPolicyOptions_SetZeroValues(t *testing.T) {
	rootOpts := options.GetRootOptions()
	opts := GetBucketPolicyOptions()
	opts.RootOptions = rootOpts
	assert.NotNil(t, opts)

	opts.SetZeroValues()
}
