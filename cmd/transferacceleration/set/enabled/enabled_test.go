//go:build e2e

package enabled

import (
	"context"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/pkg/aws/types"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/constants"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/prompt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

func TestExecuteEnabledCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()

	ctx := context.Background()
	EnabledCmd.SetContext(ctx)

	cases := []struct {
		caseName                  string
		args                      []string
		shouldPass                bool
		getBucketAccelerationFunc func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error)
		putBucketAccelerationFunc func(ctx context.Context, params *s3.PutBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.PutBucketAccelerateConfigurationOutput, error)
		prompt.PromptRunner
		dryRun      bool
		autoApprove bool
	}{
		{
			"Too many arguments",
			[]string{"enabled", "foo"},
			false,
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return &s3.GetBucketAccelerateConfigurationOutput{
					Status: types.BucketAccelerateStatusEnabled,
				}, nil
			},
			internalawstypes.DefaultPutBucketAccelerationFunc,
			nil,
			false,
			false,
		},
		{
			"Success when disabled",
			[]string{},
			true,
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return &s3.GetBucketAccelerateConfigurationOutput{
					Status: types.BucketAccelerateStatusSuspended,
				}, nil
			},
			internalawstypes.DefaultPutBucketAccelerationFunc,
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Success already enabled",
			[]string{},
			true,
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return &s3.GetBucketAccelerateConfigurationOutput{
					Status: types.BucketAccelerateStatusEnabled,
				}, nil
			},
			internalawstypes.DefaultPutBucketAccelerationFunc,
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Success when auto-approve enabled",
			[]string{},
			true,
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return &s3.GetBucketAccelerateConfigurationOutput{
					Status: types.BucketAccelerateStatusSuspended,
				}, nil
			},
			internalawstypes.DefaultPutBucketAccelerationFunc,
			nil,
			false,
			true,
		},
		{
			"Success when dry-run enabled",
			[]string{},
			true,
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return &s3.GetBucketAccelerateConfigurationOutput{
					Status: types.BucketAccelerateStatusSuspended,
				}, nil
			},
			internalawstypes.DefaultPutBucketAccelerationFunc,
			nil,
			true,
			false,
		},
		{
			"Failure unknown status",
			[]string{},
			false,
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return &s3.GetBucketAccelerateConfigurationOutput{
					Status: "Enableddd",
				}, nil
			},
			internalawstypes.DefaultPutBucketAccelerationFunc,
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by prompt error",
			[]string{},
			false,
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return &s3.GetBucketAccelerateConfigurationOutput{
					Status: types.BucketAccelerateStatusEnabled,
				}, nil
			},
			internalawstypes.DefaultPutBucketAccelerationFunc,
			prompt.PromptMock{
				Msg: "asdfadsf",
				Err: constants.ErrInjected,
			},
			false,
			false,
		},
		{
			"Failure caused by user terminated the process",
			[]string{},
			false,
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return &s3.GetBucketAccelerateConfigurationOutput{
					Status: types.BucketAccelerateStatusEnabled,
				}, nil
			},
			internalawstypes.DefaultPutBucketAccelerationFunc,
			prompt.PromptMock{
				Msg: "n",
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
		mockS3.GetBucketAccelerateConfigurationAPI = tc.getBucketAccelerationFunc
		mockS3.PutBucketAccelerateConfigurationAPI = tc.putBucketAccelerationFunc

		EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.S3ClientKey{}, mockS3))
		EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.OptsKey{}, rootOpts))
		EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.ConfirmRunnerKey{}, tc.PromptRunner))
		EnabledCmd.SetArgs(tc.args)

		err := EnabledCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		transferAccelerationOpts.SetZeroValues()
	}
}
