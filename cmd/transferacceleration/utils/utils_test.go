//go:build unit

package utils

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	options3 "github.com/bilalcaliskan/s3-manager/cmd/transferacceleration/options"
	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/stretchr/testify/assert"
)

func createSvc(rootOpts *options.RootOptions) (*s3.S3, error) {
	return internalaws.CreateAwsService(rootOpts)
}

func TestCheckArgsSuccess(t *testing.T) {
	err := CheckArgs([]string{})
	assert.Nil(t, err)
}

func TestCheckArgsFailure(t *testing.T) {
	err := CheckArgs([]string{"foo"})
	assert.NotNil(t, err)
}

func TestPrepareConstants(t *testing.T) {
	var (
		svc    s3iface.S3API
		taOpts *options3.TransferAccelerationOptions
		logger zerolog.Logger
	)

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	cmd.SetContext(context.WithValue(context.Background(), options.OptsKey{}, rootOpts))
	cmd.SetContext(context.WithValue(cmd.Context(), options.S3SvcKey{}, svc))

	svc, taOpts, logger = PrepareConstants(cmd, options3.GetTransferAccelerationOptions())
	assert.NotNil(t, svc)
	assert.NotNil(t, taOpts)
	assert.NotNil(t, logger)
}
