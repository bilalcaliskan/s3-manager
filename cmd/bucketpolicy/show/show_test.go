//go:build e2e

package show

import (
	"context"
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

// GetBucketPolicy mocks the GetBucketPolicy method of s3iface.S3API
func (m *mockS3Client) GetBucketPolicy(input *s3.GetBucketPolicyInput) (*s3.GetBucketPolicyOutput, error) {
	// Return the mocked output values using the `On` method of testify/mock
	args := m.Called(input)
	return args.Get(0).(*s3.GetBucketPolicyOutput), args.Error(1)
}

func TestExecuteShowCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	ctx := context.Background()
	ShowCmd.SetContext(ctx)

	cases := []struct {
		caseName              string
		args                  []string
		shouldPass            bool
		getBucketPolicyErr    error
		getBucketPolicyOutput *s3.GetBucketPolicyOutput
	}{
		{
			"Too many arguments",
			[]string{"enabled", "foo"},
			false,
			nil,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String("{}"),
			},
		},
		{
			"Success",
			[]string{},
			true,
			nil,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String("{}"),
			},
		},
		{
			"Json failure",
			[]string{},
			false,
			nil,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String(""),
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(mockS3Client)
		mockS3.On("GetBucketPolicy", mock.AnythingOfType("*s3.GetBucketPolicyInput")).Return(tc.getBucketPolicyOutput, tc.getBucketPolicyErr)

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

	rootOpts.SetZeroValues()
	bucketPolicyOpts.SetZeroValues()
}
