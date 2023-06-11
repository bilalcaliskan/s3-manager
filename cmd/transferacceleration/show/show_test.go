package show

import (
	"context"
	"errors"
	"testing"

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
	defaultGetBucketAccelerationErr error
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

func TestExecuteTooManyArguments(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	ShowCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3SvcKey{}, svc))
	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))

	args := []string{"enabled"}
	ShowCmd.SetArgs(args)

	err = ShowCmd.Execute()
	assert.NotNil(t, err)

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
	ShowCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3SvcKey{}, svc))
	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))

	args := []string{"eeenabled"}
	ShowCmd.SetArgs(args)

	err = ShowCmd.Execute()
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
		ShowCmd.SetContext(ctx)
		svc, err := createSvc(rootOpts)
		assert.NotNil(t, svc)
		assert.Nil(t, err)

		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3SvcKey{}, svc))
		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))

		ShowCmd.SetArgs([]string{})
		err = ShowCmd.Execute()
		assert.NotNil(t, err)
		assert.Equal(t, ErrNoArgument, err.Error())

		rootOpts.SetZeroValues()
		versioningOpts.SetZeroValues()
	}
*/
func TestExecuteSuccessAlreadyDisabled(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	ShowCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketAccelerationErr = nil
	defaultGetBucketAccelerationOutput.Status = aws.String("Suspended")

	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3SvcKey{}, svc))
	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))

	ShowCmd.SetArgs([]string{})
	err := ShowCmd.Execute()
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
	ShowCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketAccelerationErr = nil
	defaultGetBucketAccelerationOutput.Status = aws.String("Enabled")

	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3SvcKey{}, svc))
	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))

	ShowCmd.SetArgs([]string{})
	err := ShowCmd.Execute()
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
	ShowCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketAccelerationErr = errors.New("dummy error")
	defaultGetBucketAccelerationOutput.Status = aws.String("Suspended")

	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3SvcKey{}, svc))
	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))

	ShowCmd.SetArgs([]string{})
	err := ShowCmd.Execute()
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
	ShowCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketAccelerationErr = nil
	defaultGetBucketAccelerationOutput.Status = aws.String("Enableddd")

	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3SvcKey{}, svc))
	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))

	ShowCmd.SetArgs([]string{})
	err := ShowCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	transferAccelerationOpts.SetZeroValues()
}

/*
func TestExecuteSuccessEnabledWrongVersioning(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	ShowCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketAccelerationErr = nil
	defaultGetBucketAccelerationOutput.Status = aws.String("Suspended")
	defaultPutBucketAccelerationErr = nil

	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3SvcKey{}, svc))
	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))

	ShowCmd.SetArgs([]string{})
	err := ShowCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	transferAccelerationOpts.SetZeroValues()
}
*/
