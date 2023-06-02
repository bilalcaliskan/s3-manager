package clean

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/internal/aws"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

func createSvc(rootOpts *options.RootOptions) (*s3.S3, error) {
	return aws.CreateAwsService(rootOpts)
}

func TestExecuteMissingRegion(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = ""

	ctx := context.Background()
	CleanCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	CleanCmd.SetContext(context.WithValue(CleanCmd.Context(), options.S3SvcKey{}, svc))
	CleanCmd.SetContext(context.WithValue(CleanCmd.Context(), options.OptsKey{}, rootOpts))

	err = CleanCmd.Execute()
	assert.NotNil(t, err)

	cleanOpts.SetZeroValues()
}

func TestExecuteInvalidSortByOption(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	CleanCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	CleanCmd.SetContext(context.WithValue(CleanCmd.Context(), options.S3SvcKey{}, svc))
	CleanCmd.SetContext(context.WithValue(CleanCmd.Context(), options.OptsKey{}, rootOpts))

	err = CleanCmd.Flags().Set("sortBy", "nonexistedsortbyflag")
	assert.Nil(t, err)

	err = CleanCmd.Execute()
	assert.NotNil(t, err)

	cleanOpts.SetZeroValues()
}

func TestExecuteInvalidMinMaxValues(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	CleanCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	CleanCmd.SetContext(context.WithValue(CleanCmd.Context(), options.S3SvcKey{}, svc))
	CleanCmd.SetContext(context.WithValue(CleanCmd.Context(), options.OptsKey{}, rootOpts))

	err = CleanCmd.Flags().Set("minFileSizeInMb", "20")
	assert.Nil(t, err)

	err = CleanCmd.Flags().Set("maxFileSizeInMb", "10")
	assert.Nil(t, err)

	err = CleanCmd.Execute()
	assert.NotNil(t, err)

	cleanOpts.SetZeroValues()
}

func TestExecute(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	CleanCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	CleanCmd.SetContext(context.WithValue(CleanCmd.Context(), options.S3SvcKey{}, svc))
	CleanCmd.SetContext(context.WithValue(CleanCmd.Context(), options.OptsKey{}, rootOpts))

	err = CleanCmd.Execute()
	assert.NotNil(t, err)

	cleanOpts.SetZeroValues()
}

func TestExecuteSuccess(t *testing.T) {
	// TODO: implement
}
