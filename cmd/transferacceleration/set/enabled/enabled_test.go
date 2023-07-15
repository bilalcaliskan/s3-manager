//go:build e2e

package enabled

import (
	"context"
	"testing"

	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"

	"github.com/stretchr/testify/mock"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

type promptMock struct {
	msg string
	err error
}

func (p promptMock) Run() (string, error) {
	// return expected result
	return p.msg, p.err
}

func TestExecuteEnabledCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()

	ctx := context.Background()
	EnabledCmd.SetContext(ctx)

	cases := []struct {
		caseName                    string
		args                        []string
		shouldPass                  bool
		getBucketAccelerationErr    error
		getBucketAccelerationOutput *s3.GetBucketAccelerateConfigurationOutput
		putBucketAccelerationErr    error
		putBucketAccelerationOutput *s3.PutBucketAccelerateConfigurationOutput
		promptMock                  *promptMock
		dryRun                      bool
		autoApprove                 bool
	}{
		{
			"Too many arguments",
			[]string{"enabled", "foo"},
			false,
			nil,
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Enabled"),
			},
			nil,
			&s3.PutBucketAccelerateConfigurationOutput{},
			nil,
			false,
			false,
		},
		{
			"Success when disabled",
			[]string{},
			true,
			nil,
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Suspended"),
			},
			nil,
			&s3.PutBucketAccelerateConfigurationOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
		{
			"Success already enabled",
			[]string{},
			true,
			nil,
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Enabled"),
			},
			nil,
			&s3.PutBucketAccelerateConfigurationOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
		{
			"Success when auto-approve enabled",
			[]string{},
			true,
			nil,
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Suspended"),
			},
			nil,
			&s3.PutBucketAccelerateConfigurationOutput{},
			nil,
			false,
			true,
		},
		{
			"Success when dry-run enabled",
			[]string{},
			true,
			nil,
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Suspended"),
			},
			nil,
			&s3.PutBucketAccelerateConfigurationOutput{},
			nil,
			true,
			false,
		},
		{
			"Failure unknown status",
			[]string{},
			false,
			nil,
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Enableddd"),
			},
			nil,
			&s3.PutBucketAccelerateConfigurationOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by prompt error",
			[]string{},
			false,
			nil,
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Enabled"),
			},
			nil,
			&s3.PutBucketAccelerateConfigurationOutput{},
			&promptMock{
				msg: "asdfadsf",
				err: constants.ErrInjected,
			},
			false,
			false,
		},
		{
			"Failure caused by user terminated the process",
			[]string{},
			false,
			nil,
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Enabled"),
			},
			nil,
			&s3.PutBucketAccelerateConfigurationOutput{},
			&promptMock{
				msg: "n",
				err: constants.ErrInjected,
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
		mockS3.On("GetBucketAccelerateConfiguration", mock.AnythingOfType("*s3.GetBucketAccelerateConfigurationInput")).Return(tc.getBucketAccelerationOutput, tc.getBucketAccelerationErr)
		mockS3.On("PutBucketAccelerateConfiguration", mock.AnythingOfType("*s3.PutBucketAccelerateConfigurationInput")).Return(tc.putBucketAccelerationOutput, tc.putBucketAccelerationErr)

		EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.S3SvcKey{}, mockS3))
		EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.OptsKey{}, rootOpts))
		EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.ConfirmRunnerKey{}, tc.promptMock))
		EnabledCmd.SetArgs(tc.args)

		err := EnabledCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}

	transferAccelerationOpts.SetZeroValues()
}
