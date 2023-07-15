//go:build e2e

package text

import (
	"context"
	"io"
	"strings"
	"sync"
	"testing"

	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"

	"github.com/stretchr/testify/mock"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

var mu sync.Mutex

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
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(internalaws.MockS3Client)
		mockS3.On("ListObjects", mock.AnythingOfType("*s3.ListObjectsInput")).Return(tc.listObjectsOutput, tc.listObjectsErr)
		mockS3.On("GetObject", mock.AnythingOfType("*s3.GetObjectInput")).Return(tc.getObjectOutput, tc.getObjectErr)

		TextCmd.SetContext(context.WithValue(TextCmd.Context(), options.S3SvcKey{}, mockS3))
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
