//go:build unit
// +build unit

package utils

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	options2 "github.com/bilalcaliskan/s3-manager/cmd/versioning/options"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// Define a testdata struct to be used in your unit tests
type mockS3Client struct {
	s3iface.S3API
}

func TestCheckArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected error
	}{
		{
			name:     "Failure caused by too many arguments",
			args:     []string{"foo", "bar"},
			expected: errors.New(ErrTooManyArguments),
		},
		{
			name:     "Success",
			args:     []string{},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CheckArgs(test.args)
			assert.Equal(t, test.expected, err)
		})
	}
}

func TestPrepareConstants(t *testing.T) {
	var (
		svc            s3iface.S3API
		versioningOpts *options2.VersioningOptions
		logger         zerolog.Logger
	)

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	rootOpts := rootopts.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	mockSvc := &mockS3Client{}
	svc = mockSvc
	assert.NotNil(t, mockSvc)

	cmd.SetContext(context.WithValue(context.Background(), rootopts.OptsKey{}, rootOpts))
	cmd.SetContext(context.WithValue(cmd.Context(), rootopts.S3SvcKey{}, svc))

	svc, versioningOpts, logger = PrepareConstants(cmd, options2.GetVersioningOptions())
	assert.NotNil(t, svc)
	assert.NotNil(t, versioningOpts)
	assert.NotNil(t, logger)
}

func TestDecideActualState(t *testing.T) {
	tests := []struct {
		name     string
		res      *s3.GetBucketVersioningOutput
		expected error
	}{
		{
			name: "Success enabled",
			res: &s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			},
			expected: nil,
		},
		{
			name: "Success suspended",
			res: &s3.GetBucketVersioningOutput{
				Status: aws.String("Suspended"),
			},
			expected: nil,
		},
		{
			name: "Failure caused by unknown state",
			res: &s3.GetBucketVersioningOutput{
				Status: aws.String("Suspendeddd"),
			},
			expected: fmt.Errorf(ErrUnknownStatus, "Suspendeddd"),
		},
	}

	opts := options2.GetVersioningOptions()
	for _, test := range tests {
		err := DecideActualState(test.res, opts)
		assert.Equal(t, test.expected, err)
	}
}
