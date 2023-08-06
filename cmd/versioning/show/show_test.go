//go:build e2e

package show

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/aws/types"
	"github.com/bilalcaliskan/s3-manager/internal/constants"
	"testing"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

func TestExecuteShowCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()

	ctx := context.Background()
	ShowCmd.SetContext(ctx)

	cases := []struct {
		caseName                string
		args                    []string
		shouldPass              bool
		getBucketVersioningFunc func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error)
	}{
		{
			"Too many arguments",
			[]string{"enabled", "foo"},
			false,
			internalawstypes.DefaultGetBucketVersioningFunc,
		},
		{
			"Success while already enabled",
			[]string{},
			true,
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return &s3.GetBucketVersioningOutput{
					Status: types.BucketVersioningStatusEnabled,
				}, nil
			},
		},
		{
			"Success while disabled",
			[]string{},
			true,
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return &s3.GetBucketVersioningOutput{
					Status: types.BucketVersioningStatusSuspended,
				}, nil
			},
		},
		{
			"Failure caused by GetBucketVersioning error",
			[]string{},
			false,
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return nil, constants.ErrInjected
			},
		},
		{
			"Failure caused by unknown status returned by external call",
			[]string{},
			false,
			func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
				return &s3.GetBucketVersioningOutput{
					Status: "Enableddd",
				}, nil
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(internalawstypes.MockS3v2Client)
		mockS3.GetBucketVersioningAPI = tc.getBucketVersioningFunc

		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3ClientKey{}, mockS3))
		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))
		ShowCmd.SetArgs(tc.args)

		err := ShowCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		versioningOpts.SetZeroValues()
	}
}
