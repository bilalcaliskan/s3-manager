//go:build e2e

package enabled

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

type promptMock struct {
	msg string
	err error
}

func (p promptMock) Run() (string, error) {
	return p.msg, p.err
}

type mockS3Client struct {
	mock.Mock
	s3iface.S3API
}

// GetBucketVersioning mocks the GetBucketVersioning method of s3iface.S3API
func (m *mockS3Client) GetBucketVersioning(input *s3.GetBucketVersioningInput) (*s3.GetBucketVersioningOutput, error) {
	// Return the mocked output values using the `On` method of testify/mock
	args := m.Called(input)
	return args.Get(0).(*s3.GetBucketVersioningOutput), args.Error(1)
}

// PutBucketVersioning mocks the PutBucketVersioning method of s3iface.S3API
func (m *mockS3Client) PutBucketVersioning(input *s3.PutBucketVersioningInput) (*s3.PutBucketVersioningOutput, error) {
	// Return the mocked output values using the `On` method of testify/mock
	args := m.Called(input)
	return args.Get(0).(*s3.PutBucketVersioningOutput), args.Error(1)
}

func TestExecuteEnabledCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()

	ctx := context.Background()
	EnabledCmd.SetContext(ctx)

	cases := []struct {
		caseName                  string
		args                      []string
		shouldPass                bool
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
		t.Logf("starting case %s", tc.caseName)

		rootOpts.DryRun = tc.dryRun
		rootOpts.AutoApprove = tc.autoApprove

		mockS3 := new(mockS3Client)
		mockS3.On("GetBucketVersioning", mock.AnythingOfType("*s3.GetBucketVersioningInput")).Return(tc.getBucketVersioningOutput, tc.getBucketVersioningErr)
		mockS3.On("PutBucketVersioning", mock.AnythingOfType("*s3.PutBucketVersioningInput")).Return(tc.putBucketVersioningOutput, tc.putBucketVersioningErr)

		EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.S3SvcKey{}, mockS3))
		EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.ConfirmRunnerKey{}, tc.promptMock))
		EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.OptsKey{}, rootOpts))
		EnabledCmd.SetArgs(tc.args)

		err := EnabledCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}

	versioningOpts.SetZeroValues()
}
