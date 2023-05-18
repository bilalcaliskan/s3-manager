package versioning

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/aws/aws-sdk-go/service/s3"
	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
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
	//err = CleanCmd.PersistentFlags().Set("verbose", "true")
	//assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	VersioningCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	VersioningCmd.SetContext(context.WithValue(VersioningCmd.Context(), options.S3SvcKey{}, svc))
	VersioningCmd.SetContext(context.WithValue(VersioningCmd.Context(), options.OptsKey{}, rootOpts))

	args := []string{"enabled", "foo"}
	VersioningCmd.SetArgs(args)

	err = VersioningCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, ErrTooManyArguments, err.Error())

	rootOpts.SetZeroValues()
	versioningOpts.SetZeroValues()
}

func TestExecuteWrongArguments(t *testing.T) {
	//err = CleanCmd.PersistentFlags().Set("verbose", "true")
	//assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	VersioningCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	VersioningCmd.SetContext(context.WithValue(VersioningCmd.Context(), options.S3SvcKey{}, svc))
	VersioningCmd.SetContext(context.WithValue(VersioningCmd.Context(), options.OptsKey{}, rootOpts))

	args := []string{"eeenabled"}
	VersioningCmd.SetArgs(args)

	err = VersioningCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, ErrWrongArgumentProvided, err.Error())

	rootOpts.SetZeroValues()
	versioningOpts.SetZeroValues()
}

func TestExecuteNoArgument(t *testing.T) {
	//err = CleanCmd.PersistentFlags().Set("verbose", "true")
	//assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	VersioningCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	VersioningCmd.SetContext(context.WithValue(VersioningCmd.Context(), options.S3SvcKey{}, svc))
	VersioningCmd.SetContext(context.WithValue(VersioningCmd.Context(), options.OptsKey{}, rootOpts))

	VersioningCmd.SetArgs([]string{})
	err = VersioningCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, ErrNoArgument, err.Error())

	rootOpts.SetZeroValues()
	versioningOpts.SetZeroValues()
}

func TestExecuteSuccessEnabled(t *testing.T) {
	//err = CleanCmd.PersistentFlags().Set("verbose", "true")
	//assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	VersioningCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketVersioningErr = nil
	defaultGetBucketVersioningOutput.Status = aws.String("Suspended")
	defaultPutBucketVersioningErr = nil

	VersioningCmd.SetContext(context.WithValue(VersioningCmd.Context(), options.S3SvcKey{}, svc))
	VersioningCmd.SetContext(context.WithValue(VersioningCmd.Context(), options.OptsKey{}, rootOpts))

	VersioningCmd.SetArgs([]string{"enabled"})
	err := VersioningCmd.Execute()
	assert.Nil(t, err)

	rootOpts.SetZeroValues()
	versioningOpts.SetZeroValues()
}

func TestExecuteSuccessAlreadyEnabled(t *testing.T) {
	//err = CleanCmd.PersistentFlags().Set("verbose", "true")
	//assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	VersioningCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketVersioningErr = nil
	defaultGetBucketVersioningOutput.Status = aws.String("Enabled")
	defaultPutBucketVersioningErr = nil

	VersioningCmd.SetContext(context.WithValue(VersioningCmd.Context(), options.S3SvcKey{}, svc))
	VersioningCmd.SetContext(context.WithValue(VersioningCmd.Context(), options.OptsKey{}, rootOpts))

	VersioningCmd.SetArgs([]string{"enabled"})
	err := VersioningCmd.Execute()
	assert.Nil(t, err)

	rootOpts.SetZeroValues()
	versioningOpts.SetZeroValues()
}

func TestExecuteGetBucketVersioningErr(t *testing.T) {
	//err = CleanCmd.PersistentFlags().Set("verbose", "true")
	//assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	VersioningCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketVersioningErr = errors.New("dummy error")

	VersioningCmd.SetContext(context.WithValue(VersioningCmd.Context(), options.S3SvcKey{}, svc))
	VersioningCmd.SetContext(context.WithValue(VersioningCmd.Context(), options.OptsKey{}, rootOpts))

	VersioningCmd.SetArgs([]string{"enabled"})
	err := VersioningCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	versioningOpts.SetZeroValues()
}

func TestExecuteSetBucketVersioningErr(t *testing.T) {
	//err = CleanCmd.PersistentFlags().Set("verbose", "true")
	//assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	VersioningCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketVersioningErr = nil
	defaultPutBucketVersioningErr = errors.New("new dummy error")

	VersioningCmd.SetContext(context.WithValue(VersioningCmd.Context(), options.S3SvcKey{}, svc))
	VersioningCmd.SetContext(context.WithValue(VersioningCmd.Context(), options.OptsKey{}, rootOpts))

	VersioningCmd.SetArgs([]string{"disabled"})
	err := VersioningCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	versioningOpts.SetZeroValues()
}

/*func TestExecuteSuccess(t *testing.T) {
	//err = CleanCmd.PersistentFlags().Set("verbose", "true")
	//assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	VersioningCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultListObjectsOutput.Contents = []*s3.Object{
		{
			ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
			Key:          aws.String("../../mock/file1.txt"),
			StorageClass: aws.String("STANDARD"),
		},
		{
			ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
			Key:          aws.String("../../mock/file2.txt"),
			StorageClass: aws.String("STANDARD"),
		},
		{
			ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
			Key:          aws.String("../../mock/file3.txt"),
			StorageClass: aws.String("STANDARD"),
		},
	}

	SearchCmd.SetContext(context.WithValue(SearchCmd.Context(), options.S3SvcKey{}, mockSvc))
	SearchCmd.SetContext(context.WithValue(SearchCmd.Context(), options.OptsKey{}, rootOpts))

	searchOpts.Substring = "akqASmLLlK"
	err := SearchCmd.Execute()
	assert.Nil(t, err)

	rootOpts.SetZeroValues()
	searchOpts.SetZeroValues()
}

func TestExecuteEmptyList(t *testing.T) {
	//err = CleanCmd.PersistentFlags().Set("verbose", "true")
	//assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	SearchCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultListObjectsOutput.Contents = []*s3.Object{}

	SearchCmd.SetContext(context.WithValue(SearchCmd.Context(), options.S3SvcKey{}, mockSvc))
	SearchCmd.SetContext(context.WithValue(SearchCmd.Context(), options.OptsKey{}, rootOpts))

	searchOpts.Substring = "akqASmLLlK"
	err := SearchCmd.Execute()
	assert.Nil(t, err)

	rootOpts.SetZeroValues()
	searchOpts.SetZeroValues()
}
*/
