//go:build e2e

package file

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/stretchr/testify/assert"
)

var (
	defaultListObjectsErr    error
	defaultListObjectsOutput = &s3.ListObjectsOutput{
		Contents: []*s3.Object{
			{
				ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
				Key:          aws.String("../../../testdata/file1.txt"),
				StorageClass: aws.String("STANDARD"),
			},
			{
				ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
				Key:          aws.String("../../../testdata/file2.txt"),
				StorageClass: aws.String("STANDARD"),
			},
			{
				ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
				Key:          aws.String("../../../testdata/file3.txt"),
				StorageClass: aws.String("STANDARD"),
			},
		},
	}
	defaultGetObjectErr    error
	defaultGetObjectOutput = &s3.GetObjectOutput{}
)

func createSvc(rootOpts *options.RootOptions) (*s3.S3, error) {
	return internalaws.CreateAwsService(rootOpts)
}

// Define a testdata struct to be used in your unit tests
type mockS3Client struct {
	s3iface.S3API
}

// ListObjects mocks the S3API ListObjects method
func (m *mockS3Client) ListObjects(obj *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	return defaultListObjectsOutput, defaultListObjectsErr
}

// GetObject mocks the S3API GetObject method
func (m *mockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return defaultGetObjectOutput, defaultGetObjectErr
}

func TestExecuteFileCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	ctx := context.Background()
	FileCmd.SetContext(ctx)

	cases := []struct {
		caseName          string
		args              []string
		shouldMock        bool
		shouldPass        bool
		listObjectsErr    error
		listObjectsOutput *s3.ListObjectsOutput
		getObjectErr      error
		getObjectOutput   *s3.GetObjectOutput
	}{
		{"Failure caused by too many arguments",
			[]string{"text1", "text2.txt"},
			true,
			false,
			nil,
			&s3.ListObjectsOutput{},
			nil,
			&s3.GetObjectOutput{},
		},
		{"Success matching files",
			[]string{"file1.txt"},
			true,
			true,
			nil,
			&s3.ListObjectsOutput{
				Contents: []*s3.Object{
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
						Key:          aws.String("../../../testdata/file1.txt"),
						StorageClass: aws.String("STANDARD"),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
						Key:          aws.String("../../../testdata/file2.txt"),
						StorageClass: aws.String("STANDARD"),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
						Key:          aws.String("../../../testdata/file3.txt"),
						StorageClass: aws.String("STANDARD"),
					},
				},
			},
			nil,
			&s3.GetObjectOutput{
				Body: getMockBody("BrYKzqcTqn"),
			},
		},
		{"Failure caused by ListObjects error",
			[]string{"file1.txt"},
			true,
			false,
			errors.New("injected error"),
			&s3.ListObjectsOutput{},
			nil,
			&s3.GetObjectOutput{},
		},
		{"Success no matching files",
			[]string{"file1.txt"},
			true,
			true,
			nil,
			&s3.ListObjectsOutput{
				Contents: []*s3.Object{
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
						Key:          aws.String("../../../testdata/file11.txt"),
						StorageClass: aws.String("STANDARD"),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
						Key:          aws.String("../../../testdata/file2.txt"),
						StorageClass: aws.String("STANDARD"),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
						Key:          aws.String("../../../testdata/file3.txt"),
						StorageClass: aws.String("STANDARD"),
					},
				},
			},
			nil,
			&s3.GetObjectOutput{
				Body: getMockBody("BrYKzqcTqn"),
			},
		},
	}

	for _, tc := range cases {
		defaultListObjectsErr = tc.listObjectsErr
		defaultListObjectsOutput = tc.listObjectsOutput

		defaultGetObjectErr = tc.getObjectErr
		defaultGetObjectOutput = tc.getObjectOutput

		var err error
		if tc.shouldMock {
			mockSvc := &mockS3Client{}
			svc = mockSvc
			assert.NotNil(t, mockSvc)
		} else {
			svc, err = createSvc(rootOpts)
			assert.NotNil(t, svc)
			assert.Nil(t, err)
		}

		FileCmd.SetContext(context.WithValue(FileCmd.Context(), options.S3SvcKey{}, svc))
		FileCmd.SetContext(context.WithValue(FileCmd.Context(), options.OptsKey{}, rootOpts))
		FileCmd.SetArgs(tc.args)

		err = FileCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		rootOpts.SetZeroValues()
		searchOpts.SetZeroValues()
	}
}

// getMockBody returns a mock implementation of io.ReadCloser
func getMockBody(str string) io.ReadCloser {
	// Create a ReadCloser implementation using strings.NewReader
	// strings.NewReader returns a new Reader reading from the provided string
	body := strings.NewReader(str)

	// Return the ReadCloser implementation
	return io.NopCloser(body)
}

/*func TestExecuteNotEnoughArgument(t *testing.T) {
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
	FileCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	FileCmd.SetContext(context.WithValue(FileCmd.Context(), options.S3SvcKey{}, svc))
	FileCmd.SetContext(context.WithValue(FileCmd.Context(), options.OptsKey{}, rootOpts))

	FileCmd.SetArgs([]string{})
	err = FileCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	searchOpts.SetZeroValues()
}

func TestExecuteTooManyArguments(t *testing.T) {
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
	FileCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	FileCmd.SetContext(context.WithValue(FileCmd.Context(), options.S3SvcKey{}, svc))
	FileCmd.SetContext(context.WithValue(FileCmd.Context(), options.OptsKey{}, rootOpts))

	FileCmd.SetArgs([]string{"adsfasdf", "asdfasdfadsf"})
	err = FileCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	searchOpts.SetZeroValues()
}

func TestExecuteSuccessWithNoMatchingFiles(t *testing.T) {
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
	FileCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultListObjectsErr = nil
	defaultListObjectsOutput.Contents = []*s3.Object{
		{
			ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
			Key:          aws.String("../../../testdata/file1.txt"),
			StorageClass: aws.String("STANDARD"),
		},
		{
			ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
			Key:          aws.String("../../../testdata/file2.txt"),
			StorageClass: aws.String("STANDARD"),
		},
		{
			ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
			Key:          aws.String("../../../testdata/file3.txt"),
			StorageClass: aws.String("STANDARD"),
		},
	}

	FileCmd.SetContext(context.WithValue(FileCmd.Context(), options.S3SvcKey{}, mockSvc))
	FileCmd.SetContext(context.WithValue(FileCmd.Context(), options.OptsKey{}, rootOpts))

	FileCmd.SetArgs([]string{"aslkdads"})
	err := FileCmd.Execute()
	assert.Nil(t, err)

	rootOpts.SetZeroValues()
	searchOpts.SetZeroValues()
}

func TestExecuteSuccessWithMatchingFiles(t *testing.T) {
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
	FileCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultListObjectsErr = nil
	defaultListObjectsOutput.Contents = []*s3.Object{
		{
			ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
			Key:          aws.String("../../../testdata/file1.txt"),
			StorageClass: aws.String("STANDARD"),
		},
		{
			ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
			Key:          aws.String("../../../testdata/file2.txt"),
			StorageClass: aws.String("STANDARD"),
		},
		{
			ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
			Key:          aws.String("../../../testdata/file3.txt"),
			StorageClass: aws.String("STANDARD"),
		},
	}

	FileCmd.SetContext(context.WithValue(FileCmd.Context(), options.S3SvcKey{}, mockSvc))
	FileCmd.SetContext(context.WithValue(FileCmd.Context(), options.OptsKey{}, rootOpts))

	FileCmd.SetArgs([]string{".*.txt"})
	err := FileCmd.Execute()
	assert.Nil(t, err)

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
	FileCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultListObjectsErr = errors.New("dummy error")
	defaultListObjectsOutput.Contents = []*s3.Object{
		{
			ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
			Key:          aws.String("../../../testdata/file1.txt"),
			StorageClass: aws.String("STANDARD"),
		},
		{
			ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
			Key:          aws.String("../../../testdata/file2.txt"),
			StorageClass: aws.String("STANDARD"),
		},
		{
			ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
			Key:          aws.String("../../../testdata/file3.txt"),
			StorageClass: aws.String("STANDARD"),
		},
	}

	FileCmd.SetContext(context.WithValue(FileCmd.Context(), options.S3SvcKey{}, mockSvc))
	FileCmd.SetContext(context.WithValue(FileCmd.Context(), options.OptsKey{}, rootOpts))

	FileCmd.SetArgs([]string{".*.txt"})
	err := FileCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	searchOpts.SetZeroValues()
}*/
