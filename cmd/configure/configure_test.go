package configure

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

/*
	func TestExecuteMissingRegion(t *testing.T) {
		rootOpts := options.GetRootOptions()
		rootOpts.AccessKey = "thisisaccesskey"
		rootOpts.SecretKey = "thisissecretkey"
		rootOpts.Region = ""
		rootOpts.BucketName = "thisisbucketname"
		rootOpts.Interactive = false

		svc, err := createSvc(rootOpts)


		ConfigureCmd.SetContext(context.WithValue(context.Background(), options.OptsKey{}, rootOpts))
		err := ConfigureCmd.Execute()
		assert.NotNil(t, err)

		rootOpts.SetZeroValues()
		configureOpts.SetZeroValues()
	}
*/
func TestExecute(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"
	rootOpts.Interactive = false

	ctx := context.Background()
	ConfigureCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	ConfigureCmd.SetContext(context.WithValue(ConfigureCmd.Context(), options.S3SvcKey{}, svc))
	ConfigureCmd.SetContext(context.WithValue(ConfigureCmd.Context(), options.OptsKey{}, rootOpts))

	err = ConfigureCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	configureOpts.SetZeroValues()
}

func TestExecuteFailingPutRequest(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"
	rootOpts.Interactive = true

	ctx := context.Background()
	ConfigureCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	ConfigureCmd.SetContext(context.WithValue(ConfigureCmd.Context(), options.S3SvcKey{}, svc))
	ConfigureCmd.SetContext(context.WithValue(ConfigureCmd.Context(), options.OptsKey{}, rootOpts))

	err = ConfigureCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	configureOpts.SetZeroValues()
}
