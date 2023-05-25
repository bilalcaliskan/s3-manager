package substring

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	options2 "github.com/bilalcaliskan/s3-manager/cmd/search/options"
	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	defaultListObjectsOutput = &s3.ListObjectsOutput{
		Name:        aws.String(""),
		Marker:      aws.String(""),
		MaxKeys:     aws.Int64(1000),
		Prefix:      aws.String(""),
		IsTruncated: aws.Bool(false),
	}
)

func createSvc(rootOpts *options.RootOptions) (*s3.S3, error) {
	return internalaws.CreateAwsService(rootOpts)
}

// Define a mock struct to be used in your unit tests
type mockS3Client struct {
	s3iface.S3API
}

// ListObjects mocks the S3API ListObjects method
func (m *mockS3Client) ListObjects(obj *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	return defaultListObjectsOutput, nil
}

// GetObject mocks the S3API GetObject method
func (m *mockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	bytes, err := os.Open(*input.Key)
	if err != nil {
		return nil, err
	}

	return &s3.GetObjectOutput{
		AcceptRanges:  aws.String("bytes"),
		Body:          bytes,
		ContentLength: aws.Int64(1000),
		ContentType:   aws.String("substring/plain"),
		ETag:          aws.String("d73a503d212d9279e6b2ed8ac6bb81f3"),
	}, nil
}

func TestExecuteMissingRegion(t *testing.T) {
	//err = CleanCmd.PersistentFlags().Set("verbose", "true")
	//assert.Nil(t, err)

	rootOpts := options.GetRootOptions()

	searchOpts = options2.GetSearchOptions()
	searchOpts.RootOptions = rootOpts

	ctx := context.Background()
	SubstringCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.Nil(t, svc)
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	searchOpts.SetZeroValues()
}

func TestExecuteFailure(t *testing.T) {
	//err = CleanCmd.PersistentFlags().Set("verbose", "true")
	//assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	searchOpts = options2.GetSearchOptions()
	searchOpts.RootOptions = rootOpts

	ctx := context.Background()
	SubstringCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	SubstringCmd.SetContext(context.WithValue(SubstringCmd.Context(), options.S3SvcKey{}, svc))
	SubstringCmd.SetContext(context.WithValue(SubstringCmd.Context(), options.OptsKey{}, rootOpts))

	err = SubstringCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	searchOpts.SetZeroValues()
}

func TestExecuteFailure2(t *testing.T) {
	//err = CleanCmd.PersistentFlags().Set("verbose", "true")
	//assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	searchOpts = options2.GetSearchOptions()
	searchOpts.RootOptions = rootOpts

	ctx := context.Background()
	SubstringCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	SubstringCmd.SetContext(context.WithValue(SubstringCmd.Context(), options.S3SvcKey{}, svc))
	SubstringCmd.SetContext(context.WithValue(SubstringCmd.Context(), options.OptsKey{}, rootOpts))

	SubstringCmd.SetArgs([]string{"aslkdads", "asdadsadsfdafs"})
	err = SubstringCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	searchOpts.SetZeroValues()
}

func TestExecuteSuccess(t *testing.T) {
	//err = CleanCmd.PersistentFlags().Set("verbose", "true")
	//assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	searchOpts = options2.GetSearchOptions()
	searchOpts.RootOptions = rootOpts

	ctx := context.Background()
	SubstringCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultListObjectsOutput.Contents = []*s3.Object{
		{
			ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
			Key:          aws.String("../../../mock/file1.txt"),
			StorageClass: aws.String("STANDARD"),
		},
		{
			ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
			Key:          aws.String("../../../mock/file2.txt"),
			StorageClass: aws.String("STANDARD"),
		},
		{
			ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
			Key:          aws.String("../../../mock/file3.txt"),
			StorageClass: aws.String("STANDARD"),
		},
	}

	SubstringCmd.SetContext(context.WithValue(SubstringCmd.Context(), options.S3SvcKey{}, mockSvc))
	SubstringCmd.SetContext(context.WithValue(SubstringCmd.Context(), options.OptsKey{}, rootOpts))

	SubstringCmd.SetArgs([]string{"aslkdads"})
	err := SubstringCmd.Execute()
	assert.Nil(t, err)

	rootOpts.SetZeroValues()
	searchOpts.SetZeroValues()
}

func TestExecuteSuccess2(t *testing.T) {
	//err = CleanCmd.PersistentFlags().Set("verbose", "true")
	//assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	searchOpts = options2.GetSearchOptions()
	searchOpts.RootOptions = rootOpts

	ctx := context.Background()
	SubstringCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultListObjectsOutput.Contents = []*s3.Object{
		{
			ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
			Key:          aws.String("../../../mock/file1.txt"),
			StorageClass: aws.String("STANDARD"),
		},
		{
			ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
			Key:          aws.String("../../../mock/file2.txt"),
			StorageClass: aws.String("STANDARD"),
		},
		{
			ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
			Key:          aws.String("../../../mock/file3.txt"),
			StorageClass: aws.String("STANDARD"),
		},
	}

	SubstringCmd.SetContext(context.WithValue(SubstringCmd.Context(), options.S3SvcKey{}, mockSvc))
	SubstringCmd.SetContext(context.WithValue(SubstringCmd.Context(), options.OptsKey{}, rootOpts))

	SubstringCmd.SetArgs([]string{"yILlXDYWyU"})
	err := SubstringCmd.Execute()
	assert.Nil(t, err)

	rootOpts.SetZeroValues()
	searchOpts.SetZeroValues()
}

func TestExecuteFailureMultipleErrors(t *testing.T) {
	//err = CleanCmd.PersistentFlags().Set("verbose", "true")
	//assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	searchOpts = options2.GetSearchOptions()
	searchOpts.RootOptions = rootOpts

	ctx := context.Background()
	SubstringCmd.SetContext(ctx)

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
			Key:          aws.String("../../../mock/file2.txt"),
			StorageClass: aws.String("STANDARD"),
		},
		{
			ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
			Key:          aws.String("../../../mock/file3.txt"),
			StorageClass: aws.String("STANDARD"),
		},
	}

	SubstringCmd.SetContext(context.WithValue(SubstringCmd.Context(), options.S3SvcKey{}, mockSvc))
	SubstringCmd.SetContext(context.WithValue(SubstringCmd.Context(), options.OptsKey{}, rootOpts))

	SubstringCmd.SetArgs([]string{"aslkdads"})
	err := SubstringCmd.Execute()
	assert.NotNil(t, err)

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

	searchOpts = options2.GetSearchOptions()
	searchOpts.RootOptions = rootOpts

	ctx := context.Background()
	SubstringCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultListObjectsOutput.Contents = []*s3.Object{}

	SubstringCmd.SetContext(context.WithValue(SubstringCmd.Context(), options.S3SvcKey{}, mockSvc))
	SubstringCmd.SetContext(context.WithValue(SubstringCmd.Context(), options.OptsKey{}, rootOpts))

	SubstringCmd.SetArgs([]string{"aslkdads"})
	err := SubstringCmd.Execute()
	assert.Nil(t, err)

	rootOpts.SetZeroValues()
	searchOpts.SetZeroValues()
}
