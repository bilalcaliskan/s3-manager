//go:build unit
// +build unit

package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetCmd(t *testing.T) {
	assert.NotNil(t, SetCmd)
}
