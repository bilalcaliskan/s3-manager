//go:build e2e

package disabled

import (
	"context"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/pkg/aws/types"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/constants"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/prompt"
	"testing"

	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

func TestExecuteDisabledCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()

	ctx := context.Background()
	DisabledCmd.SetContext(ctx)

	cases := []struct {
		caseName                string
		args                    []string
		shouldPass              bool
		getBucketVersioningFunc func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error)
		putBucketVersioningFunc func(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error)
		prompt.PromptRunner
		dryRun      bool
		autoApprove bool
	}{
		{
			"Too many arguments",
			[]string{"enabled", "foo"},
			false,
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return &s3.GetBucketVersioningOutput{
					Status: types.BucketVersioningStatusEnabled,
				}, nil
			},
			internalawstypes.DefaultPutBucketVersioningFunc,
			nil,
			false,
			false,
		},
		{
			"Success when enabled",
			[]string{},
			true,
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return &s3.GetBucketVersioningOutput{
					Status: types.BucketVersioningStatusEnabled,
				}, nil
			},
			internalawstypes.DefaultPutBucketVersioningFunc,
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Success already disabled",
			[]string{},
			true,
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return &s3.GetBucketVersioningOutput{
					Status: types.BucketVersioningStatusSuspended,
				}, nil
			},
			internalawstypes.DefaultPutBucketVersioningFunc,
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Success when dry-run enabled",
			[]string{},
			true,
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return &s3.GetBucketVersioningOutput{
					Status: types.BucketVersioningStatusEnabled,
				}, nil
			},
			internalawstypes.DefaultPutBucketVersioningFunc,
			nil,
			true,
			false,
		},
		{
			"Success when auto-approve enabled",
			[]string{},
			true,
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return &s3.GetBucketVersioningOutput{
					Status: types.BucketVersioningStatusEnabled,
				}, nil
			},
			internalawstypes.DefaultPutBucketVersioningFunc,
			nil,
			false,
			true,
		},
		{
			"Failure unknown status",
			[]string{},
			false,
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return &s3.GetBucketVersioningOutput{
					Status: "Enableddd",
				}, nil
			},
			internalawstypes.DefaultPutBucketVersioningFunc,
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by user terminated the process",
			[]string{},
			false,
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return &s3.GetBucketVersioningOutput{
					Status: types.BucketVersioningStatusEnabled,
				}, nil
			},
			internalawstypes.DefaultPutBucketVersioningFunc,
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
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return &s3.GetBucketVersioningOutput{
					Status: types.BucketVersioningStatusEnabled,
				}, nil
			},
			internalawstypes.DefaultPutBucketVersioningFunc,
			prompt.PromptMock{
				Msg: "asdfasfd",
				Err: constants.ErrInjected,
			},
			false,
			false,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		rootOpts.DryRun = tc.dryRun
		rootOpts.AutoApprove = tc.autoApprove

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.GetBucketVersioningAPI = tc.getBucketVersioningFunc
		mockS3.PutBucketVersioningAPI = tc.putBucketVersioningFunc

		DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.S3ClientKey{}, mockS3))
		DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.ConfirmRunnerKey{}, tc.PromptRunner))
		DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.OptsKey{}, rootOpts))
		DisabledCmd.SetArgs(tc.args)

		err := DisabledCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		versioningOpts.SetZeroValues()
	}
}
