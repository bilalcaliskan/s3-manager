//go:build e2e

package remove

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/aws/types"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

var policyStr = `
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

func TestExecuteRemoveCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()

	ctx := context.Background()
	RemoveCmd.SetContext(ctx)
	cases := []struct {
		caseName               string
		args                   []string
		shouldPass             bool
		getBucketPolicyFunc    func(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error)
		deleteBucketPolicyFunc func(ctx context.Context, params *s3.DeleteBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketPolicyOutput, error)
		prompt.PromptRunner
		dryRun      bool
		autoApprove bool
	}{
		{
			"Too many arguments",
			[]string{"enabled", "foo"},
			false,
			nil,
			nil,
			nil,
			false,
			false,
		},
		{
			"Success",
			[]string{},
			true,
			func(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
				return &s3.GetBucketPolicyOutput{
					Policy: aws.String(policyStr),
				}, nil
			},
			func(ctx context.Context, params *s3.DeleteBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketPolicyOutput, error) {
				return &s3.DeleteBucketPolicyOutput{}, nil
			},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Success with dry run",
			[]string{},
			true,
			func(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
				return &s3.GetBucketPolicyOutput{
					Policy: aws.String(policyStr),
				}, nil
			},
			func(ctx context.Context, params *s3.DeleteBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketPolicyOutput, error) {
				return &s3.DeleteBucketPolicyOutput{}, nil
			},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			true,
			false,
		},
		{
			"Success with auto approve",
			[]string{},
			true,
			func(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
				return &s3.GetBucketPolicyOutput{
					Policy: aws.String(policyStr),
				}, nil
			},
			func(ctx context.Context, params *s3.DeleteBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketPolicyOutput, error) {
				return &s3.DeleteBucketPolicyOutput{}, nil
			},
			nil,
			false,
			true,
		},
		{
			"Failure caused by delete error",
			[]string{},
			false,
			func(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
				return &s3.GetBucketPolicyOutput{
					Policy: aws.String(policyStr),
				}, nil
			},
			func(ctx context.Context, params *s3.DeleteBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketPolicyOutput, error) {
				return nil, constants.ErrInjected
			},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by get bucket policy error",
			[]string{},
			false,
			func(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
				return nil, constants.ErrInjected
			},
			func(ctx context.Context, params *s3.DeleteBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketPolicyOutput, error) {
				return nil, constants.ErrInjected
			},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by user terminated process",
			[]string{},
			false,
			func(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
				return nil, constants.ErrInjected
			},
			nil,
			prompt.PromptMock{
				Msg: "n",
				Err: constants.ErrInjected,
			},
			false,
			false,
		},
		{
			"Failure caused by prompt error",
			[]string{},
			false,
			func(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
				return nil, constants.ErrInjected
			},
			nil,
			prompt.PromptMock{
				Msg: "nasdfadf",
				Err: constants.ErrInjected,
			},
			false,
			false,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.GetBucketPolicyAPI = tc.getBucketPolicyFunc
		mockS3.DeleteBucketPolicyAPI = tc.deleteBucketPolicyFunc

		rootOpts.DryRun = tc.dryRun
		rootOpts.AutoApprove = tc.autoApprove

		RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.S3ClientKey{}, mockS3))
		RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.OptsKey{}, rootOpts))
		RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.ConfirmRunnerKey{}, tc.PromptRunner))
		RemoveCmd.SetArgs(tc.args)

		err := RemoveCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		bucketPolicyOpts.SetZeroValues()
	}
}
