//go:build unit

package transferacceleration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferAccelerationCmd(t *testing.T) {
	assert.NotNil(t, TransferAccelerationCmd)
}
