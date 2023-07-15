//go:build e2e

package show

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

type mockS3Client struct {
	mock.Mock
	s3iface.S3API
}

// GetBucketAccelerateConfiguration mocks the GetBucketAccelerateConfiguration method of s3iface.S3API
func (m *mockS3Client) GetBucketAccelerateConfiguration(input *s3.GetBucketAccelerateConfigurationInput) (*s3.GetBucketAccelerateConfigurationOutput, error) {
	// Return the mocked output values using the `On` method of testify/mock
	args := m.Called(input)
	return args.Get(0).(*s3.GetBucketAccelerateConfigurationOutput), args.Error(1)
}

func TestExecuteShowCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()

	ctx := context.Background()
	ShowCmd.SetContext(ctx)

	cases := []struct {
		caseName                    string
		args                        []string
		shouldPass                  bool
		getBucketAccelerationErr    error
		getBucketAccelerationOutput *s3.GetBucketAccelerateConfigurationOutput
	}{
		{
			"Too many arguments",
			[]string{"enabled", "foo"},
			false,
			nil,
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Enabled"),
			},
		},
		{
			"Success enabled",
			[]string{},
			true,
			nil,
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Enabled"),
			},
		},
		{
			"Success suspended",
			[]string{},
			true,
			nil,
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Suspended"),
			},
		},
		{
			"Failure get bucket acceleration",
			[]string{},
			false,
			errors.New("injected error"),
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Enabled"),
			},
		},
		{
			"Failure unknown status",
			[]string{},
			false,
			nil,
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Enableddd"),
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(mockS3Client)
		mockS3.On("GetBucketAccelerateConfiguration", mock.AnythingOfType("*s3.GetBucketAccelerateConfigurationInput")).Return(tc.getBucketAccelerationOutput, tc.getBucketAccelerationErr)

		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3SvcKey{}, mockS3))
		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))
		ShowCmd.SetArgs(tc.args)

		err := ShowCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}

	transferAccelerationOpts.SetZeroValues()
}
