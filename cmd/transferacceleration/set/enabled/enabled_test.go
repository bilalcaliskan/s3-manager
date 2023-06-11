package enabled

import (
	"context"
	"errors"
	"testing"

	"github.com/bilalcaliskan/s3-manager/cmd/versioning/set/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/stretchr/testify/assert"
)

var (
	defaultGetBucketAccelerationOutput = &s3.GetBucketAccelerateConfigurationOutput{
		Status: aws.String("Enabled"),
	}
	defaultGetBucketAccelerationErr    error
	defaultPutBucketAccelerationOutput = &s3.PutBucketAccelerateConfigurationOutput{}
	defaultPutBucketAccelerationErr    error
)

func createSvc(rootOpts *options.RootOptions) (*s3.S3, error) {
	return internalaws.CreateAwsService(rootOpts)
}

// Define a mock struct to be used in your unit tests
type mockS3Client struct {
	s3iface.S3API
}

func (m *mockS3Client) GetBucketAccelerateConfiguration(input *s3.GetBucketAccelerateConfigurationInput) (*s3.GetBucketAccelerateConfigurationOutput, error) {
	return defaultGetBucketAccelerationOutput, defaultGetBucketAccelerationErr
}

func (m *mockS3Client) PutBucketAccelerateConfiguration(input *s3.PutBucketAccelerateConfigurationInput) (*s3.PutBucketAccelerateConfigurationOutput, error) {
	return defaultPutBucketAccelerationOutput, defaultPutBucketAccelerationErr
}

func TestExecuteTooManyArguments(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	EnabledCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.S3SvcKey{}, svc))
	EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.OptsKey{}, rootOpts))

	args := []string{"enabled"}
	EnabledCmd.SetArgs(args)

	err = EnabledCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, utils.ErrTooManyArguments, err.Error())

	rootOpts.SetZeroValues()
	transferAccelerationOpts.SetZeroValues()
}

/*func TestExecuteWrongArguments(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	EnabledCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.S3SvcKey{}, svc))
	EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.OptsKey{}, rootOpts))

	args := []string{"eeenabled"}
	EnabledCmd.SetArgs(args)

	err = EnabledCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, ErrWrongArgumentProvided, err.Error())

	rootOpts.SetZeroValues()
	versioningOpts.SetZeroValues()
}*/

/*
	func TestExecuteNoArgument(t *testing.T) {
		rootOpts := options.GetRootOptions()
		rootOpts.AccessKey = "thisisaccesskey"
		rootOpts.SecretKey = "thisissecretkey"
		rootOpts.Region = "thisisregion"
		rootOpts.BucketName = "thisisbucketname"

		ctx := context.Background()
		EnabledCmd.SetContext(ctx)
		svc, err := createSvc(rootOpts)
		assert.NotNil(t, svc)
		assert.Nil(t, err)

		EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.S3SvcKey{}, svc))
		EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.OptsKey{}, rootOpts))

		EnabledCmd.SetArgs([]string{})
		err = EnabledCmd.Execute()
		assert.NotNil(t, err)
		assert.Equal(t, ErrNoArgument, err.Error())

		rootOpts.SetZeroValues()
		versioningOpts.SetZeroValues()
	}
*/
func TestExecuteSuccessAlreadyEnabled(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	EnabledCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketAccelerationErr = nil
	defaultGetBucketAccelerationOutput.Status = aws.String("Enabled")
	defaultPutBucketAccelerationErr = nil

	EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.S3SvcKey{}, svc))
	EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.OptsKey{}, rootOpts))

	EnabledCmd.SetArgs([]string{})
	err := EnabledCmd.Execute()
	assert.Nil(t, err)

	rootOpts.SetZeroValues()
	transferAccelerationOpts.SetZeroValues()
}

func TestExecuteSuccess(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	EnabledCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketAccelerationErr = nil
	defaultGetBucketAccelerationOutput.Status = aws.String("Suspended")
	defaultPutBucketAccelerationErr = nil

	EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.S3SvcKey{}, svc))
	EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.OptsKey{}, rootOpts))

	EnabledCmd.SetArgs([]string{})
	err := EnabledCmd.Execute()
	assert.Nil(t, err)

	rootOpts.SetZeroValues()
	transferAccelerationOpts.SetZeroValues()
}

func TestExecuteGetBucketAccelerationErr(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	EnabledCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketAccelerationErr = errors.New("dummy error")
	defaultPutBucketAccelerationErr = nil

	EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.S3SvcKey{}, svc))
	EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.OptsKey{}, rootOpts))

	EnabledCmd.SetArgs([]string{})
	err := EnabledCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	transferAccelerationOpts.SetZeroValues()
}

func TestExecuteSetBucketAccelerationErr(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	EnabledCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketAccelerationErr = nil
	defaultGetBucketAccelerationOutput.Status = aws.String("Suspended")
	defaultPutBucketAccelerationErr = errors.New("dummy error")

	EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.S3SvcKey{}, svc))
	EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.OptsKey{}, rootOpts))

	EnabledCmd.SetArgs([]string{})
	err := EnabledCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	transferAccelerationOpts.SetZeroValues()
}

func TestExecuteUnknownErr(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	EnabledCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketAccelerationErr = nil
	defaultGetBucketAccelerationOutput.Status = aws.String("Enableddd")
	defaultPutBucketAccelerationErr = errors.New("dummy error")

	EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.S3SvcKey{}, svc))
	EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.OptsKey{}, rootOpts))

	EnabledCmd.SetArgs([]string{})
	err := EnabledCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	transferAccelerationOpts.SetZeroValues()
}
