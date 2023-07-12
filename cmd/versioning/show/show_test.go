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

// Define a testdata struct to be used in your unit tests
type mockS3Client struct {
	s3iface.S3API
}

func (m *mockS3Client) GetBucketVersioning(input *s3.GetBucketVersioningInput) (*s3.GetBucketVersioningOutput, error) {
	return defaultGetBucketVersioningOutput, defaultGetBucketVersioningErr
}

func TestExecuteShowCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	svc, err := internalaws.CreateAwsService(rootOpts)
	assert.Nil(t, err)
	assert.NotNil(t, svc)

	ctx := context.Background()
	ShowCmd.SetContext(ctx)

	cases := []struct {
		caseName                  string
		args                      []string
		shouldPass                bool
		svc                       s3iface.S3API
		getBucketVersioningErr    error
		getBucketVersioningOutput *s3.GetBucketVersioningOutput
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
		},
		{
			"Success while disabled",
			[]string{},
			true,
			&mockS3Client{},
			nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Suspended"),
			},
		},
		{
			"Failure caused by GetBucketVersioning error",
			[]string{},
			false,
			&mockS3Client{},
			errors.New("dummy error"), &s3.GetBucketVersioningOutput{
				Status: aws.String("Suspended"),
			},
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
		},
	}

	for _, tc := range cases {
		defaultGetBucketVersioningErr = tc.getBucketVersioningErr
		defaultGetBucketVersioningOutput = tc.getBucketVersioningOutput

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

	versioningOpts.SetZeroValues()
}
