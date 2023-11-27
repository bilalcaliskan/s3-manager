//go:build unit

package cleaner

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/aws/types"

	"github.com/bilalcaliskan/s3-manager/internal/constants"
	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	rootoptions "github.com/bilalcaliskan/s3-manager/cmd/root/options"

	"github.com/bilalcaliskan/s3-manager/internal/logging"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/bilalcaliskan/s3-manager/cmd/clean/options"
	"github.com/stretchr/testify/assert"
)

// TestStartCleaning is a unit test function that tests the StartCleaning function.
//
// It tests various scenarios and expected outcomes of the StartCleaning function.
func TestStartCleaning(t *testing.T) {
	cases := []struct {
		caseName string
		expected error
		*options.CleanOptions
		prompt.PromptRunner
		listObjectsFunc  func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error)
		deleteObjectFunc func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
		dryRun           bool
		autoApprove      bool
	}{
		{
			"Success while sorting by lastModificationDate descending",
			nil,
			&options.CleanOptions{
				MinFileSizeInMb: 0,
				MaxFileSizeInMb: 0,
				FileExtensions:  "",
				Regex:           "",
				KeepLastNFiles:  2,
				SortBy:          "lastModificationDate",
				Order:           "descending",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{
					Name:        aws.String(""),
					Marker:      aws.String(""),
					MaxKeys:     aws.Int32(1000),
					Prefix:      aws.String(""),
					IsTruncated: aws.Bool(false),
					Contents: []types.Object{
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
							Key:          aws.String("file4.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(10000000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
							Key:          aws.String("file5.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(20000000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("file6.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(30000000),
							LastModified: aws.Time(time.Now()),
						},
					},
				}, nil
			},
			func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
				return &s3.DeleteObjectOutput{}, nil
			},
			false,
			false,
		},
		{
			"Success while sorting by lastModificationDate ascending",
			nil,
			&options.CleanOptions{
				MinFileSizeInMb: 0,
				MaxFileSizeInMb: 0,
				FileExtensions:  "",
				Regex:           "",
				KeepLastNFiles:  2,
				SortBy:          "lastModificationDate",
				Order:           "ascending",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{
					Name:        aws.String(""),
					Marker:      aws.String(""),
					MaxKeys:     aws.Int32(1000),
					Prefix:      aws.String(""),
					IsTruncated: aws.Bool(false),
					Contents: []types.Object{
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
							Key:          aws.String("file4.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(10000000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
							Key:          aws.String("file5.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(20000000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("file6.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(30000000),
							LastModified: aws.Time(time.Now()),
						},
					},
				}, nil
			},
			func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
				return &s3.DeleteObjectOutput{}, nil
			},
			false,
			false,
		},
		{
			"Success while sorting by lastModificationDate default order",
			nil,
			&options.CleanOptions{
				MinFileSizeInMb: 0,
				MaxFileSizeInMb: 0,
				FileExtensions:  "",
				Regex:           "",
				KeepLastNFiles:  2,
				SortBy:          "lastModificationDate",
				Order:           "",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{
					Name:        aws.String(""),
					Marker:      aws.String(""),
					MaxKeys:     aws.Int32(1000),
					Prefix:      aws.String(""),
					IsTruncated: aws.Bool(false),
					Contents: []types.Object{
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
							Key:          aws.String("file4.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(10000000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
							Key:          aws.String("file5.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(20000000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("file6.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(30000000),
							LastModified: aws.Time(time.Now()),
						},
					},
				}, nil
			},
			func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
				return &s3.DeleteObjectOutput{}, nil
			},
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
				Regex:           "",
				KeepLastNFiles:  2,
				SortBy:          "lastModificationDate",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{
					Name:        aws.String(""),
					Marker:      aws.String(""),
					MaxKeys:     aws.Int32(1000),
					Prefix:      aws.String(""),
					IsTruncated: aws.Bool(false),
					Contents: []types.Object{
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
							Key:          aws.String("file4.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(100000000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
							Key:          aws.String("file5.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(200000000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("file6.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(300000000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("file4.json"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(300000000),
							LastModified: aws.Time(time.Now()),
						},
					},
				}, nil
			},
			func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
				return &s3.DeleteObjectOutput{}, nil
			},
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
				Regex:           "",
				KeepLastNFiles:  2,
				SortBy:          "lastModificationDate",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{
					Name:        aws.String(""),
					Marker:      aws.String(""),
					MaxKeys:     aws.Int32(1000),
					Prefix:      aws.String(""),
					IsTruncated: aws.Bool(false),
					Contents: []types.Object{
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
							Key:          aws.String("file4.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(100000000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
							Key:          aws.String("file5.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(200000000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("file6.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(300000000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("/"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(300000000),
							LastModified: aws.Time(time.Now()),
						},
					},
				}, nil
			},
			func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
				return &s3.DeleteObjectOutput{}, nil
			},
			false,
			false,
		},
		{
			"Success while sorting by size descending",
			nil,
			&options.CleanOptions{
				MinFileSizeInMb: 0,
				MaxFileSizeInMb: 10000,
				FileExtensions:  "",
				Regex:           "",
				KeepLastNFiles:  2,
				SortBy:          "size",
				Order:           "descending",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{
					Name:        aws.String(""),
					Marker:      aws.String(""),
					MaxKeys:     aws.Int32(1000),
					Prefix:      aws.String(""),
					IsTruncated: aws.Bool(false),
					Contents: []types.Object{
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
							Key:          aws.String("file4.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(1000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
							Key:          aws.String("file5.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(2000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("file6.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(3000),
							LastModified: aws.Time(time.Now()),
						},
					},
				}, nil
			},
			func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
				return &s3.DeleteObjectOutput{}, nil
			},
			false,
			false,
		},
		{
			"Success while sorting by size ascending",
			nil,
			&options.CleanOptions{
				MinFileSizeInMb: 0,
				MaxFileSizeInMb: 10000,
				FileExtensions:  "",
				Regex:           "",
				KeepLastNFiles:  2,
				SortBy:          "size",
				Order:           "ascending",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{
					Name:        aws.String(""),
					Marker:      aws.String(""),
					MaxKeys:     aws.Int32(1000),
					Prefix:      aws.String(""),
					IsTruncated: aws.Bool(false),
					Contents: []types.Object{
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
							Key:          aws.String("file4.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(1000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
							Key:          aws.String("file5.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(2000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("file6.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(3000),
							LastModified: aws.Time(time.Now()),
						},
					},
				}, nil
			},
			func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
				return &s3.DeleteObjectOutput{}, nil
			},
			false,
			false,
		},
		{
			"Success while sorting by size with default order",
			nil,
			&options.CleanOptions{
				MinFileSizeInMb: 0,
				MaxFileSizeInMb: 10000,
				FileExtensions:  "",
				Regex:           "",
				KeepLastNFiles:  2,
				SortBy:          "size",
				Order:           "",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{
					Name:        aws.String(""),
					Marker:      aws.String(""),
					MaxKeys:     aws.Int32(1000),
					Prefix:      aws.String(""),
					IsTruncated: aws.Bool(false),
					Contents: []types.Object{
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
							Key:          aws.String("file4.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(1000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
							Key:          aws.String("file5.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(2000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("file6.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(3000),
							LastModified: aws.Time(time.Now()),
						},
					},
				}, nil
			},
			func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
				return &s3.DeleteObjectOutput{}, nil
			},
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
				Regex:           "",
				KeepLastNFiles:  2,
				SortBy:          "lastModificationDate",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{
					Name:        aws.String(""),
					Marker:      aws.String(""),
					MaxKeys:     aws.Int32(1000),
					Prefix:      aws.String(""),
					IsTruncated: aws.Bool(false),
					Contents: []types.Object{
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
							Key:          aws.String("file4.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(1000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
							Key:          aws.String("file5.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(2000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("file6.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(3000),
							LastModified: aws.Time(time.Now()),
						},
					},
				}, nil
			},
			func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
				return &s3.DeleteObjectOutput{}, nil
			},
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
				Regex:           "",
				KeepLastNFiles:  2,
				SortBy:          "lastModificationDate",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return nil, constants.ErrInjected
			},
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
				Regex:           "",
				KeepLastNFiles:  5,
				SortBy:          "lastModificationDate",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{
					Name:        aws.String(""),
					Marker:      aws.String(""),
					MaxKeys:     aws.Int32(1000),
					Prefix:      aws.String(""),
					IsTruncated: aws.Bool(false),
					Contents: []types.Object{
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
							Key:          aws.String("file4.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(1000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
							Key:          aws.String("file5.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(2000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("file6.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(3000),
							LastModified: aws.Time(time.Now()),
						},
					},
				}, nil
			},
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
				Regex:           "",
				KeepLastNFiles:  3,
				SortBy:          "lastModificationDate",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{
					Name:        aws.String(""),
					Marker:      aws.String(""),
					MaxKeys:     aws.Int32(1000),
					Prefix:      aws.String(""),
					IsTruncated: aws.Bool(false),
					Contents: []types.Object{
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
							Key:          aws.String("file4.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(1000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
							Key:          aws.String("file5.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(2000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("file6.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(3000),
							LastModified: aws.Time(time.Now()),
						},
					},
				}, nil
			},
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
				Regex:           "",
				KeepLastNFiles:  2,
				SortBy:          "lastModificationDate",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			prompt.PromptMock{
				Msg: "n",
				Err: constants.ErrInjected,
			},
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{
					Name:        aws.String(""),
					Marker:      aws.String(""),
					MaxKeys:     aws.Int32(1000),
					Prefix:      aws.String(""),
					IsTruncated: aws.Bool(false),
					Contents: []types.Object{
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
							Key:          aws.String("file4.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(1000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
							Key:          aws.String("file5.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(2000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("file6.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(3000),
							LastModified: aws.Time(time.Now()),
						},
					},
				}, nil
			},
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
				Regex:           "",
				KeepLastNFiles:  2,
				SortBy:          "lastModificationDate",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			prompt.PromptMock{
				Msg: "asdfadsf",
				Err: constants.ErrInjected,
			},
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{
					Name:        aws.String(""),
					Marker:      aws.String(""),
					MaxKeys:     aws.Int32(1000),
					Prefix:      aws.String(""),
					IsTruncated: aws.Bool(false),
					Contents: []types.Object{
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
							Key:          aws.String("file4.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(1000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
							Key:          aws.String("file5.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(2000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("file6.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(3000),
							LastModified: aws.Time(time.Now()),
						},
					},
				}, nil
			},
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
				Regex:           "",
				KeepLastNFiles:  2,
				SortBy:          "lastModificationDate",
				RootOptions:     rootoptions.GetMockedRootOptions(),
			},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{
					Name:        aws.String(""),
					Marker:      aws.String(""),
					MaxKeys:     aws.Int32(1000),
					Prefix:      aws.String(""),
					IsTruncated: aws.Bool(false),
					Contents: []types.Object{
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
							Key:          aws.String("file4.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(1000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
							Key:          aws.String("file5.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(2000),
							LastModified: aws.Time(time.Now()),
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("file6.txt"),
							StorageClass: types.ObjectStorageClassStandard,
							Size:         aws.Int64(3000),
							LastModified: aws.Time(time.Now()),
						},
					},
				}, nil
			},
			func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
				return nil, constants.ErrInjected
			},
			false,
			true,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)
		tc.DryRun = tc.dryRun
		tc.AutoApprove = tc.autoApprove

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.ListObjectsAPI = tc.listObjectsFunc
		mockS3.DeleteObjectAPI = tc.deleteObjectFunc

		err := StartCleaning(mockS3, tc.PromptRunner, tc.CleanOptions, logging.GetLogger(tc.CleanOptions.RootOptions))
		assert.Equal(t, tc.expected, err)
	}
}
