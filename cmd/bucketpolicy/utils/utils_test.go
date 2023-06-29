//go:build unit

package utils

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/bucketpolicy/options"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func createSvc(rootOpts *rootopts.RootOptions) (*s3.S3, error) {
	return internalaws.CreateAwsService(rootOpts)
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
			expected: errors.New(errTooManyArgsProvided),
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
		svc              s3iface.S3API
		bucketPolicyOpts *options.BucketPolicyOptions
		logger           zerolog.Logger
	)

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	rootOpts := rootopts.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	cmd.SetContext(context.WithValue(context.Background(), rootopts.OptsKey{}, rootOpts))
	cmd.SetContext(context.WithValue(cmd.Context(), rootopts.S3SvcKey{}, svc))

	svc, bucketPolicyOpts, logger = PrepareConstants(cmd, options.GetBucketPolicyOptions())
	assert.NotNil(t, svc)
	assert.NotNil(t, bucketPolicyOpts)
	assert.NotNil(t, logger)
}
