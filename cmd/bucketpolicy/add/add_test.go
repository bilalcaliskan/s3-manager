//go:build e2e

package add

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
	defaultPutBucketPolicyOutput = &s3.PutBucketPolicyOutput{}
	defaultPutBucketPolicyErr    error
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

func (m *mockS3Client) PutBucketPolicy(input *s3.PutBucketPolicyInput) (*s3.PutBucketPolicyOutput, error) {
	return defaultPutBucketPolicyOutput, defaultPutBucketPolicyErr
}

func TestExecuteAddCmd(t *testing.T) {
	ctx := context.Background()
	AddCmd.SetContext(ctx)

	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	cases := []struct {
		caseName              string
		args                  []string
		shouldPass            bool
		shouldMock            bool
		putBucketPolicyErr    error
		putBucketPolicyOutput *s3.PutBucketPolicyOutput
		promptMock            *promptMock
		dryRun                bool
		autoApprove           bool
	}{
		{"Success", []string{"../../../testdata/bucketpolicy.json"},
			true, true,
			nil, &s3.PutBucketPolicyOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			}, false, false,
		},
		{"Success with dry-run",
			[]string{"../../../testdata/bucketpolicy.json"},
			true, true,
			nil, &s3.PutBucketPolicyOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			}, true, false,
		},
		{"Failure", []string{"../../../testdata/bucketpolicy.json"},
			false, true,
			errors.New("dummy error"),
			&s3.PutBucketPolicyOutput{}, nil, false, false,
		},
		{"Failure caused by user terminated process", []string{"../../../testdata/bucketpolicy.json"},
			false, true,
			nil, &s3.PutBucketPolicyOutput{},
			&promptMock{
				msg: "n",
				err: nil,
			}, false, false,
		},
		{"Failure caused by prompt error", []string{"../../../testdata/bucketpolicy.json"},
			false, true,
			nil, &s3.PutBucketPolicyOutput{},
			&promptMock{
				msg: "nasdasd",
				err: errors.New("injected error"),
			}, false, false,
		},
		{"Failure caused by target file not found", []string{"../../../testdata/bucketpolicy.jsonnnn"},
			false, true, nil,
			&s3.PutBucketPolicyOutput{}, nil, false, false,
		},
		{"Failure caused by too many arguments error", []string{"enabled", "foo"},
			false, false, nil, &s3.PutBucketPolicyOutput{},
			nil, false, false,
		},
		{"Failure caused by no arguments provided error", []string{}, false, false,
			nil, &s3.PutBucketPolicyOutput{}, nil, false, false,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		rootOpts.DryRun = tc.dryRun
		rootOpts.AutoApprove = tc.autoApprove

		defaultPutBucketPolicyErr = tc.putBucketPolicyErr
		defaultPutBucketPolicyOutput = tc.putBucketPolicyOutput

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

		AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.S3SvcKey{}, svc))
		AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.OptsKey{}, rootOpts))
		AddCmd.SetArgs(tc.args)

		err = AddCmd.Execute()
		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}

	bucketPolicyOpts.SetZeroValues()
}
