//go:build unit

package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTransferAccelerationOptions(t *testing.T) {
	opts := GetTransferAccelerationOptions()
	assert.NotNil(t, opts)
}
