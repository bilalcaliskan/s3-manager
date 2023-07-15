//go:build e2e

package clean

import (
	"context"
	"testing"

	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"

	"github.com/stretchr/testify/mock"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

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
)

func TestExecuteCleanCmd(t *testing.T) {
	ctx := context.Background()
	CleanCmd.SetContext(ctx)

	rootOpts := options.GetMockedRootOptions()
	cases := []struct {
		caseName           string
		args               []string
		shouldPass         bool
		listObjectsErr     error
		listObjectsOutput  *s3.ListObjectsOutput
		deleteObjectErr    error
		deleteObjectOutput *s3.DeleteObjectOutput
	}{
		{
			"Success",
			[]string{},
			true,
			nil,
			&s3.ListObjectsOutput{},
			nil,
			&s3.DeleteObjectOutput{},
		},
		{
			"Failure caused by invalid 'sortBy' flag",
			[]string{"--sort-by=asldkfjalsdkf"},
			false,
			nil,
			dummyListObjectsOutput,
			nil,
			&s3.DeleteObjectOutput{},
		},
		{
			"Failure caused by wrong size flags",
			[]string{"--max-size-mb=10", "--min-size-mb=20"},
			false,
			nil,
			dummyListObjectsOutput,
			nil,
			&s3.DeleteObjectOutput{},
		},
		{
			"Failure caused by ListObjects error",
			[]string{},
			false,
			constants.ErrInjected,
			&s3.ListObjectsOutput{},
			nil,
			&s3.DeleteObjectOutput{},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case '%s'", tc.caseName)

		mockS3 := new(internalaws.MockS3Client)
		mockS3.On("DeleteObject", mock.AnythingOfType("*s3.DeleteObjectInput")).Return(tc.deleteObjectOutput, tc.deleteObjectErr)
		mockS3.On("ListObjects", mock.AnythingOfType("*s3.ListObjectsInput")).Return(tc.listObjectsOutput, tc.listObjectsErr)

		CleanCmd.SetArgs(tc.args)
		CleanCmd.SetContext(context.WithValue(CleanCmd.Context(), options.S3SvcKey{}, mockS3))
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
