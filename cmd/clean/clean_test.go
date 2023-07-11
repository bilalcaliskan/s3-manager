//go:build e2e

package clean

import (
	"context"
	"testing"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

var (
	dummyListObjectsOutput = &s3.ListObjectsOutput{
		CommonPrefixes: nil,
		Contents: []*s3.Object{
			{
				ChecksumAlgorithm: nil,
				ETag:              nil,
				Key:               aws.String("foo.json"),
				LastModified:      nil,
				Owner:             nil,
				Size:              aws.Int64(500),
				StorageClass:      nil,
			},
			{
				ChecksumAlgorithm: nil,
				ETag:              nil,
				Key:               aws.String("bar.json"),
				LastModified:      nil,
				Owner:             nil,
				Size:              aws.Int64(600),
				StorageClass:      nil,
			},
		},
		Delimiter:      nil,
		EncodingType:   nil,
		IsTruncated:    nil,
		Marker:         nil,
		MaxKeys:        nil,
		Name:           nil,
		NextMarker:     nil,
		Prefix:         nil,
		RequestCharged: nil,
	}

	defaultListObjectsErr    error
	defaultListObjectsOutput = &s3.ListObjectsOutput{}

	defaultDeleteObjectErr    error
	defaultDeleteObjectOutput = &s3.DeleteObjectOutput{}
)

// Define a testdata struct to be used in your unit tests
type mockS3Client struct {
	s3iface.S3API
}

func (m *mockS3Client) DeleteObject(*s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return defaultDeleteObjectOutput, defaultDeleteObjectErr
}

func (m *mockS3Client) ListObjects(*s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	return defaultListObjectsOutput, defaultListObjectsErr
}

func TestExecuteCleanCmd(t *testing.T) {
	ctx := context.Background()
	CleanCmd.SetContext(ctx)

	rootOpts := options.GetMockedRootOptions()
	cases := []struct {
		caseName           string
		args               []string
		svc                s3iface.S3API
		shouldPass         bool
		listObjectsErr     error
		listObjectsOutput  *s3.ListObjectsOutput
		deleteObjectErr    error
		deleteObjectOutput *s3.DeleteObjectOutput
	}{
		{
			"Success",
			[]string{},
			&mockS3Client{},
			true,
			nil,
			&s3.ListObjectsOutput{},
			nil,
			&s3.DeleteObjectOutput{},
		},
		{
			"Failure caused by invalid 'sortBy' flag",
			[]string{"--sort-by=asldkfjalsdkf"},
			&mockS3Client{},
			false,
			nil,
			dummyListObjectsOutput,
			nil,
			&s3.DeleteObjectOutput{},
		},
		{
			"Failure caused by wrong size flags",
			[]string{"--max-size-mb=10", "--min-size-mb=20"},
			&mockS3Client{},
			false,
			nil,
			dummyListObjectsOutput,
			nil,
			&s3.DeleteObjectOutput{},
		},
		{
			"Failure caused by ListObjects error",
			[]string{},
			&mockS3Client{},
			false,
			constants.ErrInjected,
			&s3.ListObjectsOutput{},
			nil,
			&s3.DeleteObjectOutput{},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case '%s'", tc.caseName)

		defaultListObjectsErr = tc.listObjectsErr
		defaultListObjectsOutput = tc.listObjectsOutput

		defaultDeleteObjectErr = tc.deleteObjectErr
		defaultDeleteObjectOutput = tc.deleteObjectOutput

		CleanCmd.SetArgs(tc.args)
		CleanCmd.SetContext(context.WithValue(CleanCmd.Context(), options.S3SvcKey{}, tc.svc))
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
