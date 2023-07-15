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

// Define a testdata struct to be used in your unit tests
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

func TestExecuteShowCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()

	ctx := context.Background()
	ShowCmd.SetContext(ctx)

	cases := []struct {
		caseName                  string
		args                      []string
		shouldPass                bool
		getBucketVersioningErr    error
		getBucketVersioningOutput *s3.GetBucketVersioningOutput
	}{
		{
			"Too many arguments",
			[]string{"enabled", "foo"},
			false,
			nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			},
		},
		{
			"Success while already enabled",
			[]string{},
			true,
			nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			},
		},
		{
			"Success while disabled",
			[]string{},
			true,
			nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Suspended"),
			},
		},
		{
			"Failure caused by GetBucketVersioning error",
			[]string{},
			false,
			errors.New("dummy error"), &s3.GetBucketVersioningOutput{
				Status: aws.String("Suspended"),
			},
		},
		{
			"Failure caused by unknown status returned by external call",
			[]string{},
			false,
			nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enableddd"),
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(mockS3Client)
		mockS3.On("GetBucketVersioning", mock.AnythingOfType("*s3.GetBucketVersioningInput")).Return(tc.getBucketVersioningOutput, tc.getBucketVersioningErr)

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

	versioningOpts.SetZeroValues()
}
