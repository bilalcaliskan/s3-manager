//go:build e2e

package file

import (
	"context"
	"io"
	"strings"
	"testing"

	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"

	"github.com/stretchr/testify/mock"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

func TestExecuteFileCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	ctx := context.Background()
	FileCmd.SetContext(ctx)

	cases := []struct {
		caseName          string
		args              []string
		shouldPass        bool
		listObjectsErr    error
		listObjectsOutput *s3.ListObjectsOutput
		getObjectErr      error
		getObjectOutput   *s3.GetObjectOutput
	}{
		{
			"Failure caused by too many arguments",
			[]string{"text1", "text2.txt"},
			false,
			nil,
			&s3.ListObjectsOutput{},
			nil,
			&s3.GetObjectOutput{},
		},
		{
			"Success matching files",
			[]string{"file1.txt"},
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
		{
			"Failure caused by ListObjects error",
			[]string{"file1.txt"},
			false,
			constants.ErrInjected,
			&s3.ListObjectsOutput{},
			nil,
			&s3.GetObjectOutput{},
		},
		{
			"Success no matching files",
			[]string{"file1.txt"},
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
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(internalaws.MockS3Client)
		mockS3.On("ListObjects", mock.AnythingOfType("*s3.ListObjectsInput")).Return(tc.listObjectsOutput, tc.listObjectsErr)
		mockS3.On("GetObject", mock.AnythingOfType("*s3.GetObjectInput")).Return(tc.getObjectOutput, tc.getObjectErr)

		FileCmd.SetContext(context.WithValue(FileCmd.Context(), options.S3SvcKey{}, mockS3))
		FileCmd.SetContext(context.WithValue(FileCmd.Context(), options.OptsKey{}, rootOpts))
		FileCmd.SetArgs(tc.args)

		err := FileCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
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
