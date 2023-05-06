package aws

import (
	"errors"
	"os"
	"testing"
	"time"

	options2 "github.com/bilalcaliskan/s3-manager/cmd/search/options"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/internal/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

var (
	listObjectsErr           error
	getObjectsErr            error
	deleteObjectsErr         error
	fileNamePrefix           string
	defaultListObjectsOutput = &s3.ListObjectsOutput{
		Name:        aws.String(""),
		Marker:      aws.String(""),
		MaxKeys:     aws.Int64(1000),
		Prefix:      aws.String(""),
		IsTruncated: aws.Bool(false),
	}
	defaultDeleteObjectOutput = &s3.DeleteObjectOutput{
		DeleteMarker:   nil,
		RequestCharged: nil,
		VersionId:      nil,
	}
	mockLogger = logging.GetLogger(options.GetRootOptions())
)

type mockS3Client struct {
	s3iface.S3API
}

// ListObjects mocks the S3API ListObjects method
func (m *mockS3Client) ListObjects(obj *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	return defaultListObjectsOutput, listObjectsErr
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
		ContentType:   aws.String("text/plain"),
		ETag:          aws.String("d73a503d212d9279e6b2ed8ac6bb81f3"),
	}, getObjectsErr
}

func (m *mockS3Client) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return defaultDeleteObjectOutput, deleteObjectsErr
}

func TestGetAllFilesHappyPath(t *testing.T) {
	m := &mockS3Client{}
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

	result, err := GetAllFiles(m, options.GetRootOptions(), fileNamePrefix)
	assert.NotEmpty(t, result)
	assert.Nil(t, err)
}

func TestGetAllFilesFailedListObjectsCall(t *testing.T) {
	m := &mockS3Client{}
	listObjectsErr = errors.New("dummy error thrown")
	_, err := GetAllFiles(m, options.GetRootOptions(), fileNamePrefix)
	assert.NotNil(t, err)
	listObjectsErr = nil
}

func TestDeleteFilesHappyPath(t *testing.T) {
	var input []*s3.Object
	m := &mockS3Client{}
	deleteObjectsErr = nil

	err := DeleteFiles(m, "dummy bucket", input, false, mockLogger)
	assert.Nil(t, err)
}

func TestDeleteFilesHappyPathDryRun(t *testing.T) {
	var input []*s3.Object
	m := &mockS3Client{}
	deleteObjectsErr = nil

	err := DeleteFiles(m, "dummy bucket", input, true, mockLogger)
	assert.Nil(t, err)
}

func TestDeleteFilesFailedDeleteObjectCall(t *testing.T) {
	var input []*s3.Object
	for i := 0; i < 3; i++ {
		o := s3.Object{Key: aws.String("hello-world"), LastModified: aws.Time(time.Now()), Size: aws.Int64(10000000)}
		input = append(input, &o)
	}

	m := &mockS3Client{}
	deleteObjectsErr = errors.New("dummy error")
	err := DeleteFiles(m, "dummy bucket", input, false, mockLogger)
	assert.NotNil(t, err)
	deleteObjectsErr = nil
}

func TestCreateAwsService(t *testing.T) {
	opts := options.GetRootOptions()
	opts.AccessKey = "thisisaccesskey"
	opts.SecretKey = "thisissecretkey"
	opts.Region = "thisisregion"
	opts.BucketName = "thisisbucketname"

	svc, err := CreateAwsService(opts)
	assert.Nil(t, err)
	assert.NotNil(t, svc)
}

func TestFind(t *testing.T) {
	mockSvc := &mockS3Client{}
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
	searchOpts := options2.GetSearchOptions()
	rootOpts := options.GetRootOptions()
	searchOpts.RootOptions = rootOpts

	searchOpts.Substring = "akqASmLLlK"

	result, errs := Find(mockSvc, searchOpts, logging.GetLogger(options.GetRootOptions()))
	assert.NotNil(t, result)
	assert.Empty(t, errs)
}

func TestFindWrongFilePath(t *testing.T) {
	mockSvc := &mockS3Client{}
	defaultListObjectsOutput.Contents = []*s3.Object{
		{
			ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
			Key:          aws.String("file1asdfasdf.txt"),
			StorageClass: aws.String("STANDARD"),
		},
	}
	searchOpts := options2.GetSearchOptions()
	rootOpts := options.GetRootOptions()
	searchOpts.RootOptions = rootOpts

	res, err := Find(mockSvc, searchOpts, logging.GetLogger(options.GetRootOptions()))
	assert.Nil(t, res)
	assert.NotEmpty(t, err)
}
