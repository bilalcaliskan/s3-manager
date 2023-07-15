//go:build e2e

package show

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

type mockS3Client struct {
	mock.Mock
	s3iface.S3API
}

// GetBucketTagging mocks the GetBucketTagging method of s3iface.S3API
func (m *mockS3Client) GetBucketTagging(input *s3.GetBucketTaggingInput) (*s3.GetBucketTaggingOutput, error) {
	// Return the mocked output values using the `On` method of testify/mock
	args := m.Called(input)
	return args.Get(0).(*s3.GetBucketTaggingOutput), args.Error(1)
}

func TestExecuteShowCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()

	ctx := context.Background()
	ShowCmd.SetContext(ctx)

	cases := []struct {
		caseName               string
		args                   []string
		shouldPass             bool
		getBucketTaggingErr    error
		getBucketTaggingOutput *s3.GetBucketTaggingOutput
	}{
		{
			"Too many arguments",
			[]string{"enabled", "foo"},
			false,
			nil,
			&s3.GetBucketTaggingOutput{},
		},
		{
			"Success with empty TagSet",
			[]string{},
			true,
			nil,
			&s3.GetBucketTaggingOutput{},
		},
		{
			"Success with non-empty TagSet",
			[]string{},
			true,
			nil,
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
		},
		{
			"Failure",
			[]string{},
			false,
			errors.New("dummy error"),
			&s3.GetBucketTaggingOutput{},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(mockS3Client)
		mockS3.On("GetBucketTagging", mock.AnythingOfType("*s3.GetBucketTaggingInput")).Return(tc.getBucketTaggingOutput, tc.getBucketTaggingErr)

		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3SvcKey{}, mockS3))
		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))
		ShowCmd.SetArgs(tc.args)

		err := ShowCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}

	tagOpts.SetZeroValues()
}
