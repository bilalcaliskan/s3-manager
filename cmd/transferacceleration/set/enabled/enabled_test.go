//go:build e2e

package enabled

import (
	"context"
	"testing"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/stretchr/testify/assert"
)

var (
	defaultGetBucketAccelerationOutput = &s3.GetBucketAccelerateConfigurationOutput{
		Status: aws.String("Enabled"),
	}
	defaultGetBucketAccelerationErr    error
	defaultPutBucketAccelerationOutput = &s3.PutBucketAccelerateConfigurationOutput{}
	defaultPutBucketAccelerationErr    error
)

type promptMock struct {
	msg string
	err error
}

func (p promptMock) Run() (string, error) {
	// return expected result
	return p.msg, p.err
}

// Define a testdata struct to be used in your unit tests
type mockS3Client struct {
	s3iface.S3API
}

func (m *mockS3Client) GetBucketAccelerateConfiguration(input *s3.GetBucketAccelerateConfigurationInput) (*s3.GetBucketAccelerateConfigurationOutput, error) {
	return defaultGetBucketAccelerationOutput, defaultGetBucketAccelerationErr
}

func (m *mockS3Client) PutBucketAccelerateConfiguration(input *s3.PutBucketAccelerateConfigurationInput) (*s3.PutBucketAccelerateConfigurationOutput, error) {
	return defaultPutBucketAccelerationOutput, defaultPutBucketAccelerationErr
}

func TestExecuteEnabledCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	svc, err := internalaws.CreateAwsService(rootOpts)
	assert.Nil(t, err)
	assert.NotNil(t, svc)

	ctx := context.Background()
	EnabledCmd.SetContext(ctx)

	cases := []struct {
		caseName                    string
		args                        []string
		shouldPass                  bool
		svc                         s3iface.S3API
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
			svc,
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
			&mockS3Client{},
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
			&mockS3Client{},
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
			&mockS3Client{},
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
			&mockS3Client{},
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
			&mockS3Client{},
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
			&mockS3Client{},
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
			&mockS3Client{},
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
		rootOpts.DryRun = tc.dryRun
		rootOpts.AutoApprove = tc.autoApprove

		defaultGetBucketAccelerationErr = tc.getBucketAccelerationErr
		defaultGetBucketAccelerationOutput = tc.getBucketAccelerationOutput

		if tc.promptMock != nil {
			confirmRunner = tc.promptMock
		}

		EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.S3SvcKey{}, tc.svc))
		EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.OptsKey{}, rootOpts))
		EnabledCmd.SetArgs(tc.args)

		err = EnabledCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}

	transferAccelerationOpts.SetZeroValues()
}
