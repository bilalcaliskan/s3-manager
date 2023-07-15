//go:build e2e

package remove

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"
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

type promptMock struct {
	msg string
	err error
}

func (p promptMock) Run() (string, error) {
	return p.msg, p.err
}

func TestExecuteRemoveCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()

	ctx := context.Background()
	RemoveCmd.SetContext(ctx)
	cases := []struct {
		caseName                 string
		args                     []string
		shouldPass               bool
		getBucketPolicyOutput    *s3.GetBucketPolicyOutput
		getBucketPolicyErr       error
		deleteBucketPolicyErr    error
		deleteBucketPolicyOutput *s3.DeleteBucketPolicyOutput
		promptMock               *promptMock
		dryRun                   bool
		autoApprove              bool
	}{
		{
			"Too many arguments",
			[]string{"enabled", "foo"},
			false,
			&s3.GetBucketPolicyOutput{},
			nil,
			nil,
			&s3.DeleteBucketPolicyOutput{},
			nil,
			false,
			false,
		},
		{
			"Success",
			[]string{},
			true,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String(policyStr),
			},
			nil,
			nil,
			&s3.DeleteBucketPolicyOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
		{
			"Success with dry run",
			[]string{},
			true,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String(policyStr),
			},
			nil,
			nil,
			&s3.DeleteBucketPolicyOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			true,
			false,
		},
		{
			"Success with auto approve",
			[]string{},
			true,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String(policyStr),
			},
			nil,
			nil,
			&s3.DeleteBucketPolicyOutput{},
			nil,
			false,
			true,
		},
		{
			"Failure",
			[]string{},
			false,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String(policyStr),
			},
			nil,
			constants.ErrInjected,
			&s3.DeleteBucketPolicyOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by get bucket policy error",
			[]string{},
			false,
			&s3.GetBucketPolicyOutput{
				Policy: nil,
			},
			constants.ErrInjected,
			nil,
			&s3.DeleteBucketPolicyOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by user terminated process",
			[]string{},
			false,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String(policyStr),
			},
			nil,
			nil,
			&s3.DeleteBucketPolicyOutput{},
			&promptMock{
				msg: "n",
				err: constants.ErrInjected,
			},
			false,
			false,
		},
		{
			"Failure caused by prompt error",
			[]string{},
			false,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String(policyStr),
			},
			nil,
			nil,
			&s3.DeleteBucketPolicyOutput{},
			&promptMock{
				msg: "nasdfadf",
				err: constants.ErrInjected,
			},
			false,
			false,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(internalaws.MockS3Client)
		mockS3.On("GetBucketPolicy", mock.AnythingOfType("*s3.GetBucketPolicyInput")).Return(tc.getBucketPolicyOutput, tc.getBucketPolicyErr)
		mockS3.On("DeleteBucketPolicy", mock.AnythingOfType("*s3.DeleteBucketPolicyInput")).Return(tc.deleteBucketPolicyOutput, tc.deleteBucketPolicyErr)

		rootOpts.DryRun = tc.dryRun
		rootOpts.AutoApprove = tc.autoApprove

		RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.S3SvcKey{}, mockS3))
		RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.OptsKey{}, rootOpts))
		RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.ConfirmRunnerKey{}, tc.promptMock))
		RemoveCmd.SetArgs(tc.args)

		err := RemoveCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}

	bucketPolicyOpts.SetZeroValues()
}
