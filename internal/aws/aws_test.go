//go:build unit

package aws

import (
	"errors"
	"os"
	"testing"
	"time"

	options6 "github.com/bilalcaliskan/s3-manager/cmd/bucketpolicy/options"
	options2 "github.com/bilalcaliskan/s3-manager/cmd/search/options"
	options5 "github.com/bilalcaliskan/s3-manager/cmd/transferacceleration/options"

	options4 "github.com/bilalcaliskan/s3-manager/cmd/tags/options"

	options3 "github.com/bilalcaliskan/s3-manager/cmd/versioning/options"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/internal/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

var (
	injectedErr            = errors.New("injected error")
	defaultListObjectsErr  error
	defaultGetObjectErr    error
	defaultDeleteObjectErr error
	fileNamePrefix         string
	/*defaultListObjectsOutput = &s3.ListObjectsOutput{
		Name:        aws.String(""),
		Marker:      aws.String(""),
		MaxKeys:     aws.Int64(1000),
		Prefix:      aws.String(""),
		IsTruncated: aws.Bool(false),
	}*/
	defaultListObjectsOutput = &s3.ListObjectsOutput{}
	/*defaultDeleteObjectOutput = &s3.DeleteObjectOutput{
		DeleteMarker:   nil,
		RequestCharged: nil,
		VersionId:      nil,
	}*/
	defaultDeleteObjectOutput        = &s3.DeleteObjectOutput{}
	defaultGetBucketVersioningOutput = &s3.GetBucketVersioningOutput{
		Status: aws.String("Enabled"),
	}
	defaultGetBucketVersioningErr      error
	defaultPutBucketVersioningOutput   = &s3.PutBucketVersioningOutput{}
	defaultPutBucketVersioningErr      error
	defaultGetBucketTaggingErr         error
	defaultGetBucketTaggingOutput      = &s3.GetBucketTaggingOutput{}
	defaultPutBucketTaggingErr         error
	defaultPutBucketTaggingOutput      = &s3.PutBucketTaggingOutput{}
	defaultDeleteBucketTaggingErr      error
	defaultDeleteBucketTaggingOutput   = &s3.DeleteBucketTaggingOutput{}
	defaultGetBucketAccelerationOutput = &s3.GetBucketAccelerateConfigurationOutput{
		Status: aws.String("Enabled"),
	}
	defaultGetBucketAccelerationErr    error
	defaultPutBucketAccelerationOutput = &s3.PutBucketAccelerateConfigurationOutput{}
	defaultPutBucketAccelerationErr    error
	defaultGetBucketPolicyErr          error
	defaultGetBucketPolicyOutput       = &s3.GetBucketPolicyOutput{}
	defaultPutBucketPolicyErr          error
	defaultPutBucketPolicyOutput       = &s3.PutBucketPolicyOutput{}
	defaultDeleteBucketPolicyErr       error
	defaultDeleteBucketPolicyOutput    = &s3.DeleteBucketPolicyOutput{}
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
	bytes, err := os.Open(*input.Key)
	if err != nil {
		return nil, err
	}

	return &s3.GetObjectOutput{
		AcceptRanges:  aws.String("bytes"),
		Body:          bytes,
		ContentLength: aws.Int64(1000),
		ContentType:   aws.String("text/plain"),
		ETag:          aws.String("d73a503d212d9279e6b2ed8ac6bb81f3"),
	}, defaultGetObjectErr
}

func (m *mockS3Client) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return defaultDeleteObjectOutput, defaultDeleteObjectErr
}

func (m *mockS3Client) GetBucketAccelerateConfiguration(input *s3.GetBucketAccelerateConfigurationInput) (*s3.GetBucketAccelerateConfigurationOutput, error) {
	return defaultGetBucketAccelerationOutput, defaultGetBucketAccelerationErr
}

func (m *mockS3Client) PutBucketAccelerateConfiguration(input *s3.PutBucketAccelerateConfigurationInput) (*s3.PutBucketAccelerateConfigurationOutput, error) {
	return defaultPutBucketAccelerationOutput, defaultPutBucketAccelerationErr
}

func (m *mockS3Client) GetBucketVersioning(input *s3.GetBucketVersioningInput) (*s3.GetBucketVersioningOutput, error) {
	return defaultGetBucketVersioningOutput, defaultGetBucketVersioningErr
}

func (m *mockS3Client) PutBucketVersioning(input *s3.PutBucketVersioningInput) (*s3.PutBucketVersioningOutput, error) {
	return defaultPutBucketVersioningOutput, defaultPutBucketVersioningErr
}

func (m *mockS3Client) GetBucketTagging(input *s3.GetBucketTaggingInput) (*s3.GetBucketTaggingOutput, error) {
	return defaultGetBucketTaggingOutput, defaultGetBucketTaggingErr
}

func (m *mockS3Client) PutBucketTagging(input *s3.PutBucketTaggingInput) (*s3.PutBucketTaggingOutput, error) {
	return defaultPutBucketTaggingOutput, defaultPutBucketTaggingErr
}

func (m *mockS3Client) DeleteBucketTagging(input *s3.DeleteBucketTaggingInput) (*s3.DeleteBucketTaggingOutput, error) {
	return defaultDeleteBucketTaggingOutput, defaultDeleteBucketTaggingErr
}

func (m *mockS3Client) GetBucketPolicy(input *s3.GetBucketPolicyInput) (*s3.GetBucketPolicyOutput, error) {
	return defaultGetBucketPolicyOutput, defaultGetBucketPolicyErr
}

func (m *mockS3Client) PutBucketPolicy(input *s3.PutBucketPolicyInput) (*s3.PutBucketPolicyOutput, error) {
	return defaultPutBucketPolicyOutput, defaultPutBucketPolicyErr
}

func (m *mockS3Client) DeleteBucketPolicy(input *s3.DeleteBucketPolicyInput) (*s3.DeleteBucketPolicyOutput, error) {
	return defaultDeleteBucketPolicyOutput, defaultDeleteBucketPolicyErr
}

func TestGetAllFiles(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	cases := []struct {
		caseName          string
		expected          error
		listObjectsErr    error
		listObjectsOutput *s3.ListObjectsOutput
	}{
		{"Success with non-empty file list",
			nil, nil,
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
		},
		{"Failure caused by List objects error",
			injectedErr, injectedErr,
			nil,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		defaultListObjectsErr = tc.listObjectsErr
		defaultListObjectsOutput = tc.listObjectsOutput

		mockSvc := &mockS3Client{}
		assert.NotNil(t, mockSvc)

		_, err := GetAllFiles(mockSvc, rootOpts, "")
		assert.Equal(t, tc.expected, err)
	}
}

func TestDeleteFiles(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	cases := []struct {
		caseName        string
		expected        error
		deleteObjectErr error
		dryRun          bool
		objects         []*s3.Object
	}{
		{"Success with non-empty file list",
			nil, nil, false,
			[]*s3.Object{
				{
					ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
					Key:          aws.String("../../testdata/file1.txt"),
					StorageClass: aws.String("STANDARD"),
					Size:         aws.Int64(500),
					LastModified: aws.Time(time.Now().Add(-5 * time.Hour)),
				},
				{
					ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
					Key:          aws.String("../../testdata/file2.txt"),
					StorageClass: aws.String("STANDARD"),
					Size:         aws.Int64(1000),
					LastModified: aws.Time(time.Now().Add(-2 * time.Hour)),
				},
				{
					ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
					Key:          aws.String("../../testdata/file3.txt"),
					StorageClass: aws.String("STANDARD"),
					Size:         aws.Int64(1500),
					LastModified: aws.Time(time.Now().Add(-10 * time.Hour)),
				},
			},
		},
		{"Failure caused by delete object err",
			injectedErr, injectedErr, false,
			[]*s3.Object{
				{
					ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5449d"),
					Key:          aws.String("../../testdata/file1.txt"),
					StorageClass: aws.String("STANDARD"),
					Size:         aws.Int64(500),
					LastModified: aws.Time(time.Now().Add(-5 * time.Hour)),
				},
				{
					ETag:         aws.String("03c0fe42b7efa3470fc99037a8e54122"),
					Key:          aws.String("../../testdata/file2.txt"),
					StorageClass: aws.String("STANDARD"),
					Size:         aws.Int64(1000),
					LastModified: aws.Time(time.Now().Add(-2 * time.Hour)),
				},
				{
					ETag:         aws.String("03c0fe42b7efa3470fc99037a8e5443d"),
					Key:          aws.String("../../testdata/file3.txt"),
					StorageClass: aws.String("STANDARD"),
					Size:         aws.Int64(1500),
					LastModified: aws.Time(time.Now().Add(-10 * time.Hour)),
				},
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		defaultDeleteObjectErr = tc.deleteObjectErr

		mockSvc := &mockS3Client{}
		assert.NotNil(t, mockSvc)

		assert.Equal(t, tc.expected, DeleteFiles(mockSvc, "thisisdemobucket", tc.objects, tc.dryRun, logging.GetLogger(rootOpts)))
	}
}

func TestCreateAwsService(t *testing.T) {
	cases := []struct {
		caseName   string
		opts       *options.RootOptions
		shouldPass bool
	}{
		{"Success",
			&options.RootOptions{
				AccessKey:  "thisisaccesskey",
				SecretKey:  "thisissecretkey",
				BucketName: "thisisbucketname",
				Region:     "thisisregion",
			}, true,
		},
		{"Failure caused by missing required field",
			&options.RootOptions{
				AccessKey:  "thisisaccesskey",
				SecretKey:  "thisissecretkey",
				BucketName: "thisisbucketname",
				Region:     "",
			}, false,
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
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	cases := []struct {
		caseName          string
		searchOpts        *options2.SearchOptions
		shouldPass        bool
		listObjectsErr    error
		listObjectsOutput *s3.ListObjectsOutput
		getObjectErr      error
		matchCount        int
	}{
		{"Success with specific text",
			&options2.SearchOptions{
				Text:        "pvRRTaigmb",
				FileName:    "",
				RootOptions: nil,
			}, true, nil,
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
			}, nil, 2,
		},
		{"Success with file name regex",
			&options2.SearchOptions{
				Text:        "pvRRTaigmb",
				FileName:    "file2.*.",
				RootOptions: nil,
			}, true, nil,
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
			}, nil, 1,
		},
		{"Failure caused by list objects error",
			&options2.SearchOptions{
				Text:        "pvRRTaigmb",
				FileName:    "",
				RootOptions: nil,
			}, false, errors.New("injected error"),
			nil, nil, 0,
		},
		{"Failure caused by get object error",
			&options2.SearchOptions{
				Text:        "pvRRTaigmb",
				FileName:    "",
				RootOptions: nil,
			}, false, nil,
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
			}, errors.New("injected error"), 0,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		tc.searchOpts.RootOptions = rootOpts
		defaultListObjectsErr = tc.listObjectsErr
		defaultListObjectsOutput = tc.listObjectsOutput
		defaultGetObjectErr = tc.getObjectErr

		mockSvc := &mockS3Client{}
		assert.NotNil(t, mockSvc)

		res, err := SearchString(mockSvc, tc.searchOpts)

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		assert.Equal(t, tc.matchCount, len(res))
	}
}

func TestSetBucketVersioning(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	cases := []struct {
		caseName string
		*options3.VersioningOptions
		*s3.GetBucketVersioningOutput
		getBucketVersioningErr error
		putBucketVersioningErr error
		expected               error
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
			}, nil, nil, nil,
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
			}, nil, nil, nil,
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
			}, nil, nil, nil,
		},
		{
			"Failure caused by get versioning error",
			&options3.VersioningOptions{
				ActualState:  "",
				DesiredState: "disabled",
				RootOptions:  rootOpts,
			},
			nil, injectedErr, nil, injectedErr,
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
			}, nil, nil, errors.New("unknown status 'Enableddddd' returned from AWS SDK"),
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		defaultGetBucketVersioningOutput = tc.GetBucketVersioningOutput
		defaultGetBucketVersioningErr = tc.getBucketVersioningErr
		defaultPutBucketVersioningErr = tc.putBucketVersioningErr

		mockSvc := &mockS3Client{}
		assert.NotNil(t, mockSvc)

		err := SetBucketVersioning(mockSvc, tc.VersioningOptions, logging.GetLogger(tc.VersioningOptions.RootOptions))
		assert.Equal(t, tc.expected, err)
	}
}

func TestGetBucketVersioning(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	cases := []struct {
		caseName string
		expected error
		*s3.GetBucketVersioningOutput
		getBucketVersioningErr error
	}{
		{
			"Success", nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enableddddd"),
			}, nil,
		},
		{
			"Failure", injectedErr,
			nil, injectedErr,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		defaultGetBucketVersioningOutput = tc.GetBucketVersioningOutput
		defaultGetBucketVersioningErr = tc.getBucketVersioningErr

		mockSvc := &mockS3Client{}
		assert.NotNil(t, mockSvc)

		_, err := GetBucketVersioning(mockSvc, rootOpts)
		assert.Equal(t, tc.expected, err)
	}
}

func TestGetBucketTags(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	cases := []struct {
		caseName string
		expected error
		*options4.TagOptions
		*s3.GetBucketTaggingOutput
		getBucketTaggingErr error
	}{
		{
			"Success", nil,
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
			}, nil,
		},
		{
			"Failure", injectedErr,
			&options4.TagOptions{
				RootOptions: rootOpts,
			},
			nil, injectedErr,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		defaultGetBucketTaggingOutput = tc.GetBucketTaggingOutput
		defaultGetBucketTaggingErr = tc.getBucketTaggingErr

		mockSvc := &mockS3Client{}
		assert.NotNil(t, mockSvc)

		_, err := GetBucketTags(mockSvc, tc.TagOptions)
		assert.Equal(t, tc.expected, err)
	}
}

/*func TestSetBucketTags(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	cases := []struct {
		caseName string
		expected error
		*options4.TagOptions
		*s3.GetBucketTaggingOutput
		getBucketTaggingErr error
	}{
		{
			"Success", nil,
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
			}, nil,
		},
		{
			"Failure", injectedErr,
			&options4.TagOptions{
				RootOptions: rootOpts,
			},
			nil, injectedErr,
		},
	}
}*/

func TestPutBucketTaggingSuccess(t *testing.T) {
	tagOpts := options4.GetTagOptions()
	defer func() {
		tagOpts.SetZeroValues()
	}()
	rootOpts := options.GetRootOptions()
	rootOpts.Region = "us-east-1"
	tagOpts.RootOptions = rootOpts

	mockSvc := &mockS3Client{}
	var tags []*s3.Tag
	tags = append(tags, &s3.Tag{Key: aws.String("foo"), Value: aws.String("bar")})

	tagOpts.TagsToAdd = make(map[string]string)

	for _, v := range tags {
		tagOpts.TagsToAdd[*v.Key] = *v.Value
	}

	defaultGetBucketTaggingOutput = &s3.GetBucketTaggingOutput{TagSet: tags}
	defaultGetBucketTaggingErr = nil

	defaultPutBucketTaggingErr = nil

	_, err := SetBucketTags(mockSvc, tagOpts)
	assert.Nil(t, err)
}

func TestPutBucketTaggingFailure(t *testing.T) {
	tagOpts := options4.GetTagOptions()
	defer func() {
		tagOpts.SetZeroValues()
	}()
	rootOpts := options.GetRootOptions()
	rootOpts.Region = "us-east-1"
	tagOpts.RootOptions = rootOpts

	mockSvc := &mockS3Client{}
	var tags []*s3.Tag
	tags = append(tags, &s3.Tag{Key: aws.String("foo"), Value: aws.String("bar")})
	for _, v := range tags {
		tagOpts.TagsToAdd[*v.Key] = *v.Value
	}

	defaultGetBucketTaggingOutput = &s3.GetBucketTaggingOutput{TagSet: tags}
	defaultGetBucketTaggingErr = nil

	defaultPutBucketTaggingErr = errors.New("dummy error")

	_, err := SetBucketTags(mockSvc, tagOpts)
	assert.NotNil(t, err)
}

func TestDeleteBucketTaggingSuccess(t *testing.T) {
	tagOpts := options4.GetTagOptions()
	defer func() {
		tagOpts.SetZeroValues()
	}()
	rootOpts := options.GetRootOptions()
	rootOpts.Region = "us-east-1"
	tagOpts.RootOptions = rootOpts

	defaultDeleteBucketTaggingErr = nil
	_, err := DeleteAllBucketTags(&mockS3Client{}, tagOpts)
	assert.Nil(t, err)
}

func TestGetTransferAcceleration(t *testing.T) {
	taOpts := options5.GetTransferAccelerationOptions()
	defer func() {
		taOpts.SetZeroValues()
	}()
	rootOpts := options.GetRootOptions()
	rootOpts.Region = "us-east-1"
	taOpts.RootOptions = rootOpts

	defaultGetBucketAccelerationErr = nil

	res, err := GetTransferAcceleration(&mockS3Client{}, taOpts)
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestSetTransferAcceleration(t *testing.T) {
	taOpts := options5.GetTransferAccelerationOptions()
	defer func() {
		taOpts.SetZeroValues()
	}()
	rootOpts := options.GetRootOptions()
	rootOpts.Region = "us-east-1"
	taOpts.RootOptions = rootOpts

	taOpts.DesiredState = "enabled"

	defaultGetBucketAccelerationErr = nil
	defaultPutBucketAccelerationErr = nil

	err := SetTransferAcceleration(&mockS3Client{}, taOpts, logging.GetLogger(taOpts.RootOptions))
	assert.Nil(t, err)
}

func TestSetTransferAcceleration2(t *testing.T) {
	taOpts := options5.GetTransferAccelerationOptions()
	defer func() {
		taOpts.SetZeroValues()
	}()
	rootOpts := options.GetRootOptions()
	rootOpts.Region = "us-east-1"
	taOpts.RootOptions = rootOpts

	taOpts.DesiredState = "disabled"

	defaultGetBucketAccelerationErr = nil
	defaultPutBucketAccelerationErr = nil

	err := SetTransferAcceleration(&mockS3Client{}, taOpts, logging.GetLogger(taOpts.RootOptions))
	assert.Nil(t, err)
}

func TestSetTransferAcceleration3(t *testing.T) {
	taOpts := options5.GetTransferAccelerationOptions()
	defer func() {
		taOpts.SetZeroValues()
	}()
	rootOpts := options.GetRootOptions()
	rootOpts.Region = "us-east-1"
	taOpts.RootOptions = rootOpts

	taOpts.DesiredState = "disabled"

	defaultGetBucketAccelerationErr = errors.New("dummy error")
	defaultPutBucketAccelerationErr = nil

	err := SetTransferAcceleration(&mockS3Client{}, taOpts, logging.GetLogger(taOpts.RootOptions))
	assert.NotNil(t, err)
}

func TestSetTransferAcceleration4(t *testing.T) {
	taOpts := options5.GetTransferAccelerationOptions()
	defer func() {
		taOpts.SetZeroValues()
	}()
	rootOpts := options.GetRootOptions()
	rootOpts.Region = "us-east-1"
	taOpts.RootOptions = rootOpts

	taOpts.DesiredState = "enabled"

	defaultGetBucketAccelerationErr = nil
	defaultGetBucketAccelerationOutput.Status = aws.String("Suspended")
	defaultPutBucketAccelerationErr = nil

	err := SetTransferAcceleration(&mockS3Client{}, taOpts, logging.GetLogger(taOpts.RootOptions))
	assert.Nil(t, err)
}

func TestSetTransferAcceleration5(t *testing.T) {
	taOpts := options5.GetTransferAccelerationOptions()
	defer func() {
		taOpts.SetZeroValues()
	}()
	rootOpts := options.GetRootOptions()
	rootOpts.Region = "us-east-1"
	taOpts.RootOptions = rootOpts

	taOpts.DesiredState = "enabled"

	defaultGetBucketAccelerationErr = nil
	defaultGetBucketAccelerationOutput.Status = aws.String("Suspendedddd")
	defaultPutBucketAccelerationErr = nil

	err := SetTransferAcceleration(&mockS3Client{}, taOpts, logging.GetLogger(taOpts.RootOptions))
	assert.NotNil(t, err)
}

func TestSetTransferAcceleration6(t *testing.T) {
	taOpts := options5.GetTransferAccelerationOptions()
	defer func() {
		taOpts.SetZeroValues()
	}()
	rootOpts := options.GetRootOptions()
	rootOpts.Region = "us-east-1"
	taOpts.RootOptions = rootOpts

	taOpts.DesiredState = "enabled"

	defaultGetBucketAccelerationErr = nil
	defaultGetBucketAccelerationOutput.Status = aws.String("Suspended")
	defaultPutBucketAccelerationErr = errors.New("dummy error")

	err := SetTransferAcceleration(&mockS3Client{}, taOpts, logging.GetLogger(taOpts.RootOptions))
	assert.NotNil(t, err)
}

func TestGetBucketPolicy(t *testing.T) {
	bpOpts := options6.GetBucketPolicyOptions()
	defer func() {
		bpOpts.SetZeroValues()
	}()
	rootOpts := options.GetRootOptions()
	rootOpts.Region = "us-east-1"
	bpOpts.RootOptions = rootOpts

	res, err := GetBucketPolicy(&mockS3Client{}, bpOpts)
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestGetBucketPolicyStringSuccess(t *testing.T) {
	bpOpts := options6.GetBucketPolicyOptions()
	defer func() {
		bpOpts.SetZeroValues()
	}()

	rootOpts := options.GetRootOptions()
	rootOpts.Region = "us-east-1"
	bpOpts.RootOptions = rootOpts

	defaultGetBucketPolicyErr = nil
	defaultGetBucketPolicyOutput = &s3.GetBucketPolicyOutput{Policy: aws.String(`
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
`)}

	res, err := GetBucketPolicyString(&mockS3Client{}, bpOpts)
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestGetBucketPolicyStringError(t *testing.T) {
	bpOpts := options6.GetBucketPolicyOptions()
	defer func() {
		bpOpts.SetZeroValues()
	}()

	rootOpts := options.GetRootOptions()
	rootOpts.Region = "us-east-1"
	bpOpts.RootOptions = rootOpts

	defaultGetBucketPolicyErr = errors.New("injected error")
	defaultGetBucketPolicyOutput = nil

	res, err := GetBucketPolicyString(&mockS3Client{}, bpOpts)
	assert.Equal(t, "", res)
	assert.NotNil(t, err)
}

func TestSetBucketPolicy(t *testing.T) {
	bpOpts := options6.GetBucketPolicyOptions()
	defer func() {
		bpOpts.SetZeroValues()
	}()
	rootOpts := options.GetRootOptions()
	rootOpts.Region = "us-east-1"
	bpOpts.RootOptions = rootOpts

	res, err := SetBucketPolicy(&mockS3Client{}, bpOpts)
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestDeleteBucketPolicy(t *testing.T) {
	bpOpts := options6.GetBucketPolicyOptions()
	defer func() {
		bpOpts.SetZeroValues()
	}()
	rootOpts := options.GetRootOptions()
	rootOpts.Region = "us-east-1"
	bpOpts.RootOptions = rootOpts

	res, err := DeleteBucketPolicy(&mockS3Client{}, bpOpts)
	assert.NotNil(t, res)
	assert.Nil(t, err)
}
