//go:build unit

package cleaner

import (
	"os"
	"testing"
	"time"

	"github.com/bilalcaliskan/s3-manager/internal/constants"
	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	rootoptions "github.com/bilalcaliskan/s3-manager/cmd/root/options"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/internal/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/clean/options"
	"github.com/stretchr/testify/assert"
)

var (
	defaultListObjectsErr    error
	defaultGetObjectErr      error
	defaultDeleteObjectErr   error
	defaultListObjectsOutput = &s3.ListObjectsOutput{
		Name:        aws.String(""),
		Marker:      aws.String(""),
		MaxKeys:     aws.Int64(1000),
		Prefix:      aws.String(""),
		IsTruncated: aws.Bool(false),
		Contents: []*s3.Object{
			{
				ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
				Key:          aws.String("file1.txt"),
				StorageClass: aws.String("STANDARD"),
				Size:         aws.Int64(1000),
				LastModified: aws.Time(time.Now()),
			},
			{
				ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
				Key:          aws.String("file2.txt"),
				StorageClass: aws.String("STANDARD"),
				Size:         aws.Int64(2000),
				LastModified: aws.Time(time.Now()),
			},
			{
				ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
				Key:          aws.String("file3.txt"),
				StorageClass: aws.String("STANDARD"),
				Size:         aws.Int64(3000),
				LastModified: aws.Time(time.Now()),
			},
		},
	}
	defaultDeleteObjectOutput = &s3.DeleteObjectOutput{
		DeleteMarker:   nil,
		RequestCharged: nil,
		VersionId:      nil,
	}
)

type promptMock struct {
	msg string
	err error
}

func (p promptMock) Run() (string, error) {
	// return expected result
	return p.msg, p.err
}

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

func TestStartCleaning(t *testing.T) {
	cases := []struct {
		caseName string
		expected error
		*options.CleanOptions
		prompt.PromptRunner
		*s3.ListObjectsOutput
		listObjectsErr  error
		deleteObjectErr error
		dryRun          bool
		autoApprove     bool
	}{
		{
			"Success while sorting by lastModificationDate",
			nil,
			&options.CleanOptions{
				MinFileSizeInMb: 0,
				MaxFileSizeInMb: 0,
				FileExtensions:  "",
				FileNamePrefix:  "",
				KeepLastNFiles:  2,
				SortBy:          "lastModificationDate",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			promptMock{
				msg: "y",
				err: nil,
			},
			&s3.ListObjectsOutput{
				Name:        aws.String(""),
				Marker:      aws.String(""),
				MaxKeys:     aws.Int64(1000),
				Prefix:      aws.String(""),
				IsTruncated: aws.Bool(false),
				Contents: []*s3.Object{
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
						Key:          aws.String("file1.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(10000000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
						Key:          aws.String("file2.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(20000000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
						Key:          aws.String("file3.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(30000000),
						LastModified: aws.Time(time.Now()),
					},
				},
			},
			nil,
			nil,
			false,
			false,
		},
		{
			"Success while specifying minfilesizeinmb",
			nil,
			&options.CleanOptions{
				MinFileSizeInMb: 100,
				MaxFileSizeInMb: 0,
				FileExtensions:  "txt",
				FileNamePrefix:  "",
				KeepLastNFiles:  2,
				SortBy:          "lastModificationDate",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			promptMock{
				msg: "y",
				err: nil,
			},
			&s3.ListObjectsOutput{
				Name:        aws.String(""),
				Marker:      aws.String(""),
				MaxKeys:     aws.Int64(1000),
				Prefix:      aws.String(""),
				IsTruncated: aws.Bool(false),
				Contents: []*s3.Object{
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
						Key:          aws.String("file1.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(100000000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
						Key:          aws.String("file2.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(200000000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
						Key:          aws.String("file3.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(300000000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
						Key:          aws.String("file4.json"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(300000000),
						LastModified: aws.Time(time.Now()),
					},
				},
			},
			nil,
			nil,
			false,
			false,
		},
		{
			"Success while specifying both minfilesizeinmb and maxfilesizeinmb",
			nil,
			&options.CleanOptions{
				MinFileSizeInMb: 100,
				MaxFileSizeInMb: 500,
				FileExtensions:  "",
				FileNamePrefix:  "",
				KeepLastNFiles:  2,
				SortBy:          "lastModificationDate",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			promptMock{
				msg: "y",
				err: nil,
			},
			&s3.ListObjectsOutput{
				Name:        aws.String(""),
				Marker:      aws.String(""),
				MaxKeys:     aws.Int64(1000),
				Prefix:      aws.String(""),
				IsTruncated: aws.Bool(false),
				Contents: []*s3.Object{
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
						Key:          aws.String("file1.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(100000000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
						Key:          aws.String("file2.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(200000000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
						Key:          aws.String("file3.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(300000000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
						Key:          aws.String("/"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(300000000),
						LastModified: aws.Time(time.Now()),
					},
				},
			},
			nil,
			nil,
			false,
			false,
		},
		{
			"Success while sorting by size",
			nil,
			&options.CleanOptions{
				MinFileSizeInMb: 0,
				MaxFileSizeInMb: 10000,
				FileExtensions:  "",
				FileNamePrefix:  "",
				KeepLastNFiles:  2,
				SortBy:          "size",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			promptMock{
				msg: "y",
				err: nil,
			},
			&s3.ListObjectsOutput{
				Name:        aws.String(""),
				Marker:      aws.String(""),
				MaxKeys:     aws.Int64(1000),
				Prefix:      aws.String(""),
				IsTruncated: aws.Bool(false),
				Contents: []*s3.Object{
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
						Key:          aws.String("file1.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(1000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
						Key:          aws.String("file2.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(2000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
						Key:          aws.String("file3.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(3000),
						LastModified: aws.Time(time.Now()),
					},
				},
			},
			nil,
			nil,
			false,
			false,
		},
		{
			"Success when dry-run enabled",
			nil,
			&options.CleanOptions{
				MinFileSizeInMb: 0,
				MaxFileSizeInMb: 0,
				FileExtensions:  "",
				FileNamePrefix:  "",
				KeepLastNFiles:  2,
				SortBy:          "lastModificationDate",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			promptMock{
				msg: "y",
				err: nil,
			},
			&s3.ListObjectsOutput{
				Name:        aws.String(""),
				Marker:      aws.String(""),
				MaxKeys:     aws.Int64(1000),
				Prefix:      aws.String(""),
				IsTruncated: aws.Bool(false),
				Contents: []*s3.Object{
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
						Key:          aws.String("file1.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(1000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
						Key:          aws.String("file2.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(2000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
						Key:          aws.String("file3.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(3000),
						LastModified: aws.Time(time.Now()),
					},
				},
			},
			nil,
			nil,
			true,
			false,
		},
		{
			"Failure caused by get all files error",
			constants.ErrInjected,
			&options.CleanOptions{
				MinFileSizeInMb: 0,
				MaxFileSizeInMb: 0,
				FileExtensions:  "",
				FileNamePrefix:  "",
				KeepLastNFiles:  2,
				SortBy:          "lastModificationDate",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			promptMock{
				msg: "y",
				err: nil,
			},
			nil,
			constants.ErrInjected,
			nil,
			false,
			false,
		},
		{
			"Warning caused by no file to remove caused by --keep-last-n-file flag 1",
			nil,
			&options.CleanOptions{
				MinFileSizeInMb: 0,
				MaxFileSizeInMb: 0,
				FileExtensions:  "",
				FileNamePrefix:  "",
				KeepLastNFiles:  5,
				SortBy:          "lastModificationDate",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			promptMock{
				msg: "y",
				err: nil,
			},
			&s3.ListObjectsOutput{
				Name:        aws.String(""),
				Marker:      aws.String(""),
				MaxKeys:     aws.Int64(1000),
				Prefix:      aws.String(""),
				IsTruncated: aws.Bool(false),
				Contents: []*s3.Object{
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
						Key:          aws.String("file1.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(1000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
						Key:          aws.String("file2.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(2000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
						Key:          aws.String("file3.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(3000),
						LastModified: aws.Time(time.Now()),
					},
				},
			},
			nil,
			nil,
			false,
			false,
		},
		{
			"Warning caused by no file to remove caused by --keep-last-n-file flag 2",
			nil,
			&options.CleanOptions{
				MinFileSizeInMb: 0,
				MaxFileSizeInMb: 0,
				FileExtensions:  "",
				FileNamePrefix:  "",
				KeepLastNFiles:  3,
				SortBy:          "lastModificationDate",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			promptMock{
				msg: "y",
				err: nil,
			},
			&s3.ListObjectsOutput{
				Name:        aws.String(""),
				Marker:      aws.String(""),
				MaxKeys:     aws.Int64(1000),
				Prefix:      aws.String(""),
				IsTruncated: aws.Bool(false),
				Contents: []*s3.Object{
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
						Key:          aws.String("file1.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(1000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
						Key:          aws.String("file2.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(2000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
						Key:          aws.String("file3.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(3000),
						LastModified: aws.Time(time.Now()),
					},
				},
			},
			nil,
			nil,
			false,
			false,
		},
		{
			"Failure caused by user terminated the process",
			constants.ErrUserTerminated,
			&options.CleanOptions{
				MinFileSizeInMb: 0,
				MaxFileSizeInMb: 0,
				FileExtensions:  "",
				FileNamePrefix:  "",
				KeepLastNFiles:  2,
				SortBy:          "lastModificationDate",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			promptMock{
				msg: "n",
				err: constants.ErrInjected,
			},
			&s3.ListObjectsOutput{
				Name:        aws.String(""),
				Marker:      aws.String(""),
				MaxKeys:     aws.Int64(1000),
				Prefix:      aws.String(""),
				IsTruncated: aws.Bool(false),
				Contents: []*s3.Object{
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
						Key:          aws.String("file1.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(1000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
						Key:          aws.String("file2.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(2000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
						Key:          aws.String("file3.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(3000),
						LastModified: aws.Time(time.Now()),
					},
				},
			},
			nil,
			nil,
			false,
			false,
		},
		{
			"Failure caused by prompt error",
			constants.ErrInvalidInput,
			&options.CleanOptions{
				MinFileSizeInMb: 0,
				MaxFileSizeInMb: 0,
				FileExtensions:  "",
				FileNamePrefix:  "",
				KeepLastNFiles:  2,
				SortBy:          "lastModificationDate",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			promptMock{
				msg: "asdfadsf",
				err: constants.ErrInjected,
			},
			&s3.ListObjectsOutput{
				Name:        aws.String(""),
				Marker:      aws.String(""),
				MaxKeys:     aws.Int64(1000),
				Prefix:      aws.String(""),
				IsTruncated: aws.Bool(false),
				Contents: []*s3.Object{
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
						Key:          aws.String("file1.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(1000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
						Key:          aws.String("file2.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(2000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
						Key:          aws.String("file3.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(3000),
						LastModified: aws.Time(time.Now()),
					},
				},
			},
			nil,
			nil,
			false,
			false,
		},
		{
			"Failure caused by delete files error",
			constants.ErrInjected,
			&options.CleanOptions{
				MinFileSizeInMb: 0,
				MaxFileSizeInMb: 0,
				FileExtensions:  "",
				FileNamePrefix:  "",
				KeepLastNFiles:  2,
				SortBy:          "lastModificationDate",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			promptMock{
				msg: "y",
				err: nil,
			},
			&s3.ListObjectsOutput{
				Name:        aws.String(""),
				Marker:      aws.String(""),
				MaxKeys:     aws.Int64(1000),
				Prefix:      aws.String(""),
				IsTruncated: aws.Bool(false),
				Contents: []*s3.Object{
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
						Key:          aws.String("file1.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(1000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
						Key:          aws.String("file2.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(2000),
						LastModified: aws.Time(time.Now()),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
						Key:          aws.String("file3.txt"),
						StorageClass: aws.String("STANDARD"),
						Size:         aws.Int64(3000),
						LastModified: aws.Time(time.Now()),
					},
				},
			},
			nil,
			constants.ErrInjected,
			false,
			true,
		},
	}

	for _, tc := range cases {
		tc.DryRun = tc.dryRun
		tc.AutoApprove = tc.autoApprove

		defaultListObjectsOutput = tc.ListObjectsOutput
		defaultListObjectsErr = tc.listObjectsErr

		defaultDeleteObjectErr = tc.deleteObjectErr

		mockSvc := &mockS3Client{}
		assert.NotNil(t, mockSvc)

		err := StartCleaning(mockSvc, tc.PromptRunner, tc.CleanOptions, logging.GetLogger(tc.CleanOptions.RootOptions))
		assert.Equal(t, tc.expected, err)
	}
}
