package aws

import (
	"errors"
	"os"
	"testing"
	"time"

	options4 "github.com/bilalcaliskan/s3-manager/cmd/tags/options"

	options3 "github.com/bilalcaliskan/s3-manager/cmd/versioning/options"

	options2 "github.com/bilalcaliskan/s3-manager/cmd/search/options"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/internal/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

var (
	defaultListObjectsErr    error
	defaultGetObjectErr      error
	defaultDeleteObjectErr   error
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
	mockLogger                       = logging.GetLogger(options.GetRootOptions())
	defaultGetBucketVersioningOutput = &s3.GetBucketVersioningOutput{
		Status: aws.String("Enabled"),
	}
	defaultGetBucketVersioningErr    error
	defaultPutBucketVersioningOutput = &s3.PutBucketVersioningOutput{}
	defaultPutBucketVersioningErr    error
	defaultGetBucketTaggingErr       error
	defaultGetBucketTaggingOutput    = &s3.GetBucketTaggingOutput{}
	defaultPutBucketTaggingErr       error
	defaultPutBucketTaggingOutput    = &s3.PutBucketTaggingOutput{}
	defaultDeleteBucketTaggingErr    error
	defaultDeleteBucketTaggingOutput = &s3.DeleteBucketTaggingOutput{}
)

type mockS3Client struct {
	s3iface.S3API
}

// ListObjects mocks the S3API ListObjects method
func (m *mockS3Client) ListObjects(obj *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	return defaultListObjectsOutput, defaultListObjectsErr
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
	}, defaultGetObjectErr
}

func (m *mockS3Client) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return defaultDeleteObjectOutput, defaultDeleteObjectErr
}

func (m *mockS3Client) GetBucketVersioning(input *s3.GetBucketVersioningInput) (*s3.GetBucketVersioningOutput, error) {
	return defaultGetBucketVersioningOutput, defaultGetBucketVersioningErr
}

func (m *mockS3Client) PutBucketVersioning(input *s3.PutBucketVersioningInput) (*s3.PutBucketVersioningOutput, error) {
	return defaultPutBucketVersioningOutput, defaultPutBucketVersioningErr
}

func (m *mockS3Client) GetBucketTagging(input *s3.GetBucketTaggingInput) (*s3.GetBucketTaggingOutput, error) {
	return defaultGetBucketTaggingOutput, defaultGetBucketTaggingErr
}

func (m *mockS3Client) PutBucketTagging(input *s3.PutBucketTaggingInput) (*s3.PutBucketTaggingOutput, error) {
	return defaultPutBucketTaggingOutput, defaultPutBucketTaggingErr
}

func (m *mockS3Client) DeleteBucketTagging(input *s3.DeleteBucketTaggingInput) (*s3.DeleteBucketTaggingOutput, error) {
	return defaultDeleteBucketTaggingOutput, defaultDeleteBucketTaggingErr
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
	defaultListObjectsErr = errors.New("dummy error thrown")
	_, err := GetAllFiles(m, options.GetRootOptions(), fileNamePrefix)
	assert.NotNil(t, err)
}

func TestDeleteFilesHappyPath(t *testing.T) {
	var input []*s3.Object
	m := &mockS3Client{}
	defaultDeleteObjectErr = nil

	err := DeleteFiles(m, "dummy bucket", input, false, mockLogger)
	assert.Nil(t, err)
}

func TestDeleteFilesHappyPathDryRun(t *testing.T) {
	var input []*s3.Object
	m := &mockS3Client{}
	defaultDeleteObjectErr = nil

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
	defaultDeleteObjectErr = errors.New("dummy error")
	err := DeleteFiles(m, "dummy bucket", input, false, mockLogger)
	assert.NotNil(t, err)
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

func TestSearchStringSuccess(t *testing.T) {
	mockSvc := &mockS3Client{}
	defaultListObjectsErr = nil
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

	result, errs := SearchString(mockSvc, searchOpts, logging.GetLogger(options.GetRootOptions()))
	assert.NotNil(t, result)
	assert.Empty(t, errs)
}

func TestSearchStringFailure(t *testing.T) {
	mockSvc := &mockS3Client{}
	defaultListObjectsErr = errors.New("dummy error")
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

	result, errs := SearchString(mockSvc, searchOpts, logging.GetLogger(options.GetRootOptions()))
	assert.Nil(t, result)
	assert.NotEmpty(t, errs)
}

func TestSearchStringGetObjectFailure(t *testing.T) {
	mockSvc := &mockS3Client{}
	defaultListObjectsErr = nil
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
	defaultGetObjectErr = errors.New("dummy error")
	searchOpts := options2.GetSearchOptions()
	rootOpts := options.GetRootOptions()
	searchOpts.RootOptions = rootOpts

	searchOpts.Substring = "akqASmLLlK"

	result, errs := SearchString(mockSvc, searchOpts, logging.GetLogger(options.GetRootOptions()))
	assert.Nil(t, result)
	assert.NotEmpty(t, errs)
}

func TestSearchStringWrongFilePath(t *testing.T) {
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

	res, err := SearchString(mockSvc, searchOpts, logging.GetLogger(options.GetRootOptions()))
	assert.NotNil(t, err)
	assert.Empty(t, res)
}

func TestGetDesiredFilesSuccess(t *testing.T) {
	mockSvc := &mockS3Client{}
	defaultListObjectsErr = nil
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
	searchOpts.FileName = "file1.*"

	res, err := GetDesiredFiles(mockSvc, searchOpts)
	assert.Nil(t, err)
	assert.NotEmpty(t, res)
}

func TestGetDesiredFilesFailure(t *testing.T) {
	mockSvc := &mockS3Client{}
	defaultListObjectsErr = errors.New("dummy error")
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
	searchOpts.FileName = "file1.*"

	res, err := GetDesiredFiles(mockSvc, searchOpts)
	assert.NotNil(t, err)
	assert.Empty(t, res)
}

func TestSetBucketVersioningSuccessEnabled(t *testing.T) {
	mockSvc := &mockS3Client{}
	defaultPutBucketVersioningErr = nil
	versioningOpts := &options3.VersioningOptions{
		DesiredState: "enabled",
		RootOptions: &options.RootOptions{
			BucketName: "demo-bucket",
		},
	}

	assert.Nil(t, SetBucketVersioning(mockSvc, versioningOpts, logging.GetLogger(versioningOpts.RootOptions)))
}

func TestSetBucketVersioningSuccessDisabled(t *testing.T) {
	mockSvc := &mockS3Client{}
	defaultPutBucketVersioningErr = nil
	versioningOpts := &options3.VersioningOptions{
		DesiredState: "disabled",
		RootOptions: &options.RootOptions{
			BucketName: "demo-bucket",
		},
	}

	assert.Nil(t, SetBucketVersioning(mockSvc, versioningOpts, logging.GetLogger(versioningOpts.RootOptions)))
}

func TestSetBucketVersioningError(t *testing.T) {
	mockSvc := &mockS3Client{}
	defaultGetBucketVersioningOutput = &s3.GetBucketVersioningOutput{Status: aws.String("disabledd")}
	defaultPutBucketVersioningErr = errors.New("asdflkjasdf")
	versioningOpts := &options3.VersioningOptions{
		DesiredState: "enabled",
		RootOptions: &options.RootOptions{
			BucketName: "demo-bucket",
		},
	}
	versioningOpts.ActualState = "disableddd"

	assert.NotNil(t, SetBucketVersioning(mockSvc, versioningOpts, logging.GetLogger(versioningOpts.RootOptions)))
}

func TestGetBucketVersioningSuccess(t *testing.T) {
	mockSvc := &mockS3Client{}
	defaultGetBucketVersioningErr = nil
	versioningOpts := &options3.VersioningOptions{
		DesiredState: "enabled",
		RootOptions: &options.RootOptions{
			BucketName: "demo-bucket",
		},
	}

	_, err := GetBucketVersioning(mockSvc, versioningOpts.RootOptions)
	assert.Nil(t, err)
}

func TestGetBucketVersioningFailure(t *testing.T) {
	mockSvc := &mockS3Client{}
	defaultGetBucketVersioningErr = errors.New("adsfafdsadsf")
	versioningOpts := &options3.VersioningOptions{
		DesiredState: "enabled",
		RootOptions: &options.RootOptions{
			BucketName: "demo-bucket",
		},
	}

	_, err := GetBucketVersioning(mockSvc, versioningOpts.RootOptions)
	assert.NotNil(t, err)
}

func TestCreateAwsServiceErr(t *testing.T) {
	opts := &options.RootOptions{
		AccessKey:  "asdfadsf",
		SecretKey:  "asdfsadf",
		BucketName: "asdfasdf",
		Region:     "",
	}

	svc, err := CreateAwsService(opts)
	assert.Nil(t, svc)
	assert.NotNil(t, err)
}

func TestGetBucketTaggingSuccess(t *testing.T) {
	tagOpts := options4.GetTagOptions()
	defer func() {
		tagOpts.SetZeroValues()
	}()
	rootOpts := options.GetRootOptions()
	rootOpts.Region = "us-east-1"
	tagOpts.RootOptions = rootOpts

	mockSvc := &mockS3Client{}
	var tags []*s3.Tag
	tags = append(tags, &s3.Tag{Key: aws.String("foo"), Value: aws.String("bar")})
	defaultGetBucketTaggingOutput = &s3.GetBucketTaggingOutput{TagSet: tags}
	defaultGetBucketTaggingErr = nil

	_, err := GetBucketTags(mockSvc, tagOpts)
	assert.Nil(t, err)
}

func TestGetBucketTaggingFailure(t *testing.T) {
	tagOpts := options4.GetTagOptions()
	defer func() {
		tagOpts.SetZeroValues()
	}()
	rootOpts := options.GetRootOptions()
	rootOpts.Region = "us-east-1"
	tagOpts.RootOptions = rootOpts

	mockSvc := &mockS3Client{}
	var tags []*s3.Tag
	defaultGetBucketTaggingOutput = &s3.GetBucketTaggingOutput{TagSet: tags}
	defaultGetBucketTaggingErr = errors.New("dummy error")

	_, err := GetBucketTags(mockSvc, tagOpts)
	assert.NotNil(t, err)
}

func TestPutBucketTaggingSuccess(t *testing.T) {
	tagOpts := options4.GetTagOptions()
	defer func() {
		tagOpts.SetZeroValues()
	}()
	rootOpts := options.GetRootOptions()
	rootOpts.Region = "us-east-1"
	tagOpts.RootOptions = rootOpts

	mockSvc := &mockS3Client{}
	var tags []*s3.Tag
	tags = append(tags, &s3.Tag{Key: aws.String("foo"), Value: aws.String("bar")})
	for _, v := range tags {
		tagOpts.TagsToAdd[*v.Key] = *v.Value
	}

	defaultGetBucketTaggingOutput = &s3.GetBucketTaggingOutput{TagSet: tags}
	defaultGetBucketTaggingErr = nil

	defaultPutBucketTaggingErr = nil

	_, err := SetBucketTags(mockSvc, tagOpts)
	assert.Nil(t, err)
}

func TestPutBucketTaggingFailure(t *testing.T) {
	tagOpts := options4.GetTagOptions()
	defer func() {
		tagOpts.SetZeroValues()
	}()
	rootOpts := options.GetRootOptions()
	rootOpts.Region = "us-east-1"
	tagOpts.RootOptions = rootOpts

	mockSvc := &mockS3Client{}
	var tags []*s3.Tag
	tags = append(tags, &s3.Tag{Key: aws.String("foo"), Value: aws.String("bar")})
	for _, v := range tags {
		tagOpts.TagsToAdd[*v.Key] = *v.Value
	}

	defaultGetBucketTaggingOutput = &s3.GetBucketTaggingOutput{TagSet: tags}
	defaultGetBucketTaggingErr = nil

	defaultPutBucketTaggingErr = errors.New("dummy error")

	_, err := SetBucketTags(mockSvc, tagOpts)
	assert.NotNil(t, err)
}

func TestDeleteBucketTaggingSuccess(t *testing.T) {
	tagOpts := options4.GetTagOptions()
	defer func() {
		tagOpts.SetZeroValues()
	}()
	rootOpts := options.GetRootOptions()
	rootOpts.Region = "us-east-1"
	tagOpts.RootOptions = rootOpts

	defaultDeleteBucketTaggingErr = nil
	_, err := DeleteAllBucketTags(&mockS3Client{}, tagOpts)
	assert.Nil(t, err)
}
