//go:build unit

package aws

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	"github.com/pkg/errors"

	options6 "github.com/bilalcaliskan/s3-manager/cmd/bucketpolicy/options"
	options2 "github.com/bilalcaliskan/s3-manager/cmd/search/options"
	options5 "github.com/bilalcaliskan/s3-manager/cmd/transferacceleration/options"

	options4 "github.com/bilalcaliskan/s3-manager/cmd/tags/options"

	options3 "github.com/bilalcaliskan/s3-manager/cmd/versioning/options"

	"github.com/bilalcaliskan/s3-manager/internal/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

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

// TestGetAllFiles is a test function that tests the behavior of the GetAllFiles function.
//
// It creates test cases with different scenarios and verifies the expected results.
// The test cases include both success and failure cases.
// For the success case, it sets up a mocked S3 client and defines a list of objects to be returned.
// It expects GetAllFiles to return a nil error.
// For the failure case, it injects an error in the ListObjects operation of the mocked S3 client.
// It expects GetAllFiles to return a specific error.
//
// The test function iterates through all the test cases and performs the necessary assertions.
func TestGetAllFiles(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	cases := []struct {
		caseName          string
		expected          error
		listObjectsErr    error
		listObjectsOutput *s3.ListObjectsOutput
	}{
		{
			"Success with non-empty file list",
			nil,
			nil,
			&s3.ListObjectsOutput{
				Contents: []*s3.Object{
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
						Key:          aws.String("../../testdata/file4.txt"),
						StorageClass: aws.String("STANDARD"),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
						Key:          aws.String("../../testdata/file5.txt"),
						StorageClass: aws.String("STANDARD"),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
						Key:          aws.String("../../testdata/file6.txt"),
						StorageClass: aws.String("STANDARD"),
					},
				},
			},
		},
		{
			"Failure caused by List objects error",
			constants.ErrInjected,
			constants.ErrInjected,
			nil,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(MockS3Client)
		mockS3.On("ListObjects", mock.AnythingOfType("*s3.ListObjectsInput")).Return(tc.listObjectsOutput, tc.listObjectsErr)

		_, err := GetAllFiles(mockS3, rootOpts, "")
		assert.Equal(t, tc.expected, err)
	}
}

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
		caseName        string
		expected        error
		deleteObjectErr error
		dryRun          bool
		objects         []*s3.Object
	}{
		{
			"Success with non-empty file list",
			nil,
			nil,
			false,
			[]*s3.Object{
				{
					ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
					Key:          aws.String("../../testdata/file4.txt"),
					StorageClass: aws.String("STANDARD"),
					Size:         aws.Int64(500),
					LastModified: aws.Time(time.Now().Add(-5 * time.Hour)),
				},
				{
					ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
					Key:          aws.String("../../testdata/file5.txt"),
					StorageClass: aws.String("STANDARD"),
					Size:         aws.Int64(1000),
					LastModified: aws.Time(time.Now().Add(-2 * time.Hour)),
				},
				{
					ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
					Key:          aws.String("../../testdata/file6.txt"),
					StorageClass: aws.String("STANDARD"),
					Size:         aws.Int64(1500),
					LastModified: aws.Time(time.Now().Add(-10 * time.Hour)),
				},
			},
		},
		{
			"Failure caused by delete object err",
			constants.ErrInjected,
			constants.ErrInjected,
			false,
			[]*s3.Object{
				{
					ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
					Key:          aws.String("../../testdata/file4.txt"),
					StorageClass: aws.String("STANDARD"),
					Size:         aws.Int64(500),
					LastModified: aws.Time(time.Now().Add(-5 * time.Hour)),
				},
				{
					ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
					Key:          aws.String("../../testdata/file5.txt"),
					StorageClass: aws.String("STANDARD"),
					Size:         aws.Int64(1000),
					LastModified: aws.Time(time.Now().Add(-2 * time.Hour)),
				},
				{
					ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
					Key:          aws.String("../../testdata/file6.txt"),
					StorageClass: aws.String("STANDARD"),
					Size:         aws.Int64(1500),
					LastModified: aws.Time(time.Now().Add(-10 * time.Hour)),
				},
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(MockS3Client)
		mockS3.On("DeleteObject", mock.AnythingOfType("*s3.DeleteObjectInput")).Return(&s3.DeleteObjectOutput{}, tc.deleteObjectErr)

		assert.Equal(t, tc.expected, DeleteFiles(mockS3, "thisisdemobucket", tc.objects, tc.dryRun, logging.GetLogger(rootOpts)))
	}
}

// TestCreateAwsService is a test function that tests the behavior of the CreateAwsService function.
//
// It creates test cases with different scenarios and verifies the expected results.
// The test cases include both success and failure cases.
// For the success case, it sets up a RootOptions object with all the required fields.
// It expects CreateAwsService to return a non-nil AWS service client and a nil error.
// For the failure case, it sets up a RootOptions object with a missing required field.
// It expects CreateAwsService to return a nil AWS service client and a non-nil error.
//
// The test function iterates through all the test cases and performs the necessary assertions.
func TestCreateAwsService(t *testing.T) {
	cases := []struct {
		caseName   string
		opts       *options.RootOptions
		shouldPass bool
	}{
		{
			"Success",
			&options.RootOptions{
				AccessKey:  "thisisaccesskey",
				SecretKey:  "thisissecretkey",
				BucketName: "thisisbucketname",
				Region:     "thisisregion",
			},
			true,
		},
		{
			"Failure caused by missing required field",
			&options.RootOptions{
				AccessKey:  "thisisaccesskey",
				SecretKey:  "thisissecretkey",
				BucketName: "thisisbucketname",
				Region:     "",
			},
			false,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		_, err := CreateAwsService(tc.opts)

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
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
		caseName          string
		searchOpts        *options2.SearchOptions
		shouldPass        bool
		listObjectsErr    error
		listObjectsOutput *s3.ListObjectsOutput
		getObjectErr      error
		matchCount        int
	}{
		{
			"Success with specific text",
			&options2.SearchOptions{
				Text:        "pvRRTaigmb",
				FileName:    "",
				RootOptions: nil,
			},
			true,
			nil,
			&s3.ListObjectsOutput{
				Contents: []*s3.Object{
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
						Key:          aws.String("../../testdata/file1.txt"),
						StorageClass: aws.String("STANDARD"),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
						Key:          aws.String("../../testdata/file2.txt"),
						StorageClass: aws.String("STANDARD"),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
						Key:          aws.String("../../testdata/file3.txt"),
						StorageClass: aws.String("STANDARD"),
					},
				},
			},
			nil,
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
			nil,
			&s3.ListObjectsOutput{
				Contents: []*s3.Object{
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
						Key:          aws.String("../../testdata/file1.txt"),
						StorageClass: aws.String("STANDARD"),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
						Key:          aws.String("../../testdata/file2.txt"),
						StorageClass: aws.String("STANDARD"),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
						Key:          aws.String("../../testdata/file3.txt"),
						StorageClass: aws.String("STANDARD"),
					},
				},
			},
			nil,
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
			constants.ErrInjected,
			nil,
			nil,
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
			nil,
			&s3.ListObjectsOutput{
				Contents: []*s3.Object{
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
						Key:          aws.String("../../testdata/file1.txttt"),
						StorageClass: aws.String("STANDARD"),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
						Key:          aws.String("../../testdata/file2.txt"),
						StorageClass: aws.String("STANDARD"),
					},
					{
						ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
						Key:          aws.String("../../testdata/file3.txt"),
						StorageClass: aws.String("STANDARD"),
					},
				},
			},
			constants.ErrInjected,
			0,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		tc.searchOpts.RootOptions = rootOpts

		mockS3 := new(MockS3Client)
		mockS3.On("ListObjects", mock.AnythingOfType("*s3.ListObjectsInput")).Return(tc.listObjectsOutput, tc.listObjectsErr)
		mockS3.On("GetObject", mock.AnythingOfType("*s3.GetObjectInput")).Return(&s3.GetObjectOutput{}, tc.getObjectErr)

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
		*s3.GetBucketVersioningOutput
		getBucketVersioningErr error
		putBucketVersioningErr error
		expected               error
		dryRun                 bool
		autoApprove            bool
		prompt.PromptRunner
	}{
		{
			"Successfully enabled when disabled",
			&options3.VersioningOptions{
				ActualState:  "",
				DesiredState: "enabled",
				RootOptions:  rootOpts,
			},
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Suspended"),
			},
			nil,
			nil,
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
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			},
			nil,
			nil,
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
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			},
			nil,
			nil,
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
			nil,
			constants.ErrInjected,
			nil,
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
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Suspended"),
			},
			nil,
			constants.ErrInjected,
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
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enableddddd"),
			},
			nil,
			nil,
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
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Suspended"),
			},
			nil,
			nil,
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
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Suspended"),
			},
			nil,
			nil,
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
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Suspended"),
			},
			nil,
			nil,
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

		mockS3 := new(MockS3Client)
		mockS3.On("GetBucketVersioning", mock.AnythingOfType("*s3.GetBucketVersioningInput")).Return(tc.GetBucketVersioningOutput, tc.getBucketVersioningErr)
		mockS3.On("PutBucketVersioning", mock.AnythingOfType("*s3.PutBucketVersioningInput")).Return(&s3.PutBucketVersioningOutput{}, tc.putBucketVersioningErr)

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
		caseName string
		expected error
		*s3.GetBucketVersioningOutput
		getBucketVersioningErr error
	}{
		{
			"Success",
			nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enableddddd"),
			},
			nil,
		},
		{
			"Failure",
			constants.ErrInjected,
			nil,
			constants.ErrInjected,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(MockS3Client)
		mockS3.On("GetBucketVersioning", mock.AnythingOfType("*s3.GetBucketVersioningInput")).Return(tc.GetBucketVersioningOutput, tc.getBucketVersioningErr)

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
		*s3.GetBucketTaggingOutput
		getBucketTaggingErr error
	}{
		{
			"Success",
			nil,
			&options4.TagOptions{
				RootOptions: rootOpts,
			},
			&s3.GetBucketTaggingOutput{
				TagSet: []*s3.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
					{
						Key:   aws.String("foo2"),
						Value: aws.String("bar2"),
					},
				},
			},
			nil,
		},
		{
			"Failure",
			constants.ErrInjected,
			&options4.TagOptions{
				RootOptions: rootOpts,
			},
			nil,
			constants.ErrInjected,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(MockS3Client)
		mockS3.On("GetBucketTagging", mock.AnythingOfType("*s3.GetBucketTaggingInput")).Return(tc.GetBucketTaggingOutput, tc.getBucketTaggingErr)

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
		tags                []*s3.Tag
		putBucketTaggingErr error
		enableDryRun        bool
		enableAutoApprove   bool
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
			[]*s3.Tag{
				{
					Key:   aws.String("foo"),
					Value: aws.String("bar"),
				},
				{
					Key:   aws.String("foo2"),
					Value: aws.String("bar2"),
				},
			},
			nil,
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
			[]*s3.Tag{
				{
					Key:   aws.String("foo"),
					Value: aws.String("bar"),
				},
				{
					Key:   aws.String("foo2"),
					Value: aws.String("bar2"),
				},
			},
			nil,
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
			[]*s3.Tag{
				{
					Key:   aws.String("foo"),
					Value: aws.String("bar"),
				},
				{
					Key:   aws.String("foo2"),
					Value: aws.String("bar2"),
				},
			},
			constants.ErrInjected,
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
			[]*s3.Tag{
				{
					Key:   aws.String("foo"),
					Value: aws.String("bar"),
				},
				{
					Key:   aws.String("foo2"),
					Value: aws.String("bar2"),
				},
			},
			constants.ErrInjected,
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
			[]*s3.Tag{
				{
					Key:   aws.String("foo"),
					Value: aws.String("bar"),
				},
				{
					Key:   aws.String("foo2"),
					Value: aws.String("bar2"),
				},
			},
			constants.ErrInjected,
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

		mockS3 := new(MockS3Client)
		mockS3.On("PutBucketTagging", mock.AnythingOfType("*s3.PutBucketTaggingInput")).Return(&s3.PutBucketTaggingOutput{}, tc.putBucketTaggingErr)

		for _, v := range tc.tags {
			tc.TagOptions.TagsToAdd[*v.Key] = *v.Value
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
		deleteBucketTaggingErr error
		enableDryRun           bool
		enableAutoApprove      bool
		prompt.PromptRunner
	}{
		{
			"Success",
			nil,
			&options4.TagOptions{
				RootOptions: rootOpts,
			},
			nil,
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
			nil,
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
			constants.ErrInjected,
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
			nil,
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
			nil,
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

		mockS3 := new(MockS3Client)
		mockS3.On("DeleteBucketTagging", mock.AnythingOfType("*s3.DeleteBucketTaggingInput")).Return(&s3.DeleteBucketTaggingOutput{}, tc.deleteBucketTaggingErr)

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
		caseName string
		expected error
		*options6.BucketPolicyOptions
		getBucketPolicyErr error
	}{
		{
			"Success",
			nil,
			&options6.BucketPolicyOptions{
				RootOptions: rootOpts,
			},
			nil,
		},
		{
			"Failure",
			constants.ErrInjected,
			&options6.BucketPolicyOptions{
				RootOptions: rootOpts,
			},
			constants.ErrInjected,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(MockS3Client)
		mockS3.On("GetBucketPolicy", mock.AnythingOfType("*s3.GetBucketPolicyInput")).Return(&s3.GetBucketPolicyOutput{}, tc.getBucketPolicyErr)

		_, err := GetBucketPolicy(mockS3, tc.BucketPolicyOptions)
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
		putBucketPolicyErr error
		autoApprove        bool
		dryRun             bool
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
			nil,
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
			nil,
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
			constants.ErrInjected,
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
			nil,
			false,
			false,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		tc.DryRun = tc.dryRun
		tc.AutoApprove = tc.autoApprove

		mockS3 := new(MockS3Client)
		mockS3.On("PutBucketPolicy", mock.AnythingOfType("*s3.PutBucketPolicyInput")).Return(&s3.PutBucketPolicyOutput{}, tc.putBucketPolicyErr)

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
		*s3.GetBucketPolicyOutput
		getBucketPolicyErr error
	}{
		{
			"Success",
			nil,
			&options6.BucketPolicyOptions{
				RootOptions:         rootOpts,
				BucketPolicyContent: dummyBucketPolicyStr,
			},
			&s3.GetBucketPolicyOutput{
				Policy: aws.String(dummyBucketPolicyStr),
			},
			nil,
		},
		{
			"Failure",
			constants.ErrInjected,
			&options6.BucketPolicyOptions{
				RootOptions:         rootOpts,
				BucketPolicyContent: dummyBucketPolicyStr,
			},
			nil,
			constants.ErrInjected,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(MockS3Client)
		mockS3.On("GetBucketPolicy", mock.AnythingOfType("*s3.GetBucketPolicyInput")).Return(tc.GetBucketPolicyOutput, tc.getBucketPolicyErr)

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
		deleteBucketPolicyErr error
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
			nil,
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
			nil,
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
			constants.ErrInjected,
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
			constants.ErrInjected,
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

		mockS3 := new(MockS3Client)
		mockS3.On("DeleteBucketPolicy", mock.AnythingOfType("*s3.DeleteBucketPolicyInput")).Return(&s3.DeleteBucketPolicyOutput{}, tc.deleteBucketPolicyErr)

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
		getBucketAccelerationErr error
	}{
		{
			"Success",
			nil,
			&options5.TransferAccelerationOptions{
				RootOptions: rootOpts,
			},
			nil,
		},
		{
			"Failure",
			constants.ErrInjected,
			&options5.TransferAccelerationOptions{
				RootOptions: rootOpts,
			},
			constants.ErrInjected,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(MockS3Client)
		mockS3.On("GetBucketAccelerateConfiguration", mock.AnythingOfType("*s3.GetBucketAccelerateConfigurationInput")).Return(&s3.GetBucketAccelerateConfigurationOutput{}, tc.getBucketAccelerationErr)

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
		*s3.GetBucketAccelerateConfigurationOutput
		getBucketAccelerationErr error
		putBucketAccelerationErr error
		dryRun                   bool
		autoApprove              bool
		prompt.PromptRunner
	}{
		{
			"Success",
			nil,
			&options5.TransferAccelerationOptions{
				RootOptions:  rootOpts,
				DesiredState: "enabled",
			},
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Suspended"),
			},
			nil,
			nil,
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
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Enabled"),
			},
			nil,
			nil,
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
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Suspended"),
			},
			nil,
			nil,
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
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Suspended"),
			},
			nil,
			nil,
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
			nil,
			constants.ErrInjected,
			nil,
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
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Suspendedddd"),
			},
			nil,
			nil,
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
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Enabled"),
			},
			nil,
			constants.ErrInjected,
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
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Enabled"),
			},
			nil,
			nil,
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
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Enabled"),
			},
			nil,
			nil,
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

		mockS3 := new(MockS3Client)
		mockS3.On("GetBucketAccelerateConfiguration", mock.AnythingOfType("*s3.GetBucketAccelerateConfigurationInput")).Return(tc.GetBucketAccelerateConfigurationOutput, tc.getBucketAccelerationErr)
		mockS3.On("PutBucketAccelerateConfiguration", mock.AnythingOfType("*s3.PutBucketAccelerateConfigurationInput")).Return(&s3.PutBucketAccelerateConfigurationOutput{}, tc.putBucketAccelerationErr)

		err := SetTransferAcceleration(mockS3, tc.TransferAccelerationOptions, tc.PromptRunner, logging.GetLogger(tc.RootOptions))
		if tc.expected == nil {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}
