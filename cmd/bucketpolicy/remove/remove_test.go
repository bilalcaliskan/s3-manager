package remove

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/stretchr/testify/assert"
)

var (
	defaultDeleteBucketPolicyOutput = &s3.DeleteBucketPolicyOutput{}
	defaultDeleteBucketPolicyErr    error
)

func createSvc(rootOpts *options.RootOptions) (*s3.S3, error) {
	return internalaws.CreateAwsService(rootOpts)
}

// Define a testdata struct to be used in your unit tests
type mockS3Client struct {
	s3iface.S3API
}

func (m *mockS3Client) DeleteBucketPolicy(input *s3.DeleteBucketPolicyInput) (*s3.DeleteBucketPolicyOutput, error) {
	return defaultDeleteBucketPolicyOutput, defaultDeleteBucketPolicyErr
}

func TestExecuteRemoveCmd(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	RemoveCmd.SetContext(ctx)

	cases := []struct {
		caseName                 string
		args                     []string
		shouldPass               bool
		shouldMock               bool
		deleteBucketPolicyErr    error
		deleteBucketPolicyOutput *s3.DeleteBucketPolicyOutput
	}{
		{"Too many arguments", []string{"enabled", "foo"}, false, false,
			nil, &s3.DeleteBucketPolicyOutput{},
		},
		{"Success", []string{}, true, true,
			nil, &s3.DeleteBucketPolicyOutput{},
		},
		{"Failure", []string{}, false, true,
			errors.New("dummy error"), &s3.DeleteBucketPolicyOutput{},
		},
	}

	for _, tc := range cases {
		defaultDeleteBucketPolicyErr = tc.deleteBucketPolicyErr
		defaultDeleteBucketPolicyOutput = tc.deleteBucketPolicyOutput

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

		RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.S3SvcKey{}, svc))
		RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.OptsKey{}, rootOpts))
		RemoveCmd.SetArgs(tc.args)

		err = RemoveCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}

	rootOpts.SetZeroValues()
	bucketPolicyOpts.SetZeroValues()
}
