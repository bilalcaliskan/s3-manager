//go:build e2e

package file

import (
	"context"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/pkg/aws/types"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/constants"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

func TestExecuteFileCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	ctx := context.Background()
	FileCmd.SetContext(ctx)

	cases := []struct {
		caseName        string
		args            []string
		shouldPass      bool
		listObjectsFunc func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error)
	}{
		{
			"Failure caused by too many arguments",
			[]string{"text1", "text2.txt"},
			false,
			nil,
		},
		{
			"Success matching files",
			[]string{"file3.txt"},
			true,
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{
					Contents: []types.Object{
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
							Key:          aws.String("file1.txt"),
							StorageClass: types.ObjectStorageClassStandard,
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
							Key:          aws.String("file2.txt"),
							StorageClass: types.ObjectStorageClassStandard,
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("file3.txt"),
							StorageClass: types.ObjectStorageClassStandard,
						},
					},
				}, nil
			},
		},
		{
			"Failure caused by ListObjects error",
			[]string{"file3.txt"},
			false,
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return nil, constants.ErrInjected
			},
		},
		{
			"Success no matching files",
			[]string{"file3.txt"},
			true,
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{
					Contents: []types.Object{
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
							Key:          aws.String("../../../testdata/file11.txt"),
							StorageClass: types.ObjectStorageClassStandard,
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
							Key:          aws.String("../../../testdata/file5.txt"),
							StorageClass: types.ObjectStorageClassStandard,
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("../../../testdata/file6.txt"),
							StorageClass: types.ObjectStorageClassStandard,
						},
					},
				}, nil
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.ListObjectsAPI = tc.listObjectsFunc

		FileCmd.SetContext(context.WithValue(FileCmd.Context(), options.S3ClientKey{}, mockS3))
		FileCmd.SetContext(context.WithValue(FileCmd.Context(), options.OptsKey{}, rootOpts))
		FileCmd.SetArgs(tc.args)

		err := FileCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		searchOpts.SetZeroValues()
	}
}
