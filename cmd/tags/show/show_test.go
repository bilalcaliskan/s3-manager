//go:build e2e

package show

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/stretchr/testify/assert"
)

var (
	defaultGetBucketTaggingErr    error
	defaultGetBucketTaggingOutput = &s3.GetBucketTaggingOutput{}
)

type mockS3Client struct {
	s3iface.S3API
}

func (m *mockS3Client) GetBucketTagging(input *s3.GetBucketTaggingInput) (*s3.GetBucketTaggingOutput, error) {
	return defaultGetBucketTaggingOutput, defaultGetBucketTaggingErr
}

func TestExecuteShowCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	svc, err := internalaws.CreateAwsService(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	ctx := context.Background()
	ShowCmd.SetContext(ctx)

	cases := []struct {
		caseName                      string
		args                          []string
		shouldPass                    bool
		svc                           s3iface.S3API
		defaultGetBucketTaggingErr    error
		defaultGetBucketTaggingOutput *s3.GetBucketTaggingOutput
	}{
		{
			"Too many arguments",
			[]string{"enabled", "foo"},
			false,
			svc,
			nil,
			&s3.GetBucketTaggingOutput{},
		},
		{
			"Success with empty TagSet",
			[]string{},
			true,
			&mockS3Client{},
			nil,
			&s3.GetBucketTaggingOutput{},
		},
		{
			"Success with non-empty TagSet",
			[]string{},
			true,
			&mockS3Client{},
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
			&mockS3Client{},
			errors.New("dummy error"),
			&s3.GetBucketTaggingOutput{},
		},
	}

	for _, tc := range cases {
		defaultGetBucketTaggingErr = tc.defaultGetBucketTaggingErr
		defaultGetBucketTaggingOutput = tc.defaultGetBucketTaggingOutput

		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3SvcKey{}, tc.svc))
		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))
		ShowCmd.SetArgs(tc.args)

		err = ShowCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}

	rootOpts.SetZeroValues()
	tagOpts.SetZeroValues()
}
