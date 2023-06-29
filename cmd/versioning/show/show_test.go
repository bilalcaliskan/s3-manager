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
	defaultGetBucketVersioningOutput = &s3.GetBucketVersioningOutput{
		Status: aws.String("Enabled"),
	}
	defaultGetBucketVersioningErr error
)

func createSvc(rootOpts *options.RootOptions) (*s3.S3, error) {
	return internalaws.CreateAwsService(rootOpts)
}

// Define a testdata struct to be used in your unit tests
type mockS3Client struct {
	s3iface.S3API
}

func (m *mockS3Client) GetBucketVersioning(input *s3.GetBucketVersioningInput) (*s3.GetBucketVersioningOutput, error) {
	return defaultGetBucketVersioningOutput, defaultGetBucketVersioningErr
}

func TestExecuteShowCmd(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	ShowCmd.SetContext(ctx)

	cases := []struct {
		caseName                  string
		args                      []string
		shouldPass                bool
		shouldMock                bool
		getBucketVersioningErr    error
		getBucketVersioningOutput *s3.GetBucketVersioningOutput
	}{
		{"Too many arguments", []string{"enabled", "foo"}, false, false, nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			},
		},
		{"Success while already enabled", []string{}, true, true, nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			},
		},
		{"Success while disabled", []string{}, true, true, nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Suspended"),
			},
		},
		{"Failure caused by GetBucketVersioning error", []string{}, false, true,
			errors.New("dummy error"), &s3.GetBucketVersioningOutput{
				Status: aws.String("Suspended"),
			},
		},
		{"Failure caused by unknown status returned by external call", []string{}, false, true, nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enableddd"),
			},
		},
	}

	for _, tc := range cases {
		defaultGetBucketVersioningErr = tc.getBucketVersioningErr
		defaultGetBucketVersioningOutput = tc.getBucketVersioningOutput

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

		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3SvcKey{}, svc))
		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))
		ShowCmd.SetArgs(tc.args)

		err = ShowCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}

	rootOpts.SetZeroValues()
	versioningOpts.SetZeroValues()
}
