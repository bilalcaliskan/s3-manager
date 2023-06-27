//go:build e2e

package remove

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/stretchr/testify/assert"
)

var (
	defaultGetBucketPolicyOutput = &s3.GetBucketPolicyOutput{}
	defaultGetBucketPolicyErr    error

	defaultDeleteBucketPolicyOutput = &s3.DeleteBucketPolicyOutput{}
	defaultDeleteBucketPolicyErr    error

	policyStr = `
{
  "Statement": [
    {
      "Action": "s3:*",
      "Condition": {
        "Bool": {
          "aws:SecureTransport": "false"
        }
      },
      "Effect": "Deny",
      "Principal": "*",
      "Resource": [
        "arn:aws:s3:::thevpnbeast-releases-1",
        "arn:aws:s3:::thevpnbeast-releases-1/*"
      ],
      "Sid": "RestrictToTLSRequestsOnly"
    }
  ],
  "Version": "2012-10-17"
}
`
)

func createSvc(rootOpts *options.RootOptions) (*s3.S3, error) {
	return internalaws.CreateAwsService(rootOpts)
}

type promptMock struct {
	msg string
	err error
}

func (p promptMock) Run() (string, error) {
	// return expected result
	return p.msg, p.err
}

// Define a testdata struct to be used in your unit tests
type mockS3Client struct {
	s3iface.S3API
}

func (m *mockS3Client) GetBucketPolicy(input *s3.GetBucketPolicyInput) (*s3.GetBucketPolicyOutput, error) {
	return defaultGetBucketPolicyOutput, defaultGetBucketPolicyErr
}

func (m *mockS3Client) DeleteBucketPolicy(input *s3.DeleteBucketPolicyInput) (*s3.DeleteBucketPolicyOutput, error) {
	return defaultDeleteBucketPolicyOutput, defaultDeleteBucketPolicyErr
}

func TestExecuteRemoveCmd(t *testing.T) {
	ctx := context.Background()
	RemoveCmd.SetContext(ctx)

	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	cases := []struct {
		caseName                 string
		args                     []string
		shouldPass               bool
		shouldMock               bool
		getBucketPolicyOutput    *s3.GetBucketPolicyOutput
		getBucketPolicyErr       error
		deleteBucketPolicyErr    error
		deleteBucketPolicyOutput *s3.DeleteBucketPolicyOutput
		promptMock               *promptMock
		dryRun                   bool
		autoApprove              bool
	}{
		{"Too many arguments", []string{"enabled", "foo"}, false, false,
			&s3.GetBucketPolicyOutput{}, nil, nil,
			&s3.DeleteBucketPolicyOutput{}, nil, false, false,
		},
		{"Success", []string{}, true, true,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String(policyStr),
			}, nil,
			nil, &s3.DeleteBucketPolicyOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			}, false, false,
		},
		{"Success with dry run", []string{}, true, true,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String(policyStr),
			}, nil,
			nil, &s3.DeleteBucketPolicyOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			}, true, false,
		},
		{"Success with auto approve", []string{}, true, true,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String(policyStr),
			}, nil,
			nil, &s3.DeleteBucketPolicyOutput{},
			nil, false, true,
		},
		{"Failure", []string{}, false, true,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String(policyStr),
			}, nil,
			errors.New("injected error"),
			&s3.DeleteBucketPolicyOutput{}, nil, false, false,
		},
		{"Failure caused by get bucket policy error", []string{}, false, true,
			&s3.GetBucketPolicyOutput{
				Policy: nil,
			}, errors.New("injected error"),
			nil,
			&s3.DeleteBucketPolicyOutput{}, nil, false, false,
		},
		{"Failure caused by user terminated process", []string{}, false, true,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String(policyStr),
			}, nil,
			nil,
			&s3.DeleteBucketPolicyOutput{}, &promptMock{
				msg: "n",
				err: nil,
			}, false, false,
		},
		{"Failure caused by prompt error", []string{}, false, true,
			&s3.GetBucketPolicyOutput{
				Policy: aws.String(policyStr),
			}, nil,
			nil,
			&s3.DeleteBucketPolicyOutput{}, &promptMock{
				msg: "nasdfadf",
				err: errors.New("injected error"),
			}, false, false,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		defaultGetBucketPolicyOutput = tc.getBucketPolicyOutput
		defaultGetBucketPolicyErr = tc.getBucketPolicyErr

		defaultDeleteBucketPolicyErr = tc.deleteBucketPolicyErr
		defaultDeleteBucketPolicyOutput = tc.deleteBucketPolicyOutput

		rootOpts.DryRun = tc.dryRun
		rootOpts.AutoApprove = tc.autoApprove

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

		if tc.promptMock != nil {
			confirmRunner = tc.promptMock
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

	bucketPolicyOpts.SetZeroValues()
}
