//go:build unit
// +build unit

package bucketpolicy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBucketPolicyCmd(t *testing.T) {
	assert.NotNil(t, BucketPolicyCmd)
}
