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
	defaultGetBucketVersioningOutput = &s3.GetBucketVersioningOutput{
		Status: aws.String("Enabled"),
	}
	defaultGetBucketVersioningErr    error
	defaultPutBucketVersioningOutput = &s3.PutBucketVersioningOutput{}
	defaultPutBucketVersioningErr    error
)

type promptMock struct {
	msg string
	err error
}

func (p promptMock) Run() (string, error) {
	return p.msg, p.err
}

type mockS3Client struct {
	s3iface.S3API
}

func (m *mockS3Client) GetBucketVersioning(input *s3.GetBucketVersioningInput) (*s3.GetBucketVersioningOutput, error) {
	return defaultGetBucketVersioningOutput, defaultGetBucketVersioningErr
}

func (m *mockS3Client) PutBucketVersioning(input *s3.PutBucketVersioningInput) (*s3.PutBucketVersioningOutput, error) {
	return defaultPutBucketVersioningOutput, defaultPutBucketVersioningErr
}

func TestExecuteEnabledCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	svc, err := internalaws.CreateAwsService(rootOpts)
	assert.Nil(t, err)
	assert.NotNil(t, svc)

	ctx := context.Background()
	EnabledCmd.SetContext(ctx)

	cases := []struct {
		caseName                  string
		args                      []string
		shouldPass                bool
		svc                       s3iface.S3API
		getBucketVersioningErr    error
		getBucketVersioningOutput *s3.GetBucketVersioningOutput
		putBucketVersioningErr    error
		putBucketVersioningOutput *s3.PutBucketVersioningOutput
		promptMock                *promptMock
		dryRun                    bool
		autoApprove               bool
	}{
		{
			"Too many arguments",
			[]string{"enabled", "foo"},
			false,
			svc,
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
			"Success",
			[]string{},
			true,
			&mockS3Client{},
			nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Suspended"),
			},
			nil,
			&s3.PutBucketVersioningOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
		{
			"Success when dry-run enabled",
			[]string{},
			true,
			&mockS3Client{},
			nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Suspended"),
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
			&mockS3Client{},
			nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Suspended"),
			},
			nil,
			&s3.PutBucketVersioningOutput{},
			nil,
			false,
			true,
		},
		{
			"Success while already enabled",
			[]string{},
			true,
			&mockS3Client{},
			nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			},
			nil,
			&s3.PutBucketVersioningOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by unknown status returned by external call",
			[]string{},
			false,
			&mockS3Client{},
			nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enableddd"),
			},
			nil,
			&s3.PutBucketVersioningOutput{},
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
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			},
			nil,
			&s3.PutBucketVersioningOutput{},
			&promptMock{
				msg: "asdfafj",
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
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			},
			nil,
			&s3.PutBucketVersioningOutput{},
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

		defaultGetBucketVersioningErr = tc.getBucketVersioningErr
		defaultGetBucketVersioningOutput = tc.getBucketVersioningOutput

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

	versioningOpts.SetZeroValues()
}
