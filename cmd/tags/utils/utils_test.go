//go:build unit
// +build unit

package utils

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	options2 "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/cmd/tags/options"
	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func createSvc(rootOpts *options2.RootOptions) (*s3.S3, error) {
	return internalaws.CreateAwsService(rootOpts)
}

func TestCheckArgs(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *cobra.Command
		args     []string
		expected error
	}{
		{
			name:     "Show command with no arguments",
			cmd:      &cobra.Command{Use: "show"},
			args:     []string{},
			expected: nil,
		},
		{
			name:     "Success",
			cmd:      &cobra.Command{Use: "anothercommand"},
			args:     []string{"foo"},
			expected: nil,
		},
		{
			name:     "Show command with arguments",
			cmd:      &cobra.Command{Use: "show"},
			args:     []string{"arg1"},
			expected: errors.New("too many argument provided"),
		},
		{
			name:     "Command with no arguments",
			cmd:      &cobra.Command{},
			args:     []string{},
			expected: errors.New("no argument provided"),
		},
		{
			name:     "Command with too many arguments",
			cmd:      &cobra.Command{},
			args:     []string{"arg1", "arg2"},
			expected: errors.New("too many argument provided"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CheckArgs(test.cmd, test.args)
			assert.Equal(t, test.expected, err)
		})
	}
}

func TestPrepareConstants(t *testing.T) {
	var (
		svc     s3iface.S3API
		tagOpts *options.TagOptions
		logger  zerolog.Logger
	)

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	rootOpts := options2.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	cmd.SetContext(context.WithValue(context.Background(), options2.OptsKey{}, rootOpts))
	cmd.SetContext(context.WithValue(cmd.Context(), options2.S3SvcKey{}, svc))

	svc, tagOpts, logger = PrepareConstants(cmd, options.GetTagOptions())
	assert.NotNil(t, svc)
	assert.NotNil(t, tagOpts)
	assert.NotNil(t, logger)
}
