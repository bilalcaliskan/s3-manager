//go:build e2e

package enabled

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
	defaultGetBucketVersioningErr    error
	defaultPutBucketVersioningOutput = &s3.PutBucketVersioningOutput{}
	defaultPutBucketVersioningErr    error
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

func (m *mockS3Client) GetBucketVersioning(input *s3.GetBucketVersioningInput) (*s3.GetBucketVersioningOutput, error) {
	return defaultGetBucketVersioningOutput, defaultGetBucketVersioningErr
}

func (m *mockS3Client) PutBucketVersioning(input *s3.PutBucketVersioningInput) (*s3.PutBucketVersioningOutput, error) {
	return defaultPutBucketVersioningOutput, defaultPutBucketVersioningErr
}

func TestExecuteEnabledCmd(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	EnabledCmd.SetContext(ctx)

	cases := []struct {
		caseName                  string
		args                      []string
		shouldPass                bool
		shouldMock                bool
		getBucketVersioningErr    error
		getBucketVersioningOutput *s3.GetBucketVersioningOutput
		putBucketVersioningErr    error
		putBucketVersioningOutput *s3.PutBucketVersioningOutput
		promptMock                *promptMock
		dryRun                    bool
		autoApprove               bool
	}{
		{"Too many arguments", []string{"enabled", "foo"}, false, false, nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			}, nil, &s3.PutBucketVersioningOutput{},
			nil, false, false,
		},
		{"Success", []string{}, true, true, nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Suspended"),
			}, nil, &s3.PutBucketVersioningOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			}, false, false,
		},
		{"Success when dry-run enabled", []string{}, true, true, nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Suspended"),
			}, nil, &s3.PutBucketVersioningOutput{},
			nil, true, false,
		},
		{"Success when auto-approve enabled", []string{}, true, true, nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Suspended"),
			}, nil, &s3.PutBucketVersioningOutput{},
			nil, false, true,
		},
		{"Success while already enabled", []string{}, true, true,
			nil, &s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			}, nil, &s3.PutBucketVersioningOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			}, false, false,
		},
		{"Failure caused by unknown status returned by external call", []string{}, false, true, nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enableddd"),
			}, nil, &s3.PutBucketVersioningOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			}, false, false,
		},
		{"Failure caused by prompt error", []string{}, false, true, nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			}, nil, &s3.PutBucketVersioningOutput{},
			&promptMock{
				msg: "asdfafj",
				err: errors.New("injected error"),
			}, false, false,
		},
		{"Failure caused by user terminated the process", []string{}, false, true, nil,
			&s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			}, nil, &s3.PutBucketVersioningOutput{},
			&promptMock{
				msg: "n",
				err: nil,
			}, false, false,
		},
	}

	for _, tc := range cases {
		rootOpts.DryRun = tc.dryRun
		rootOpts.AutoApprove = tc.autoApprove

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

		if tc.promptMock != nil {
			confirmRunner = tc.promptMock
		}

		EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.S3SvcKey{}, svc))
		EnabledCmd.SetContext(context.WithValue(EnabledCmd.Context(), options.OptsKey{}, rootOpts))
		EnabledCmd.SetArgs(tc.args)

		err = EnabledCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}

	rootOpts.SetZeroValues()
	versioningOpts.SetZeroValues()
}
