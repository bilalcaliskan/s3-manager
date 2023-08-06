//go:build e2e

package show

import (
	"context"
	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws/types"
	"testing"

	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

func TestExecuteShowCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()

	ctx := context.Background()
	ShowCmd.SetContext(ctx)

	cases := []struct {
		caseName                  string
		args                      []string
		shouldPass                bool
		getBucketAccelerationFunc func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error)
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
		},
		{
			"Success enabled",
			[]string{},
			true,
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return &s3.GetBucketAccelerateConfigurationOutput{
					Status: types.BucketAccelerateStatusEnabled,
				}, nil
			},
		},
		{
			"Success suspended",
			[]string{},
			true,
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return &s3.GetBucketAccelerateConfigurationOutput{
					Status: types.BucketAccelerateStatusSuspended,
				}, nil
			},
		},
		{
			"Failure get bucket acceleration",
			[]string{},
			false,
			func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
				return nil, constants.ErrInjected
			},
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
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(internalaws.MockS3v2Client)
		mockS3.GetBucketAccelerateConfigurationAPI = tc.getBucketAccelerationFunc

		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3ClientKey{}, mockS3))
		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))
		ShowCmd.SetArgs(tc.args)

		err := ShowCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		transferAccelerationOpts.SetZeroValues()
	}
}
