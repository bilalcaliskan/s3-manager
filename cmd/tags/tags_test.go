//go:build unit
// +build unit

package tags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersioningCmd(t *testing.T) {
	assert.NotNil(t, TagsCmd)
}
