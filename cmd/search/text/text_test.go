//go:build e2e

package text

import (
	"context"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/aws/types"
	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

var mu sync.Mutex

func TestExecuteTextCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	ctx := context.Background()
	TextCmd.SetContext(ctx)

	cases := []struct {
		caseName        string
		args            []string
		shouldPass      bool
		listObjectsFunc func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error)
		getObjectFunc   func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	}{
		{
			"Success no matching files",
			[]string{"text1", "--file-name=text2.txt"},
			true,
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{}, nil
			},
			func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{}, nil
			},
		},
		{
			"Success matching files",
			[]string{"jPIrSIgOcZ", "--file-name=.*.txt"},
			true,
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{
					Contents: []types.Object{
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
							Key:          aws.String("../../../testdata/file1.txt"),
							StorageClass: types.ObjectStorageClassStandard,
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
							Key:          aws.String("../../../testdata/file2.txt"),
							StorageClass: types.ObjectStorageClassStandard,
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("../../../testdata/file3.txt"),
							StorageClass: types.ObjectStorageClassStandard,
						},
					},
				}, nil
			},
			func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				body := getMockBody(*params.Key)

				return &s3.GetObjectOutput{
					AcceptRanges:  aws.String("bytes"),
					Body:          body,
					ContentLength: 1000,
					ContentType:   aws.String("text/plain"),
					ETag:          aws.String("d73a503d212d9279e6b2ed8ac6bb81f3"),
				}, nil
			},
		},
		{
			"Failure caused by ListObjects error",
			[]string{"text1", "--file-name=text2.txt"},
			false,
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return nil, constants.ErrInjected
			},
			func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{}, nil
			},
		},
		{
			"Failure caused by no arguments",
			[]string{"--file-name=text2.txt"},
			false,
			nil,
			nil,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.ListObjectsAPI = tc.listObjectsFunc
		mockS3.GetObjectAPI = tc.getObjectFunc

		TextCmd.SetContext(context.WithValue(TextCmd.Context(), options.S3ClientKey{}, mockS3))
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
func getMockBody(path string) io.ReadCloser {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	// Create a ReadCloser implementation using strings.NewReader
	// strings.NewReader returns a new Reader reading from the provided string
	body := strings.NewReader(string(content))

	// Return the ReadCloser implementation
	return io.NopCloser(body)
}
