//go:build e2e

package add

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/aws/types"
	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

func TestExecuteAddCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()

	ctx := context.Background()
	AddCmd.SetContext(ctx)

	cases := []struct {
		caseName            string
		args                []string
		shouldPass          bool
		putBucketPolicyFunc func(ctx context.Context, params *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error)
		prompt.PromptRunner
		dryRun      bool
		autoApprove bool
	}{
		{
			"Success",
			[]string{"../../../testdata/bucketpolicy.json"},
			true,
			func(ctx context.Context, params *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error) {
				return &s3.PutBucketPolicyOutput{}, nil
			},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Failure",
			[]string{"../../../testdata/bucketpolicy.json"},
			false,
			func(ctx context.Context, params *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error) {
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
			"Failure caused by target file not found",
			[]string{"../../../testdata/bucketpolicy.jsonnnn"},
			false,
			nil,
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by too many arguments error",
			[]string{"enabled", "foo"},
			false,
			nil,
			nil,
			false,
			false,
		},
		{
			"Failure caused by no arguments provided error",
			[]string{},
			false,
			nil,
			nil,
			false,
			false,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		rootOpts.DryRun = tc.dryRun
		rootOpts.AutoApprove = tc.autoApprove

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.PutBucketPolicyAPI = tc.putBucketPolicyFunc

		AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.S3ClientKey{}, mockS3))
		AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.OptsKey{}, rootOpts))
		AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.ConfirmRunnerKey{}, tc.PromptRunner))

		AddCmd.SetArgs(tc.args)

		err := AddCmd.Execute()
		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		bucketPolicyOpts.SetZeroValues()
	}
}
