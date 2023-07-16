//go:build e2e

package disabled

import (
	"context"
	"testing"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"

	"github.com/stretchr/testify/mock"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

func TestExecuteDisabledCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()

	ctx := context.Background()
	DisabledCmd.SetContext(ctx)

	cases := []struct {
		caseName                  string
		args                      []string
		shouldPass                bool
		getBucketVersioningErr    error
		getBucketVersioningOutput *s3.GetBucketVersioningOutput
		putBucketVersioningErr    error
		putBucketVersioningOutput *s3.PutBucketVersioningOutput
		prompt.PromptRunner
		dryRun      bool
		autoApprove bool
	}{
		{
			"Too many arguments",
			[]string{"enabled", "foo"},
			false,
			nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			},
			nil,
			&s3.PutBucketVersioningOutput{},
			nil,
			false,
			false,
		},
		{
			"Success when enabled",
			[]string{},
			true,
			nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			},
			nil,
			&s3.PutBucketVersioningOutput{},
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
			nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Suspended"),
			},
			nil,
			&s3.PutBucketVersioningOutput{},
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
			nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			},
			nil,
			&s3.PutBucketVersioningOutput{},
			nil,
			true,
			false,
		},
		{
			"Success when auto-approve enabled",
			[]string{},
			true,
			nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			},
			nil,
			&s3.PutBucketVersioningOutput{},
			nil,
			false,
			true,
		},
		{
			"Failure unknown status",
			[]string{},
			false,
			nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enableddd"),
			},
			nil,
			&s3.PutBucketVersioningOutput{},
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
			nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			},
			nil,
			&s3.PutBucketVersioningOutput{},
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
			nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			},
			nil,
			&s3.PutBucketVersioningOutput{},
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

		mockS3 := new(internalaws.MockS3Client)
		mockS3.On("GetBucketVersioning", mock.AnythingOfType("*s3.GetBucketVersioningInput")).Return(tc.getBucketVersioningOutput, tc.getBucketVersioningErr)
		mockS3.On("PutBucketVersioning", mock.AnythingOfType("*s3.PutBucketVersioningInput")).Return(tc.putBucketVersioningOutput, tc.putBucketVersioningErr)

		DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.S3SvcKey{}, mockS3))
		DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.ConfirmRunnerKey{}, tc.PromptRunner))
		DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.OptsKey{}, rootOpts))
		DisabledCmd.SetArgs(tc.args)

		err := DisabledCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}
