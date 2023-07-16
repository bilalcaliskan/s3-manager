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

func TestSetBucketTags(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	cases := []struct {
		caseName string
		expected error
		*options4.TagOptions
		tags                []*s3.Tag
		putBucketTaggingErr error
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
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(MockS3Client)
		mockS3.On("PutBucketTagging", mock.AnythingOfType("*s3.PutBucketTaggingInput")).Return(&s3.PutBucketTaggingOutput{}, tc.putBucketTaggingErr)

		for _, v := range tc.tags {
			tc.TagOptions.TagsToAdd[*v.Key] = *v.Value
		}

		_, err := SetBucketTags(mockS3, tc.TagOptions)
		assert.Equal(t, tc.expected, err)
	}
}

func TestDeleteAllBucketTags(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	cases := []struct {
		caseName string
		expected error
		*options4.TagOptions
		deleteBucketTaggingErr error
	}{
		{
			"Success",
			nil,
			&options4.TagOptions{
				RootOptions: rootOpts,
			},
			nil,
		},
		{
			"Failure",
			constants.ErrInjected,
			&options4.TagOptions{
				RootOptions: rootOpts,
			},
			constants.ErrInjected,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(MockS3Client)
		mockS3.On("DeleteBucketTagging", mock.AnythingOfType("*s3.DeleteBucketTaggingInput")).Return(&s3.DeleteBucketTaggingOutput{}, tc.deleteBucketTaggingErr)

		_, err := DeleteAllBucketTags(mockS3, tc.TagOptions)
		assert.Equal(t, tc.expected, err)
	}
}

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

func TestSetBucketPolicy(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	cases := []struct {
		caseName string
		expected error
		*options6.BucketPolicyOptions
		putBucketPolicyErr error
	}{
		{
			"Success",
			nil,
			&options6.BucketPolicyOptions{
				RootOptions:         rootOpts,
				BucketPolicyContent: dummyBucketPolicyStr,
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
			constants.ErrInjected,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(MockS3Client)
		mockS3.On("PutBucketPolicy", mock.AnythingOfType("*s3.PutBucketPolicyInput")).Return(&s3.PutBucketPolicyOutput{}, tc.putBucketPolicyErr)

		_, err := SetBucketPolicy(mockS3, tc.BucketPolicyOptions)
		assert.Equal(t, tc.expected, err)
	}
}

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

func TestDeleteBucketPolicy(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	cases := []struct {
		caseName string
		expected error
		*options6.BucketPolicyOptions
		deleteBucketPolicyErr error
	}{
		{
			"Success",
			nil,
			&options6.BucketPolicyOptions{
				RootOptions:         rootOpts,
				BucketPolicyContent: dummyBucketPolicyStr,
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
			constants.ErrInjected,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(MockS3Client)
		mockS3.On("DeleteBucketPolicy", mock.AnythingOfType("*s3.DeleteBucketPolicyInput")).Return(&s3.DeleteBucketPolicyOutput{}, tc.deleteBucketPolicyErr)

		_, err := DeleteBucketPolicy(mockS3, tc.BucketPolicyOptions)
		assert.Equal(t, tc.expected, err)
	}
}

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
