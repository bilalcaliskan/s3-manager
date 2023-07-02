//go:build e2e

package text

import (
	"context"
	"errors"
	"io"
	"strings"
	"sync"
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
	lock                   sync.Mutex
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

func TestExecuteTextCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	ctx := context.Background()
	TextCmd.SetContext(ctx)

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
		{"Success no matching files",
			[]string{"text1", "--file-name=text2.txt"},
			true,
			true,
			nil,
			&s3.ListObjectsOutput{},
			nil,
			&s3.GetObjectOutput{},
		},
		{"Success matching files",
			[]string{"BrYKzqcTqn", "--file-name=../../../testdata/.*"},
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
			[]string{"text1", "--file-name=text2.txt"},
			true,
			false,
			errors.New("injected error"),
			&s3.ListObjectsOutput{},
			nil,
			&s3.GetObjectOutput{},
		},
		{"Failure caused by no arguments",
			[]string{"--file-name=text2.txt"},
			true,
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

		TextCmd.SetContext(context.WithValue(TextCmd.Context(), options.S3SvcKey{}, svc))
		TextCmd.SetContext(context.WithValue(TextCmd.Context(), options.OptsKey{}, rootOpts))
		TextCmd.SetArgs(tc.args)

		lock.Lock()
		err = TextCmd.Execute()
		lock.Unlock()

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
