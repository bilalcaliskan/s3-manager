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
	defaultGetBucketVersioningOutput = &s3.GetBucketVersioningOutput{
		Status: aws.String("Enabled"),
	}
	defaultGetBucketVersioningErr    error
	defaultPutBucketVersioningOutput = &s3.PutBucketVersioningOutput{}
	defaultPutBucketVersioningErr    error
)

func createSvc(rootOpts *options.RootOptions) (*s3.S3, error) {
	return internalaws.CreateAwsService(rootOpts)
}

// Define a mock struct to be used in your unit tests
type mockS3Client struct {
	s3iface.S3API
}

func (m *mockS3Client) GetBucketVersioning(input *s3.GetBucketVersioningInput) (*s3.GetBucketVersioningOutput, error) {
	return defaultGetBucketVersioningOutput, defaultGetBucketVersioningErr
}

func (m *mockS3Client) PutBucketVersioning(input *s3.PutBucketVersioningInput) (*s3.PutBucketVersioningOutput, error) {
	return defaultPutBucketVersioningOutput, defaultPutBucketVersioningErr
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
	versioningOpts.SetZeroValues()
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

	defaultGetBucketVersioningErr = nil
	defaultGetBucketVersioningOutput.Status = aws.String("Suspended")
	defaultPutBucketVersioningErr = nil

	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.S3SvcKey{}, svc))
	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.OptsKey{}, rootOpts))

	DisabledCmd.SetArgs([]string{})
	err := DisabledCmd.Execute()
	assert.Nil(t, err)

	rootOpts.SetZeroValues()
	versioningOpts.SetZeroValues()
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

	defaultGetBucketVersioningErr = nil
	defaultGetBucketVersioningOutput.Status = aws.String("Enabled")
	defaultPutBucketVersioningErr = nil

	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.S3SvcKey{}, svc))
	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.OptsKey{}, rootOpts))

	DisabledCmd.SetArgs([]string{})
	err := DisabledCmd.Execute()
	assert.Nil(t, err)

	rootOpts.SetZeroValues()
	versioningOpts.SetZeroValues()
}

func TestExecuteGetBucketVersioningErr(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	DisabledCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketVersioningErr = errors.New("dummy error")

	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.S3SvcKey{}, svc))
	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.OptsKey{}, rootOpts))

	DisabledCmd.SetArgs([]string{})
	err := DisabledCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	versioningOpts.SetZeroValues()
}

func TestExecuteSetBucketVersioningErr(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	DisabledCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketVersioningErr = nil
	defaultPutBucketVersioningErr = errors.New("new dummy error")

	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.S3SvcKey{}, svc))
	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.OptsKey{}, rootOpts))

	DisabledCmd.SetArgs([]string{})
	err := DisabledCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	versioningOpts.SetZeroValues()
}

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

	defaultGetBucketVersioningErr = nil
	defaultGetBucketVersioningOutput.Status = aws.String("Suspendeddd")
	defaultPutBucketVersioningErr = nil

	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.S3SvcKey{}, svc))
	DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.OptsKey{}, rootOpts))

	DisabledCmd.SetArgs([]string{})
	err := DisabledCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	versioningOpts.SetZeroValues()
}
