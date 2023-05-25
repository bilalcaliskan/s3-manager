package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchCmd(t *testing.T) {
	assert.NotNil(t, SearchCmd)
}
