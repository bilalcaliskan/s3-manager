//go:build e2e

package show

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
	defaultGetBucketPolicyOutput = &s3.GetBucketPolicyOutput{
		Policy: aws.String("{}"),
	}
	defaultGetBucketPolicyErr error
)

func createSvc(rootOpts *options.RootOptions) (*s3.S3, error) {
	return internalaws.CreateAwsService(rootOpts)
}

// Define a testdata struct to be used in your unit tests
type mockS3Client struct {
	s3iface.S3API
}

func (m *mockS3Client) GetBucketPolicy(input *s3.GetBucketPolicyInput) (*s3.GetBucketPolicyOutput, error) {
	return defaultGetBucketPolicyOutput, defaultGetBucketPolicyErr
}

func TestExecuteShowCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	ctx := context.Background()
	ShowCmd.SetContext(ctx)

	cases := []struct {
		caseName              string
		args                  []string
		shouldPass            bool
		shouldMock            bool
		getBucketPolicyErr    error
		getBucketPolicyOutput *s3.GetBucketPolicyOutput
	}{
		{"Too many arguments", []string{"enabled", "foo"}, false, false,
			nil, &s3.GetBucketPolicyOutput{
				Policy: aws.String("{}"),
			},
		},
		{"No argument", []string{}, false, false, nil,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String("{}"),
			},
		},
		{"Success", []string{}, true, true, nil,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String("{}"),
			},
		},
		{"Json failure", []string{}, false, true, nil,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String(""),
			},
		},
	}

	for _, tc := range cases {
		defaultGetBucketPolicyErr = tc.getBucketPolicyErr
		defaultGetBucketPolicyOutput = tc.getBucketPolicyOutput

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
	bucketPolicyOpts.SetZeroValues()
}
