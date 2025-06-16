//go:build e2e

package show

import (
	"context"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/pkg/aws/types"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

func TestExecuteShowCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	ctx := context.Background()
	ShowCmd.SetContext(ctx)

	cases := []struct {
		caseName            string
		args                []string
		shouldPass          bool
		getBucketPolicyFunc func(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error)
	}{
		{
			"Too many arguments",
			[]string{"enabled", "foo"},
			false,
			nil,
		},
		{
			"Success",
			[]string{},
			true,
			func(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
				return &s3.GetBucketPolicyOutput{
					Policy: aws.String("{}"),
				}, nil
			},
		},
		{
			"Json failure",
			[]string{},
			false,
			func(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
				return &s3.GetBucketPolicyOutput{
					Policy: aws.String(""),
				}, nil
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.GetBucketPolicyAPI = tc.getBucketPolicyFunc

		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3ClientKey{}, mockS3))
		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))
		ShowCmd.SetArgs(tc.args)

		err := ShowCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		bucketPolicyOpts.SetZeroValues()
	}
}
