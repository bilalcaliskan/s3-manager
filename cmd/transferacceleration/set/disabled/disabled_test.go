//go:build e2e

package disabled

import (
	"context"
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
	defaultGetBucketAccelerationErr    error
	defaultPutBucketAccelerationOutput = &s3.PutBucketAccelerateConfigurationOutput{}
	defaultPutBucketAccelerationErr    error
)

func createSvc(rootOpts *options.RootOptions) (*s3.S3, error) {
	return internalaws.CreateAwsService(rootOpts)
}

// Define a mock struct to be used in your unit tests
type mockS3Client struct {
	s3iface.S3API
}

func (m *mockS3Client) GetBucketAccelerateConfiguration(input *s3.GetBucketAccelerateConfigurationInput) (*s3.GetBucketAccelerateConfigurationOutput, error) {
	return defaultGetBucketAccelerationOutput, defaultGetBucketAccelerationErr
}

func (m *mockS3Client) PutBucketAccelerateConfiguration(input *s3.PutBucketAccelerateConfigurationInput) (*s3.PutBucketAccelerateConfigurationOutput, error) {
	return defaultPutBucketAccelerationOutput, defaultPutBucketAccelerationErr
}

func TestExecuteDisabledCmd(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	DisabledCmd.SetContext(ctx)

	cases := []struct {
		caseName                    string
		args                        []string
		shouldPass                  bool
		shouldMock                  bool
		getBucketAccelerationErr    error
		getBucketAccelerationOutput *s3.GetBucketAccelerateConfigurationOutput
		putBucketAccelerationErr    error
		putBucketAccelerationOutput *s3.PutBucketAccelerateConfigurationOutput
	}{
		{"Too many arguments", []string{"enabled", "foo"}, false, false, nil,
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Enabled"),
			}, nil, &s3.PutBucketAccelerateConfigurationOutput{},
		},
		{"Success when enabled", []string{}, true, true, nil,
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Enabled"),
			}, nil, &s3.PutBucketAccelerateConfigurationOutput{},
		},
		{"Success already disabled", []string{}, true, true,
			nil, &s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Suspended"),
			}, nil, &s3.PutBucketAccelerateConfigurationOutput{},
		},
		{"Failure unknown status", []string{}, false, true, nil,
			&s3.GetBucketAccelerateConfigurationOutput{
				Status: aws.String("Enableddd"),
			}, nil, &s3.PutBucketAccelerateConfigurationOutput{},
		},
	}

	for _, tc := range cases {
		defaultGetBucketAccelerationErr = tc.getBucketAccelerationErr
		defaultGetBucketAccelerationOutput = tc.getBucketAccelerationOutput

		var err error
		if tc.shouldMock {
			mockSvc := &mockS3Client{}
			svc = mockSvc
			assert.NotNil(t, mockSvc)
		} else {
			svc, err = createSvc(rootOpts)
			assert.NotNil(t, svc)
			assert.Nil(t, err)
		}

		DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.S3SvcKey{}, svc))
		DisabledCmd.SetContext(context.WithValue(DisabledCmd.Context(), options.OptsKey{}, rootOpts))
		DisabledCmd.SetArgs(tc.args)

		err = DisabledCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}

	rootOpts.SetZeroValues()
	transferAccelerationOpts.SetZeroValues()
}
