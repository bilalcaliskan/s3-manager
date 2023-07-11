//go:build e2e

package show

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/stretchr/testify/assert"
)

var (
	defaultGetBucketPolicyOutput = &s3.GetBucketPolicyOutput{
		Policy: aws.String("{}"),
	}
	defaultGetBucketPolicyErr error
)

type mockS3Client struct {
	s3iface.S3API
}

func (m *mockS3Client) GetBucketPolicy(input *s3.GetBucketPolicyInput) (*s3.GetBucketPolicyOutput, error) {
	return defaultGetBucketPolicyOutput, defaultGetBucketPolicyErr
}

func TestExecuteShowCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	ctx := context.Background()
	ShowCmd.SetContext(ctx)

	svc, err := internalaws.CreateAwsService(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	cases := []struct {
		caseName              string
		args                  []string
		svc                   s3iface.S3API
		shouldPass            bool
		getBucketPolicyErr    error
		getBucketPolicyOutput *s3.GetBucketPolicyOutput
	}{
		{
			"Too many arguments",
			[]string{"enabled", "foo"},
			svc,
			false,
			nil,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String("{}"),
			},
		},
		{
			"No argument",
			[]string{},
			svc,
			false,
			nil,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String("{}"),
			},
		},
		{
			"Success",
			[]string{},
			&mockS3Client{},
			true,
			nil,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String("{}"),
			},
		},
		{
			"Json failure",
			[]string{},
			&mockS3Client{},
			false,
			nil,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String(""),
			},
		},
	}

	for _, tc := range cases {
		defaultGetBucketPolicyErr = tc.getBucketPolicyErr
		defaultGetBucketPolicyOutput = tc.getBucketPolicyOutput

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
	bucketPolicyOpts.SetZeroValues()
}
