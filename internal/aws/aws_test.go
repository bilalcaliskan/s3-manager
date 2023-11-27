//go:build unit

package aws

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/aws/types"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	"github.com/pkg/errors"

	options6 "github.com/bilalcaliskan/s3-manager/cmd/bucketpolicy/options"
	options2 "github.com/bilalcaliskan/s3-manager/cmd/search/options"
	options5 "github.com/bilalcaliskan/s3-manager/cmd/transferacceleration/options"

	options4 "github.com/bilalcaliskan/s3-manager/cmd/tags/options"

	options3 "github.com/bilalcaliskan/s3-manager/cmd/versioning/options"

	"github.com/bilalcaliskan/s3-manager/internal/logging"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

// dummyBucketPolicyStr is a dummy bucket policy string used in tests.
var dummyBucketPolicyStr = `
{
  "Statement": [
    {
      "Action": "s3:*",
      "Condition": {
        "Bool": {
          "aws:SecureTransport": "false"
        }
      },
      "Effect": "Deny",
      "Principal": "*",
      "Resource": [
        "arn:aws:s3:::thevpnbeast-releases-1",
        "arn:aws:s3:::thevpnbeast-releases-1/*"
      ],
      "Sid": "RestrictToTLSRequestsOnly"
    }
  ],
  "Version": "2012-10-17"
}
`

// TestDeleteFiles is a test function that tests the behavior of the DeleteFiles function.
//
// It creates test cases with different scenarios and verifies the expected results.
// The test cases include both success and failure cases.
// For the success case, it sets up a mocked S3 client and defines a list of objects to be deleted.
// It expects DeleteFiles to return a nil error.
// For the failure case, it injects an error in the DeleteObject operation of the mocked S3 client.
// It expects DeleteFiles to return a specific error.
//
// The test function iterates through all the test cases and performs the necessary assertions.
func TestDeleteFiles(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	cases := []struct {
		caseName   string
		expected   error
		deleteFunc func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
		dryRun     bool
		objects    []types.Object
	}{
		{
			"Success with non-empty file list",
			nil,
			func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
				return &s3.DeleteObjectOutput{}, nil
			},
			false,
			[]types.Object{
				{
					ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
					Key:          aws.String("../../testdata/file4.txt"),
					StorageClass: types.ObjectStorageClass("STANDART"),
					Size:         aws.Int64(500),
					LastModified: aws.Time(time.Now().Add(-5 * time.Hour)),
				},
				{
					ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
					Key:          aws.String("../../testdata/file5.txt"),
					StorageClass: types.ObjectStorageClass("STANDART"),
					Size:         aws.Int64(1000),
					LastModified: aws.Time(time.Now().Add(-2 * time.Hour)),
				},
				{
					ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
					Key:          aws.String("../../testdata/file6.txt"),
					StorageClass: types.ObjectStorageClass("STANDART"),
					Size:         aws.Int64(1500),
					LastModified: aws.Time(time.Now().Add(-10 * time.Hour)),
				},
			},
		},
		{
			"Failure caused by delete object err",
			constants.ErrInjected,
			func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
				return nil, constants.ErrInjected
			},
			false,
			[]types.Object{
				{
					ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
					Key:          aws.String("../../testdata/file4.txt"),
					StorageClass: types.ObjectStorageClass("STANDART"),
					Size:         aws.Int64(500),
					LastModified: aws.Time(time.Now().Add(-5 * time.Hour)),
				},
				{
					ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
					Key:          aws.String("../../testdata/file5.txt"),
					StorageClass: types.ObjectStorageClass("STANDART"),
					Size:         aws.Int64(1000),
					LastModified: aws.Time(time.Now().Add(-2 * time.Hour)),
				},
				{
					ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
					Key:          aws.String("../../testdata/file6.txt"),
					StorageClass: types.ObjectStorageClass("STANDART"),
					Size:         aws.Int64(1500),
					LastModified: aws.Time(time.Now().Add(-10 * time.Hour)),
				},
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		rootOpts.DryRun = tc.dryRun

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.DeleteObjectAPI = tc.deleteFunc

		assert.Equal(t, tc.expected, DeleteFiles(mockS3, "thisisdemobucket", tc.objects, tc.dryRun, logging.GetLogger(rootOpts)))
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

// TestSearchString is a test function that tests the behavior of the SearchString function.
//
// It creates test cases with different scenarios and verifies the expected results.
// The test cases include both success and failure cases.
// For the success cases, it sets up SearchOptions objects with specific search criteria and a mocked S3 client.
// It expects SearchString to return a non-nil result and a nil error, and it asserts the match count.
// For the failure cases, it either injects an error in the ListObjects operation or the GetObject operation of the mocked S3 client.
// It expects SearchString to return a nil result and a non-nil error.
//
// The test function iterates through all the test cases and performs the necessary assertions.
func TestSearchString(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	cases := []struct {
		caseName        string
		searchOpts      *options2.SearchOptions
		shouldPass      bool
		listObjectsFunc func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error)
		getObjectFunc   func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
		matchCount      int
	}{
		{
			"Success with specific text",
			&options2.SearchOptions{
				Text:        "pvRRTaigmb",
				FileName:    "",
				RootOptions: nil,
			},
			true,
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{
					Contents: []types.Object{
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
							Key:          aws.String("../../testdata/file1.txt"),
							StorageClass: types.ObjectStorageClassStandard,
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
							Key:          aws.String("../../testdata/file2.txt"),
							StorageClass: types.ObjectStorageClassStandard,
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("../../testdata/file3.txt"),
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
					ContentLength: aws.Int64(1000),
					ContentType:   aws.String("text/plain"),
					ETag:          aws.String("d73a503d212d9279e6b2ed8ac6bb81f3"),
				}, nil
			},
			2,
		},
		{
			"Success with file name regex",
			&options2.SearchOptions{
				Text:        "pvRRTaigmb",
				FileName:    "file2.*.",
				RootOptions: nil,
			},
			true,
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{
					Contents: []types.Object{
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
							Key:          aws.String("../../testdata/file1.txt"),
							StorageClass: types.ObjectStorageClassStandard,
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
							Key:          aws.String("../../testdata/file2.txt"),
							StorageClass: types.ObjectStorageClassStandard,
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("../../testdata/file3.txt"),
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
					ContentLength: aws.Int64(1000),
					ContentType:   aws.String("text/plain"),
					ETag:          aws.String("d73a503d212d9279e6b2ed8ac6bb81f3"),
				}, nil
			},
			1,
		},
		{
			"Failure caused by list objects error",
			&options2.SearchOptions{
				Text:        "pvRRTaigmb",
				FileName:    "",
				RootOptions: nil,
			},
			false,
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return nil, constants.ErrInjected
			},
			func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{}, nil
			},
			0,
		},
		{
			"Failure caused by get object error",
			&options2.SearchOptions{
				Text:        "pvRRTaigmb",
				FileName:    "",
				RootOptions: nil,
			},
			false,
			func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{
					Contents: []types.Object{
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
							Key:          aws.String("../../testdata/file1.txttt"),
							StorageClass: types.ObjectStorageClassStandard,
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
							Key:          aws.String("../../testdata/file2.txt"),
							StorageClass: types.ObjectStorageClassStandard,
						},
						{
							ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
							Key:          aws.String("../../testdata/file3.txt"),
							StorageClass: types.ObjectStorageClassStandard,
						},
					},
				}, constants.ErrInjected
			},
			func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				return nil, constants.ErrInjected
			},
			0,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		tc.searchOpts.RootOptions = rootOpts

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.ListObjectsAPI = tc.listObjectsFunc
		mockS3.GetObjectAPI = tc.getObjectFunc

		res, err := SearchString(mockS3, tc.searchOpts)

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		assert.Equal(t, tc.matchCount, len(res))
	}
}

// TestSetBucketVersioning is a test function that tests the behavior of the SetBucketVersioning function.
//
// It creates test cases with different scenarios and verifies the expected results.
// The test cases cover various scenarios related to enabling or disabling bucket versioning.
// Each test case includes VersioningOptions, GetBucketVersioningOutput, and error parameters.
// The function tests different combinations of inputs, including success cases and failure cases.
// It sets up a mocked S3 client and mocks the GetBucketVersioning and PutBucketVersioning operations.
// The function asserts the expected error or nil value based on the scenario being tested.
// It also verifies that the function behaves correctly when dry-run or auto-approve options are enabled.
//
// The test function iterates through all the test cases and performs the necessary assertions.
func TestSetBucketVersioning(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	cases := []struct {
		caseName string
		*options3.VersioningOptions
		getBucketVersioningFunc func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error)
		putBucketVersioningFunc func(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error)
		expected                error
		dryRun                  bool
		autoApprove             bool
		prompt.PromptRunner
	}{
		{
			"Successfully enabled when disabled",
			&options3.VersioningOptions{
				ActualState:  "",
				DesiredState: "enabled",
				RootOptions:  rootOpts,
			},
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return &s3.GetBucketVersioningOutput{
					Status: types.BucketVersioningStatusSuspended,
				}, nil
			},
			func(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
				return &s3.PutBucketVersioningOutput{}, nil
			},
			nil,
			false,
			true,
			nil,
		},
		{
			"Successfully enabled when already enabled",
			&options3.VersioningOptions{
				ActualState:  "",
				DesiredState: "enabled",
				RootOptions:  rootOpts,
			},
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return &s3.GetBucketVersioningOutput{
					Status: types.BucketVersioningStatusEnabled,
				}, nil
			},
			func(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
				return &s3.PutBucketVersioningOutput{}, nil
			},
			nil,
			false,
			false,
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
		},
		{
			"Successfully disabled when enabled",
			&options3.VersioningOptions{
				ActualState:  "",
				DesiredState: "disabled",
				RootOptions:  rootOpts,
			},
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return &s3.GetBucketVersioningOutput{
					Status: types.BucketVersioningStatusEnabled,
				}, nil
			},
			func(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
				return &s3.PutBucketVersioningOutput{}, nil
			},
			nil,
			false,
			false,
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
		},
		{
			"Failure caused by get versioning error",
			&options3.VersioningOptions{
				ActualState:  "",
				DesiredState: "disabled",
				RootOptions:  rootOpts,
			},
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return nil, constants.ErrInjected
			},
			func(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
				return &s3.PutBucketVersioningOutput{}, nil
			},
			constants.ErrInjected,
			false,
			false,
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
		},
		{
			"Failure caused by put error",
			&options3.VersioningOptions{
				ActualState:  "",
				DesiredState: "enabled",
				RootOptions:  rootOpts,
			},
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return &s3.GetBucketVersioningOutput{
					Status: types.BucketVersioningStatusSuspended,
				}, nil
			},
			func(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
				return nil, constants.ErrInjected
			},
			constants.ErrInjected,
			false,
			true,
			nil,
		},
		{
			"Failure caused by unknown status",
			&options3.VersioningOptions{
				ActualState:  "",
				DesiredState: "disabled",
				RootOptions:  rootOpts,
			},
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return &s3.GetBucketVersioningOutput{
					Status: "Enableddddd",
				}, nil
			},
			func(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
				return &s3.PutBucketVersioningOutput{}, nil
			},
			errors.New("unknown status 'Enableddddd' returned from AWS SDK"),
			false,
			false,
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
		},
		{
			"Success when dry-run enabled",
			&options3.VersioningOptions{
				ActualState:  "",
				DesiredState: "enabled",
				RootOptions:  rootOpts,
			},
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return &s3.GetBucketVersioningOutput{
					Status: types.BucketVersioningStatusSuspended,
				}, nil
			},
			func(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
				return &s3.PutBucketVersioningOutput{}, nil
			},
			nil,
			true,
			false,
			nil,
		},
		{
			"Failure caused by user terminated the process",
			&options3.VersioningOptions{
				ActualState:  "",
				DesiredState: "enabled",
				RootOptions:  rootOpts,
			},
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return &s3.GetBucketVersioningOutput{
					Status: types.BucketVersioningStatusSuspended,
				}, nil
			},
			func(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
				return &s3.PutBucketVersioningOutput{}, nil
			},
			constants.ErrUserTerminated,
			false,
			false,
			prompt.PromptMock{
				Msg: "n",
				Err: constants.ErrInjected,
			},
		},
		{
			"Failure caused by prompt error",
			&options3.VersioningOptions{
				ActualState:  "",
				DesiredState: "enabled",
				RootOptions:  rootOpts,
			},
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return &s3.GetBucketVersioningOutput{
					Status: types.BucketVersioningStatusSuspended,
				}, nil
			},
			func(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
				return &s3.PutBucketVersioningOutput{}, nil
			},
			constants.ErrInvalidInput,
			false,
			false,
			prompt.PromptMock{
				Msg: "nasdf",
				Err: constants.ErrInjected,
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		tc.DryRun = tc.dryRun
		tc.AutoApprove = tc.autoApprove

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.GetBucketVersioningAPI = tc.getBucketVersioningFunc
		mockS3.PutBucketVersioningAPI = tc.putBucketVersioningFunc

		err := SetBucketVersioning(mockS3, tc.VersioningOptions, tc.PromptRunner, logging.GetLogger(tc.VersioningOptions.RootOptions))
		if tc.expected != nil {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

// TestGetBucketVersioning is a test function that tests the behavior of the GetBucketVersioning function.
//
// It creates test cases with different scenarios and verifies the expected results.
// The test cases include both success and failure cases.
// For the success case, it sets up a mocked S3 client and defines a GetBucketVersioningOutput object with a specific status.
// It expects GetBucketVersioning to return a non-nil output and a nil error.
// For the failure case, it injects an error in the GetBucketVersioning operation of the mocked S3 client.
// It expects GetBucketVersioning to return a nil output and a non-nil error.
//
// The test function iterates through all the test cases and performs the necessary assertions.
func TestGetBucketVersioning(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	cases := []struct {
		caseName                string
		expected                error
		getBucketVersioningFunc func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error)
	}{
		{
			"Success",
			nil,
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return &s3.GetBucketVersioningOutput{
					Status: types.BucketVersioningStatusEnabled,
				}, nil
			},
		},
		{
			"Failure",
			constants.ErrInjected,
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return nil, constants.ErrInjected
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.GetBucketVersioningAPI = tc.getBucketVersioningFunc

		_, err := GetBucketVersioning(mockS3, rootOpts)
		assert.Equal(t, tc.expected, err)
	}
}

// TestGetBucketTags is a test function that tests the behavior of the GetBucketTags function.
//
// It creates test cases with different scenarios and verifies the expected results.
// The test cases include both success and failure cases.
// For the success case, it sets up a mocked S3 client and defines a GetBucketTaggingOutput object with specific tags.
// It expects GetBucketTags to return the expected tags and a nil error.
// For the failure case, it injects an error in the GetBucketTagging operation of the mocked S3 client.
// It expects GetBucketTags to return a nil result and a non-nil error.
//
// The test function iterates through all the test cases and performs the necessary assertions.
func TestGetBucketTags(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	cases := []struct {
		caseName string
		expected error
		*options4.TagOptions
		getBucketTaggingFunc func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error)
	}{
		{
			"Success",
			nil,
			&options4.TagOptions{
				RootOptions: rootOpts,
			},
			func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
				return &s3.GetBucketTaggingOutput{
					TagSet: []types.Tag{
						{
							Key:   aws.String("foo"),
							Value: aws.String("bar"),
						},
						{
							Key:   aws.String("foo2"),
							Value: aws.String("bar2"),
						},
					},
				}, nil
			},
		},
		{
			"Failure",
			constants.ErrInjected,
			&options4.TagOptions{
				RootOptions: rootOpts,
			},
			func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
				return nil, constants.ErrInjected
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.GetBucketTaggingAPI = tc.getBucketTaggingFunc

		_, err := GetBucketTags(mockS3, tc.TagOptions)
		assert.Equal(t, tc.expected, err)
	}
}

// TestSetBucketTags is a test function that tests the behavior of the SetBucketTags function.
//
// It creates test cases with different scenarios and verifies the expected results.
// The test cases include both success and failure cases.
// For the success case, it sets up a mocked S3 client and defines a list of tags to be added to the bucket.
// It expects SetBucketTags to return a nil error.
// For the failure case, it injects an error in the PutBucketTagging operation of the mocked S3 client.
// It expects SetBucketTags to return a non-nil error.
//
// The test function iterates through all the test cases and performs the necessary assertions.
func TestSetBucketTags(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	logger := logging.GetLogger(rootOpts)
	cases := []struct {
		caseName string
		expected error
		*options4.TagOptions
		tags                 map[string]string
		putBucketTaggingFunc func(ctx context.Context, params *s3.PutBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.PutBucketTaggingOutput, error)
		enableDryRun         bool
		enableAutoApprove    bool
		prompt.PromptRunner
	}{
		{
			"Success",
			nil,
			&options4.TagOptions{
				RootOptions:  rootOpts,
				TagsToAdd:    make(map[string]string),
				TagsToRemove: make(map[string]string),
			},
			map[string]string{
				"foo": "bar",
			},
			func(ctx context.Context, params *s3.PutBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.PutBucketTaggingOutput, error) {
				return &s3.PutBucketTaggingOutput{}, nil
			},
			false,
			false,
			&prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
		},
		{
			"Success with dry run enabled",
			nil,
			&options4.TagOptions{
				RootOptions:  rootOpts,
				TagsToAdd:    make(map[string]string),
				TagsToRemove: make(map[string]string),
			},
			map[string]string{
				"foo":  "bar",
				"foo2": "bar2",
			},
			func(ctx context.Context, params *s3.PutBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.PutBucketTaggingOutput, error) {
				return &s3.PutBucketTaggingOutput{}, nil
			},
			true,
			false,
			nil,
		},
		{
			"Failure",
			constants.ErrInjected,
			&options4.TagOptions{
				RootOptions:  rootOpts,
				TagsToAdd:    make(map[string]string),
				TagsToRemove: make(map[string]string),
			},
			map[string]string{
				"foo":  "bar",
				"foo2": "bar2",
			},
			func(ctx context.Context, params *s3.PutBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.PutBucketTaggingOutput, error) {
				return nil, constants.ErrInjected
			},
			false,
			false,
			&prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
		},
		{
			"Failure caused by user terminated the process",
			constants.ErrUserTerminated,
			&options4.TagOptions{
				RootOptions:  rootOpts,
				TagsToAdd:    make(map[string]string),
				TagsToRemove: make(map[string]string),
			},
			map[string]string{
				"foo":  "bar",
				"foo2": "bar2",
			},
			func(ctx context.Context, params *s3.PutBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.PutBucketTaggingOutput, error) {
				return nil, constants.ErrInjected
			},
			false,
			false,
			&prompt.PromptMock{
				Msg: "n",
				Err: constants.ErrUserTerminated,
			},
		},
		{
			"Failure caused by invalid input",
			constants.ErrInvalidInput,
			&options4.TagOptions{
				RootOptions:  rootOpts,
				TagsToAdd:    make(map[string]string),
				TagsToRemove: make(map[string]string),
			},
			map[string]string{
				"foo":  "bar",
				"foo2": "bar2",
			},
			func(ctx context.Context, params *s3.PutBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.PutBucketTaggingOutput, error) {
				return nil, constants.ErrInjected
			},
			false,
			false,
			&prompt.PromptMock{
				Msg: "asdfasdfy",
				Err: constants.ErrInvalidInput,
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		tc.DryRun = tc.enableDryRun
		tc.AutoApprove = tc.enableAutoApprove

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.PutBucketTaggingAPI = tc.putBucketTaggingFunc

		for k, v := range tc.tags {
			tc.TagOptions.TagsToAdd[k] = v
		}

		assert.Equal(t, tc.expected, SetBucketTags(mockS3, tc.TagOptions, tc.PromptRunner, logger))
	}
}

// TestDeleteAllBucketTags is a test function that tests the behavior of the DeleteAllBucketTags function.
//
// It creates test cases with different scenarios and verifies the expected results.
// The test cases include both success and failure cases.
// For the success case, it sets up a mocked S3 client and expects the DeleteBucketTagging operation to succeed.
// It expects DeleteAllBucketTags to return a nil error.
// For the failure case, it injects an error in the DeleteBucketTagging operation of the mocked S3 client.
// It expects DeleteAllBucketTags to return a non-nil error.
//
// The test function iterates through all the test cases and performs the necessary assertions.
func TestDeleteAllBucketTags(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	logger := logging.GetLogger(rootOpts)
	cases := []struct {
		caseName string
		expected error
		*options4.TagOptions
		deleteBucketTaggingFunc func(ctx context.Context, params *s3.DeleteBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketTaggingOutput, error)
		enableDryRun            bool
		enableAutoApprove       bool
		prompt.PromptRunner
	}{
		{
			"Success",
			nil,
			&options4.TagOptions{
				RootOptions: rootOpts,
			},
			func(ctx context.Context, params *s3.DeleteBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketTaggingOutput, error) {
				return &s3.DeleteBucketTaggingOutput{}, nil
			},
			false,
			false,
			&prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
		},
		{
			"Success with dry-run enabled",
			nil,
			&options4.TagOptions{
				RootOptions: rootOpts,
			},
			func(ctx context.Context, params *s3.DeleteBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketTaggingOutput, error) {
				return &s3.DeleteBucketTaggingOutput{}, nil
			},
			true,
			false,
			nil,
		},
		{
			"Failure",
			constants.ErrInjected,
			&options4.TagOptions{
				RootOptions: rootOpts,
			},
			func(ctx context.Context, params *s3.DeleteBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketTaggingOutput, error) {
				return &s3.DeleteBucketTaggingOutput{}, constants.ErrInjected
			},
			false,
			false,
			&prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
		},
		{
			"Failure caused by user terminated process",
			constants.ErrUserTerminated,
			&options4.TagOptions{
				RootOptions: rootOpts,
			},
			func(ctx context.Context, params *s3.DeleteBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketTaggingOutput, error) {
				return &s3.DeleteBucketTaggingOutput{}, nil
			},
			false,
			false,
			&prompt.PromptMock{
				Msg: "n",
				Err: constants.ErrUserTerminated,
			},
		},
		{
			"Failure caused by invalid input",
			constants.ErrInvalidInput,
			&options4.TagOptions{
				RootOptions: rootOpts,
			},
			func(ctx context.Context, params *s3.DeleteBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketTaggingOutput, error) {
				return &s3.DeleteBucketTaggingOutput{}, nil
			},
			false,
			false,
			&prompt.PromptMock{
				Msg: "nasfassadads",
				Err: constants.ErrInvalidInput,
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		tc.DryRun = tc.enableDryRun
		tc.AutoApprove = tc.enableAutoApprove

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.DeleteBucketTaggingAPI = tc.deleteBucketTaggingFunc

		_, err := DeleteAllBucketTags(mockS3, tc.TagOptions, tc.PromptRunner, logger)
		assert.Equal(t, tc.expected, err)
	}
}

// TestGetBucketPolicy is a test function that tests the behavior of the GetBucketPolicy function.
//
// It creates test cases with different scenarios and verifies the expected results.
// The test cases include both success and failure cases.
// For the success case, it sets up a mocked S3 client and expects the GetBucketPolicy operation to succeed.
// It expects GetBucketPolicy to return a nil error.
// For the failure case, it injects an error in the GetBucketPolicy operation of the mocked S3 client.
// It expects GetBucketPolicy to return a non-nil error.
//
// The test function iterates through all the test cases and performs the necessary assertions.
func TestGetBucketPolicy(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	cases := []struct {
		caseName            string
		expected            error
		getBucketPolicyFunc func(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error)
		*options6.BucketPolicyOptions
		getBucketPolicyErr error
	}{
		{
			"Success",
			nil,
			func(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
				return &s3.GetBucketPolicyOutput{}, nil
			},
			&options6.BucketPolicyOptions{
				RootOptions: rootOpts,
			},
			nil,
		},
		{
			"Failure",
			constants.ErrInjected,
			func(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
				return nil, constants.ErrInjected
			},
			&options6.BucketPolicyOptions{
				RootOptions: rootOpts,
			},
			nil,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockSvc := new(internalawstypes.MockS3Client)
		mockSvc.GetBucketPolicyAPI = tc.getBucketPolicyFunc

		_, err := GetBucketPolicy(mockSvc, tc.BucketPolicyOptions)
		assert.Equal(t, tc.expected, err)
	}
}

// TestSetBucketPolicy is a test function that tests the behavior of the SetBucketPolicy function.
//
// It creates test cases with different scenarios and verifies the expected results.
// The test cases include both success and failure cases.
// For the success case, it sets up a mocked S3 client and expects the PutBucketPolicy operation to succeed.
// It expects SetBucketPolicy to return a nil error.
// For the failure case, it injects an error in the PutBucketPolicy operation of the mocked S3 client.
// It expects SetBucketPolicy to return a non-nil error.
//
// The test function iterates through all the test cases and performs the necessary assertions.
func TestSetBucketPolicy(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	logger := logging.GetLogger(rootOpts)
	cases := []struct {
		caseName string
		expected error
		prompt.PromptRunner
		*options6.BucketPolicyOptions
		putBucketPolicyFunc func(ctx context.Context, params *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error)
		autoApprove         bool
		dryRun              bool
	}{
		{
			"Success",
			nil,
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			&options6.BucketPolicyOptions{
				RootOptions:         rootOpts,
				BucketPolicyContent: dummyBucketPolicyStr,
			},
			func(ctx context.Context, params *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error) {
				return &s3.PutBucketPolicyOutput{}, nil
			},
			false,
			false,
		},
		{
			"Success with dry-run enabled",
			nil,
			nil,
			&options6.BucketPolicyOptions{
				RootOptions:         rootOpts,
				BucketPolicyContent: dummyBucketPolicyStr,
			},
			func(ctx context.Context, params *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error) {
				return &s3.PutBucketPolicyOutput{}, nil
			},
			false,
			true,
		},
		{
			"Failure",
			constants.ErrInjected,
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			&options6.BucketPolicyOptions{
				RootOptions:         rootOpts,
				BucketPolicyContent: dummyBucketPolicyStr,
			},
			func(ctx context.Context, params *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error) {
				return &s3.PutBucketPolicyOutput{}, constants.ErrInjected
			},
			false,
			false,
		},
		{
			"Failure caused by prompt error",
			constants.ErrUserTerminated,
			prompt.PromptMock{
				Msg: "n",
				Err: constants.ErrUserTerminated,
			},
			&options6.BucketPolicyOptions{
				RootOptions:         rootOpts,
				BucketPolicyContent: dummyBucketPolicyStr,
			},
			func(ctx context.Context, params *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error) {
				return &s3.PutBucketPolicyOutput{}, nil
			},
			false,
			false,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		tc.DryRun = tc.dryRun
		tc.AutoApprove = tc.autoApprove

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.PutBucketPolicyAPI = tc.putBucketPolicyFunc

		_, err := SetBucketPolicy(mockS3, tc.BucketPolicyOptions, tc.PromptRunner, logger)
		assert.Equal(t, tc.expected, err)
	}
}

// TestGetBucketPolicyString is a test function that tests the behavior of the GetBucketPolicyString function.
//
// It creates test cases with different scenarios and verifies the expected results.
// The test cases include both success and failure cases.
// For the success case, it sets up a mocked S3 client and expects the GetBucketPolicy operation to succeed.
// It expects GetBucketPolicyString to return the expected bucket policy string and a nil error.
// For the failure case, it injects an error in the GetBucketPolicy operation of the mocked S3 client.
// It expects GetBucketPolicyString to return a non-nil error.
//
// The test function iterates through all the test cases and performs the necessary assertions.
func TestGetBucketPolicyString(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	cases := []struct {
		caseName string
		expected error
		*options6.BucketPolicyOptions
		getBucketPolicyFunc func(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error)
	}{
		{
			"Success",
			nil,
			&options6.BucketPolicyOptions{
				RootOptions:         rootOpts,
				BucketPolicyContent: dummyBucketPolicyStr,
			},
			func(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
				return &s3.GetBucketPolicyOutput{
					Policy: aws.String(dummyBucketPolicyStr),
				}, nil
			},
		},
		{
			"Failure",
			constants.ErrInjected,
			&options6.BucketPolicyOptions{
				RootOptions:         rootOpts,
				BucketPolicyContent: dummyBucketPolicyStr,
			},
			func(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
				return nil, constants.ErrInjected
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.GetBucketPolicyAPI = tc.getBucketPolicyFunc

		_, err := GetBucketPolicyString(mockS3, tc.BucketPolicyOptions)
		if tc.expected == nil {
			assert.Nil(t, err)
		} else {
			assert.Contains(t, err.Error(), tc.expected.Error())
		}
	}
}

// TestDeleteBucketPolicy is a test function that tests the behavior of the DeleteBucketPolicy function.
//
// It creates test cases with different scenarios and verifies the expected results.
// The test cases include both success and failure cases.
// For the success case, it sets up a mocked S3 client and expects the DeleteBucketPolicy operation to succeed.
// It expects DeleteBucketPolicy to return a nil error.
// For the failure case, it injects an error in the DeleteBucketPolicy operation of the mocked S3 client.
// It expects DeleteBucketPolicy to return a non-nil error.
//
// The test function iterates through all the test cases and performs the necessary assertions.
func TestDeleteBucketPolicy(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	logger := logging.GetLogger(rootOpts)
	cases := []struct {
		caseName string
		expected error
		*options6.BucketPolicyOptions
		deleteBucketPolicyFunc func(ctx context.Context, params *s3.DeleteBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketPolicyOutput, error)
		prompt.PromptRunner
		enableAutoApprove bool
		enableDryRun      bool
	}{
		{
			"Success",
			nil,
			&options6.BucketPolicyOptions{
				RootOptions:         rootOpts,
				BucketPolicyContent: dummyBucketPolicyStr,
			},
			func(ctx context.Context, params *s3.DeleteBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketPolicyOutput, error) {
				return &s3.DeleteBucketPolicyOutput{}, nil
			},
			&prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Success with dry-run enabled",
			nil,
			&options6.BucketPolicyOptions{
				RootOptions:         rootOpts,
				BucketPolicyContent: dummyBucketPolicyStr,
			},
			func(ctx context.Context, params *s3.DeleteBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketPolicyOutput, error) {
				return &s3.DeleteBucketPolicyOutput{}, nil
			},
			&prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			true,
		},
		{
			"Failure",
			constants.ErrInjected,
			&options6.BucketPolicyOptions{
				RootOptions:         rootOpts,
				BucketPolicyContent: dummyBucketPolicyStr,
			},
			func(ctx context.Context, params *s3.DeleteBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketPolicyOutput, error) {
				return nil, constants.ErrInjected
			},
			&prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by user terminated process",
			constants.ErrUserTerminated,
			&options6.BucketPolicyOptions{
				RootOptions:         rootOpts,
				BucketPolicyContent: dummyBucketPolicyStr,
			},
			func(ctx context.Context, params *s3.DeleteBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketPolicyOutput, error) {
				return nil, constants.ErrInjected
			},
			&prompt.PromptMock{
				Msg: "n",
				Err: constants.ErrUserTerminated,
			},
			false,
			false,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		tc.DryRun = tc.enableDryRun
		tc.AutoApprove = tc.enableAutoApprove

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.DeleteBucketPolicyAPI = tc.deleteBucketPolicyFunc

		_, err := DeleteBucketPolicy(mockS3, tc.BucketPolicyOptions, tc.PromptRunner, logger)
		assert.Equal(t, tc.expected, err)
	}
}

// TestGetTransferAcceleration is a test function that tests the behavior of the GetTransferAcceleration function.
//
// It creates test cases with different scenarios and verifies the expected results.
// The test cases include both success and failure cases.
// For the success case, it sets up a mocked S3 client and expects the GetBucketAccelerateConfiguration operation to succeed.
// It expects GetTransferAcceleration to return a nil error.
// For the failure case, it injects an error in the GetBucketAccelerateConfiguration operation of the mocked S3 client.
// It expects GetTransferAcceleration to return a non-nil error.
//
// The test function iterates through all the test cases and performs the necessary assertions.
func TestGetTransferAcceleration(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	cases := []struct {
		caseName string
		expected error
		*options5.TransferAccelerationOptions
		getBucketAccelerationErr  error
		getBucketAccelerationFunc func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error)
	}{
		{
			"Success",
			nil,
			&options5.TransferAccelerationOptions{
				RootOptions: rootOpts,
			},
			nil,
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return &s3.GetBucketAccelerateConfigurationOutput{}, nil
			},
		},
		{
			"Failure",
			constants.ErrInjected,
			&options5.TransferAccelerationOptions{
				RootOptions: rootOpts,
			},
			constants.ErrInjected,
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return nil, constants.ErrInjected
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.GetBucketAccelerateConfigurationAPI = tc.getBucketAccelerationFunc

		_, err := GetTransferAcceleration(mockS3, tc.TransferAccelerationOptions)
		assert.Equal(t, tc.expected, err)
	}
}

// TestSetTransferAcceleration is a test function that tests the behavior of the SetTransferAcceleration function.
//
// It creates test cases with different scenarios and verifies the expected results.
// The test cases include both success and failure cases.
// For the success cases, it sets up a mocked S3 client and expects the GetBucketAccelerateConfiguration and PutBucketAccelerateConfiguration operations to succeed.
// It verifies that the function returns a nil error when the expected state is achieved.
// For the failure cases, it injects errors in the GetBucketAccelerateConfiguration and PutBucketAccelerateConfiguration operations of the mocked S3 client.
// It expects the function to return a non-nil error.
// The test function also includes cases where dry-run mode is enabled and prompts are involved.
//
// The function iterates through all the test cases and performs the necessary assertions.
func TestSetTransferAcceleration(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	cases := []struct {
		caseName string
		expected error
		*options5.TransferAccelerationOptions
		getBucketAccelerationFunc func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error)
		putBucketAccelerationFunc func(ctx context.Context, params *s3.PutBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.PutBucketAccelerateConfigurationOutput, error)
		dryRun                    bool
		autoApprove               bool
		prompt.PromptRunner
	}{
		{
			"Success",
			nil,
			&options5.TransferAccelerationOptions{
				RootOptions:  rootOpts,
				DesiredState: "enabled",
			},
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return &s3.GetBucketAccelerateConfigurationOutput{
					Status: types.BucketAccelerateStatusSuspended,
				}, nil
			},
			func(ctx context.Context, params *s3.PutBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.PutBucketAccelerateConfigurationOutput, error) {
				return &s3.PutBucketAccelerateConfigurationOutput{}, nil
			},
			false,
			true,
			nil,
		},
		{
			"Success when already enabled",
			nil,
			&options5.TransferAccelerationOptions{
				RootOptions:  rootOpts,
				DesiredState: "enabled",
			},
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return &s3.GetBucketAccelerateConfigurationOutput{
					Status: types.BucketAccelerateStatusEnabled,
				}, nil
			},
			func(ctx context.Context, params *s3.PutBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.PutBucketAccelerateConfigurationOutput, error) {
				return &s3.PutBucketAccelerateConfigurationOutput{}, nil
			},
			false,
			false,
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
		},
		{
			"Success when already disabled",
			nil,
			&options5.TransferAccelerationOptions{
				RootOptions:  rootOpts,
				DesiredState: "disabled",
			},
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return &s3.GetBucketAccelerateConfigurationOutput{
					Status: types.BucketAccelerateStatusSuspended,
				}, nil
			},
			func(ctx context.Context, params *s3.PutBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.PutBucketAccelerateConfigurationOutput, error) {
				return &s3.PutBucketAccelerateConfigurationOutput{}, nil
			},
			false,
			false,
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
		},
		{
			"Success when dry-run enabled",
			nil,
			&options5.TransferAccelerationOptions{
				RootOptions:  rootOpts,
				DesiredState: "disabled",
			},
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return &s3.GetBucketAccelerateConfigurationOutput{
					Status: types.BucketAccelerateStatusSuspended,
				}, nil
			},
			func(ctx context.Context, params *s3.PutBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.PutBucketAccelerateConfigurationOutput, error) {
				return &s3.PutBucketAccelerateConfigurationOutput{}, nil
			},
			true,
			false,
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
		},
		{
			"Failure caused by get transfer acceleration error",
			constants.ErrInjected,
			&options5.TransferAccelerationOptions{
				RootOptions:  rootOpts,
				DesiredState: "disabled",
			},
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return nil, constants.ErrInjected
			},
			func(ctx context.Context, params *s3.PutBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.PutBucketAccelerateConfigurationOutput, error) {
				return &s3.PutBucketAccelerateConfigurationOutput{}, nil
			},
			false,
			false,
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
		},
		{
			"Failure caused by unknown status returned by get transfer acceleration",
			constants.ErrInjected,
			&options5.TransferAccelerationOptions{
				RootOptions:  rootOpts,
				DesiredState: "disabled",
			},
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return &s3.GetBucketAccelerateConfigurationOutput{
					Status: "Suspendedddd",
				}, nil
			},
			func(ctx context.Context, params *s3.PutBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.PutBucketAccelerateConfigurationOutput, error) {
				return &s3.PutBucketAccelerateConfigurationOutput{}, nil
			},
			false,
			false,
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
		},
		{
			"Failure caused by put transfer acceleration error",
			constants.ErrInjected,
			&options5.TransferAccelerationOptions{
				RootOptions:  rootOpts,
				DesiredState: "disabled",
			},
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return &s3.GetBucketAccelerateConfigurationOutput{
					Status: types.BucketAccelerateStatusEnabled,
				}, nil
			},
			func(ctx context.Context, params *s3.PutBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.PutBucketAccelerateConfigurationOutput, error) {
				return &s3.PutBucketAccelerateConfigurationOutput{}, constants.ErrInjected
			},
			false,
			false,
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
		},
		{
			"Failure caused by prompt error",
			constants.ErrInvalidInput,
			&options5.TransferAccelerationOptions{
				RootOptions:  rootOpts,
				DesiredState: "disabled",
			},
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return &s3.GetBucketAccelerateConfigurationOutput{
					Status: types.BucketAccelerateStatusEnabled,
				}, nil
			},
			func(ctx context.Context, params *s3.PutBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.PutBucketAccelerateConfigurationOutput, error) {
				return &s3.PutBucketAccelerateConfigurationOutput{}, nil
			},
			false,
			false,
			prompt.PromptMock{
				Msg: "dkslfa",
				Err: constants.ErrInjected,
			},
		},
		{
			"Failure caused by user terminated the process",
			constants.ErrUserTerminated,
			&options5.TransferAccelerationOptions{
				RootOptions:  rootOpts,
				DesiredState: "disabled",
			},
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return &s3.GetBucketAccelerateConfigurationOutput{
					Status: types.BucketAccelerateStatusEnabled,
				}, nil
			},
			func(ctx context.Context, params *s3.PutBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.PutBucketAccelerateConfigurationOutput, error) {
				return &s3.PutBucketAccelerateConfigurationOutput{}, nil
			},
			false,
			false,
			prompt.PromptMock{
				Msg: "n",
				Err: constants.ErrInjected,
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		tc.DryRun = tc.dryRun
		tc.AutoApprove = tc.autoApprove

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.GetBucketAccelerateConfigurationAPI = tc.getBucketAccelerationFunc
		mockS3.PutBucketAccelerateConfigurationAPI = tc.putBucketAccelerationFunc

		err := SetTransferAcceleration(mockS3, tc.TransferAccelerationOptions, tc.PromptRunner, logging.GetLogger(tc.RootOptions))
		if tc.expected == nil {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}
