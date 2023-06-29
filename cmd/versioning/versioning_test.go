//go:build unit

package versioning

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersioningCmd(t *testing.T) {
	assert.NotNil(t, VersioningCmd)
}
