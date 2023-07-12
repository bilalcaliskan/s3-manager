//go:build e2e

package show

import (
	"context"
	"errors"
	"testing"

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
	defaultGetBucketAccelerationErr error
)

type mockS3Client struct {
	s3iface.S3API
}

func (m *mockS3Client) GetBucketAccelerateConfiguration(input *s3.GetBucketAccelerateConfigurationInput) (*s3.GetBucketAccelerateConfigurationOutput, error) {
	return defaultGetBucketAccelerationOutput, defaultGetBucketAccelerationErr
}

func TestExecuteShowCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	svc, err := internalaws.CreateAwsService(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	ctx := context.Background()
	ShowCmd.SetContext(ctx)

	cases := []struct {
		caseName                    string
		args                        []string
		shouldPass                  bool
		svc                         s3iface.S3API
		getBucketAccelerationErr    error
		getBucketAccelerationOutput *s3.GetBucketAccelerateConfigurationOutput
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
		},
		{
			"Success enabled",
			[]string{},
			true,
			&mockS3Client{},
			nil,
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Enabled"),
			},
		},
		{
			"Success suspended",
			[]string{},
			true,
			&mockS3Client{},
			nil,
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Suspended"),
			},
		},
		{
			"Failure get bucket acceleration",
			[]string{},
			false,
			&mockS3Client{},
			errors.New("dummy error"),
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Enabled"),
			},
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
		},
	}

	for _, tc := range cases {
		defaultGetBucketAccelerationErr = tc.getBucketAccelerationErr
		defaultGetBucketAccelerationOutput = tc.getBucketAccelerationOutput

		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3SvcKey{}, tc.svc))
		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))
		ShowCmd.SetArgs(tc.args)

		err = ShowCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}

	transferAccelerationOpts.SetZeroValues()
}
