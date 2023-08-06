//go:build e2e

package clean

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/aws/types"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

func TestExecuteCleanCmd(t *testing.T) {
	ctx := context.Background()
	CleanCmd.SetContext(ctx)

	rootOpts := options.GetMockedRootOptions()
	cases := []struct {
		caseName         string
		args             []string
		shouldPass       bool
		listObjectsFunc  func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error)
		deleteObjectFunc func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	}{
		{
			// TODO: refactor that test
			"Success",
			[]string{},
			true,
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{}, nil
			},
			func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
				return &s3.DeleteObjectOutput{}, nil
			},
		},
		{
			"Failure caused by invalid 'sortBy' flag",
			[]string{"--sort-by=asldkfjalsdkf"},
			false,
			nil,
			nil,
		},
		{
			"Failure caused by invalid 'order' flag",
			[]string{"--order=asldkfjalsdkf"},
			false,
			nil,
			nil,
		},
		{
			"Failure caused by wrong size flags",
			[]string{"--max-size-mb=10", "--min-size-mb=20"},
			false,
			nil,
			nil,
		},
		{
			"Failure caused by ListObjects error",
			[]string{},
			false,
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return nil, constants.ErrInjected
			},
			nil,
		},
		{
			"Failure caused by wrong number of arguments",
			[]string{"foo", "bar"},
			false,
			nil,
			nil,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case '%s'", tc.caseName)

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.ListObjectsAPI = tc.listObjectsFunc
		mockS3.DeleteObjectAPI = tc.deleteObjectFunc

		CleanCmd.SetArgs(tc.args)
		CleanCmd.SetContext(context.WithValue(CleanCmd.Context(), options.S3ClientKey{}, mockS3))
		CleanCmd.SetContext(context.WithValue(CleanCmd.Context(), options.OptsKey{}, rootOpts))

		err := CleanCmd.Execute()
		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		cleanOpts.SetZeroValues()
	}
}
