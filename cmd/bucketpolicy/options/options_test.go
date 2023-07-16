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
