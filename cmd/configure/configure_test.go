package configure

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/stretchr/testify/assert"
)

var (
	defaultGetBucketVersioningOutput = &s3.GetBucketVersioningOutput{
		MFADelete: aws.String("True"),
		Status:    aws.String("Enabled"),
	}
	defaultGetBucketVersioningErr error = nil
	defaultPutBucketVersioningErr error = nil
)

// Define a mock struct to be used in your unit tests
type mockS3Client struct {
	s3iface.S3API
}

// ListObjects mocks the S3API ListObjects method
func (m *mockS3Client) GetBucketVersioning(obj *s3.GetBucketVersioningInput) (*s3.GetBucketVersioningOutput, error) {
	return defaultGetBucketVersioningOutput, defaultGetBucketVersioningErr
}

func (m *mockS3Client) PutBucketVersioning(obj *s3.PutBucketVersioningInput) (*s3.PutBucketVersioningOutput, error) {
	return &s3.PutBucketVersioningOutput{}, defaultPutBucketVersioningErr
}

func createSvc(rootOpts *options.RootOptions) (*s3.S3, error) {
	return internalaws.CreateAwsService(rootOpts)
}

func TestExecuteVersioningAlreadyEnabled(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"
	rootOpts.Interactive = false

	ctx := context.Background()
	ConfigureCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketVersioningOutput.Status = aws.String("Enabled")
	defaultGetBucketVersioningErr = nil
	configureOpts.Versioning = true

	ConfigureCmd.SetContext(context.WithValue(ConfigureCmd.Context(), options.S3SvcKey{}, mockSvc))
	ConfigureCmd.SetContext(context.WithValue(ConfigureCmd.Context(), options.OptsKey{}, rootOpts))

	err := ConfigureCmd.Execute()
	assert.Nil(t, err)

	rootOpts.SetZeroValues()
	configureOpts.SetZeroValues()
}

func TestExecuteVersioningEnableSuccess(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"
	rootOpts.Interactive = false

	ctx := context.Background()
	ConfigureCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketVersioningOutput.Status = aws.String("Suspended")
	defaultGetBucketVersioningErr = nil
	configureOpts.Versioning = true
	defaultPutBucketVersioningErr = nil

	ConfigureCmd.SetContext(context.WithValue(ConfigureCmd.Context(), options.S3SvcKey{}, mockSvc))
	ConfigureCmd.SetContext(context.WithValue(ConfigureCmd.Context(), options.OptsKey{}, rootOpts))

	err := ConfigureCmd.Execute()
	assert.Nil(t, err)

	rootOpts.SetZeroValues()
	configureOpts.SetZeroValues()
}

func TestExecuteVersioningEnableFailure(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"
	rootOpts.Interactive = false

	ctx := context.Background()
	ConfigureCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketVersioningOutput.Status = aws.String("Suspended")
	defaultGetBucketVersioningErr = nil
	configureOpts.Versioning = true
	defaultPutBucketVersioningErr = errors.New("an error occurred while setting bucket versioning")

	ConfigureCmd.SetContext(context.WithValue(ConfigureCmd.Context(), options.S3SvcKey{}, mockSvc))
	ConfigureCmd.SetContext(context.WithValue(ConfigureCmd.Context(), options.OptsKey{}, rootOpts))

	err := ConfigureCmd.Execute()
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
