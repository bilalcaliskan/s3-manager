//go:build e2e

package text

import (
	"context"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

var (
	defaultListObjectsErr    error
	defaultListObjectsOutput = &s3.ListObjectsOutput{}
	defaultGetObjectErr      error
	defaultGetObjectOutput   = &s3.GetObjectOutput{}
	mu                       sync.Mutex
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
	return defaultGetObjectOutput, defaultGetObjectErr
}

func TestExecuteTextCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	ctx := context.Background()
	TextCmd.SetContext(ctx)

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
			"Success no matching files",
			[]string{"text1", "--file-name=text2.txt"},
			true,
			nil,
			&s3.ListObjectsOutput{},
			nil,
			&s3.GetObjectOutput{},
		},
		{
			"Success matching files",
			[]string{"BrYKzqcTqn", "--file-name=../../../testdata/.*"},
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
			[]string{"text1", "--file-name=text2.txt"},
			false,
			constants.ErrInjected,
			&s3.ListObjectsOutput{},
			nil,
			&s3.GetObjectOutput{},
		},
		{
			"Failure caused by no arguments",
			[]string{"--file-name=text2.txt"},
			false,
			nil,
			&s3.ListObjectsOutput{},
			nil,
			&s3.GetObjectOutput{},
		},
	}

	for _, tc := range cases {
		defaultListObjectsErr = tc.listObjectsErr
		defaultListObjectsOutput = tc.listObjectsOutput

		defaultGetObjectErr = tc.getObjectErr
		defaultGetObjectOutput = tc.getObjectOutput

		TextCmd.SetContext(context.WithValue(TextCmd.Context(), options.S3SvcKey{}, &mockS3Client{}))
		TextCmd.SetContext(context.WithValue(TextCmd.Context(), options.OptsKey{}, rootOpts))
		TextCmd.SetArgs(tc.args)

		mu.Lock()
		err := TextCmd.Execute()
		mu.Unlock()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

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
