//go:build unit

package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBucketPolicyOptions(t *testing.T) {
	opts := GetBucketPolicyOptions()
	assert.NotNil(t, opts)
}

func TestBucketPolicyOptions_SetZeroValues(t *testing.T) {
	opts := GetBucketPolicyOptions()
	assert.NotNil(t, opts)

	opts.SetZeroValues()
}
