package disabled

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
	DisabledCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.S3SvcKey{}, svc))
	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.OptsKey{}, rootOpts))

	args := []string{"enabled"}
	DisabledCmd.SetArgs(args)

	err = DisabledCmd.Execute()
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
	DisabledCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.S3SvcKey{}, svc))
	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.OptsKey{}, rootOpts))

	args := []string{"eeenabled"}
	DisabledCmd.SetArgs(args)

	err = DisabledCmd.Execute()
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
		DisabledCmd.SetContext(ctx)
		svc, err := createSvc(rootOpts)
		assert.NotNil(t, svc)
		assert.Nil(t, err)

		DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.S3SvcKey{}, svc))
		DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.OptsKey{}, rootOpts))

		DisabledCmd.SetArgs([]string{})
		err = DisabledCmd.Execute()
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
	DisabledCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketAccelerationErr = nil
	defaultGetBucketAccelerationOutput.Status = aws.String("Suspended")
	defaultPutBucketAccelerationErr = nil

	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.S3SvcKey{}, svc))
	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.OptsKey{}, rootOpts))

	DisabledCmd.SetArgs([]string{})
	err := DisabledCmd.Execute()
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
	DisabledCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketAccelerationErr = nil
	defaultGetBucketAccelerationOutput.Status = aws.String("Enabled")
	defaultPutBucketAccelerationErr = nil

	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.S3SvcKey{}, svc))
	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.OptsKey{}, rootOpts))

	DisabledCmd.SetArgs([]string{})
	err := DisabledCmd.Execute()
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
	DisabledCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketAccelerationErr = errors.New("dummy error")
	defaultGetBucketAccelerationOutput.Status = aws.String("Suspended")
	defaultPutBucketAccelerationErr = nil

	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.S3SvcKey{}, svc))
	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.OptsKey{}, rootOpts))

	DisabledCmd.SetArgs([]string{})
	err := DisabledCmd.Execute()
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
	DisabledCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketAccelerationErr = nil
	defaultGetBucketAccelerationOutput.Status = aws.String("Enabled")
	defaultPutBucketAccelerationErr = errors.New("dummy error")

	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.S3SvcKey{}, svc))
	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.OptsKey{}, rootOpts))

	DisabledCmd.SetArgs([]string{})
	err := DisabledCmd.Execute()
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
	DisabledCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketAccelerationErr = nil
	defaultGetBucketAccelerationOutput.Status = aws.String("Enableddd")
	defaultPutBucketAccelerationErr = errors.New("dummy error")

	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.S3SvcKey{}, svc))
	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.OptsKey{}, rootOpts))

	DisabledCmd.SetArgs([]string{})
	err := DisabledCmd.Execute()
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
	DisabledCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketAccelerationErr = nil
	defaultGetBucketAccelerationOutput.Status = aws.String("Suspended")
	defaultPutBucketAccelerationErr = nil

	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.S3SvcKey{}, svc))
	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.OptsKey{}, rootOpts))

	DisabledCmd.SetArgs([]string{})
	err := DisabledCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	transferAccelerationOpts.SetZeroValues()
}
*/
