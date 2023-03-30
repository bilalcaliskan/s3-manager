package logging

import (
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetLogger function tests if GetLogger function running properly
func TestGetLogger(t *testing.T) {
	logger := GetLogger(options.GetRootOptions())
	assert.NotNil(t, logger)
}

func TestEnableDebugLogging(t *testing.T) {
	EnableDebugLogging()
}
