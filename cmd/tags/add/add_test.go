//go:build e2e

package add

import (
	"context"
	"testing"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/stretchr/testify/assert"
)

var (
	defaultPutBucketTaggingErr    error
	defaultPutBucketTaggingOutput = &s3.PutBucketTaggingOutput{}

	defaultGetBucketTaggingErr    error
	defaultGetBucketTaggingOutput = &s3.GetBucketTaggingOutput{}
)

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

func (m *mockS3Client) PutBucketTagging(input *s3.PutBucketTaggingInput) (*s3.PutBucketTaggingOutput, error) {
	return defaultPutBucketTaggingOutput, defaultPutBucketTaggingErr
}

func (m *mockS3Client) GetBucketTagging(input *s3.GetBucketTaggingInput) (*s3.GetBucketTaggingOutput, error) {
	return defaultGetBucketTaggingOutput, defaultGetBucketTaggingErr
}

func TestExecuteAddCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	svc, err := internalaws.CreateAwsService(rootOpts)
	assert.Nil(t, err)
	assert.NotNil(t, svc)

	ctx := context.Background()
	AddCmd.SetContext(ctx)

	cases := []struct {
		caseName               string
		args                   []string
		shouldPass             bool
		svc                    s3iface.S3API
		getBucketTaggingErr    error
		getBucketTaggingOutput *s3.GetBucketTaggingOutput
		putBucketTaggingErr    error
		putBucketTaggingOutput *s3.PutBucketTaggingOutput
		promptMock             *promptMock
		dryRun                 bool
		autoApprove            bool
	}{
		{
			"No arguments provided",
			[]string{},
			false,
			svc,
			nil,
			&s3.GetBucketTaggingOutput{},
			nil,
			&s3.PutBucketTaggingOutput{},
			nil,
			false,
			false,
		},
		{
			"Success when auto-approve disabled",
			[]string{"foo=bar,foo2=bar2"},
			true,
			&mockS3Client{},
			nil,
			&s3.GetBucketTaggingOutput{
				TagSet: []*s3.Tag{
					{
						Key:   aws.String("hasan1"),
						Value: aws.String("huseyin1"),
					},
					{
						Key:   aws.String("hasan2"),
						Value: aws.String("huseyin2"),
					},
				},
			},
			nil,
			&s3.PutBucketTaggingOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
		{
			"Success when auto-approve enabled",
			[]string{"foo=bar,foo2=bar2"},
			true,
			&mockS3Client{},
			nil,
			&s3.GetBucketTaggingOutput{
				TagSet: []*s3.Tag{
					{
						Key:   aws.String("hasan1"),
						Value: aws.String("huseyin1"),
					},
					{
						Key:   aws.String("hasan2"),
						Value: aws.String("huseyin2"),
					},
				},
			},
			nil,
			&s3.PutBucketTaggingOutput{},
			nil,
			false,
			true,
		},
		{
			"Success when dry run enabled",
			[]string{"foo=bar,foo2=bar2"},
			true,
			&mockS3Client{},
			nil,
			&s3.GetBucketTaggingOutput{
				TagSet: []*s3.Tag{
					{
						Key:   aws.String("hasan1"),
						Value: aws.String("huseyin1"),
					},
					{
						Key:   aws.String("hasan2"),
						Value: aws.String("huseyin2"),
					},
				},
			},
			nil,
			&s3.PutBucketTaggingOutput{},
			nil,
			true,
			false,
		},
		{
			"Success",
			[]string{"foo=bar,foo2=bar2"},
			true,
			&mockS3Client{},
			nil,
			&s3.GetBucketTaggingOutput{
				TagSet: []*s3.Tag{
					{
						Key:   aws.String("hasan1"),
						Value: aws.String("huseyin1"),
					},
					{
						Key:   aws.String("hasan2"),
						Value: aws.String("huseyin2"),
					},
				},
			},
			nil,
			&s3.PutBucketTaggingOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by GetBucketTags error",
			[]string{"foo=bar,foo2=bar2"},
			false,
			&mockS3Client{},
			constants.ErrInjected,
			&s3.GetBucketTaggingOutput{},
			nil,
			&s3.PutBucketTaggingOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by prompt error",
			[]string{"foo=bar,foo2=bar2"},
			false,
			&mockS3Client{},
			nil,
			&s3.GetBucketTaggingOutput{},
			nil,
			&s3.PutBucketTaggingOutput{},
			&promptMock{
				msg: "asdfasdf",
				err: constants.ErrInjected,
			},
			false,
			false,
		},
		{
			"Failure caused by user terminated the process",
			[]string{"foo=bar,foo2=bar2"},
			false,
			&mockS3Client{},
			nil,
			&s3.GetBucketTaggingOutput{},
			nil,
			&s3.PutBucketTaggingOutput{},
			&promptMock{
				msg: "n",
				err: constants.ErrInjected,
			},
			false,
			false,
		},
		{
			"Failure caused by wrong provided arg",
			[]string{"foo=bar=barX,foo2=bar2"},
			false,
			&mockS3Client{},
			nil,
			&s3.GetBucketTaggingOutput{},
			nil,
			&s3.PutBucketTaggingOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by SetBucketTags error",
			[]string{"foo=bar,foo2=bar2"},
			false,
			&mockS3Client{},
			nil,
			&s3.GetBucketTaggingOutput{},
			constants.ErrInjected,
			&s3.PutBucketTaggingOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
	}

	for _, tc := range cases {
		rootOpts.DryRun = tc.dryRun
		rootOpts.AutoApprove = tc.autoApprove

		defaultGetBucketTaggingErr = tc.getBucketTaggingErr
		defaultPutBucketTaggingErr = tc.putBucketTaggingErr
		defaultGetBucketTaggingOutput = tc.getBucketTaggingOutput
		defaultPutBucketTaggingOutput = tc.putBucketTaggingOutput

		if tc.promptMock != nil {
			confirmRunner = tc.promptMock
		}

		AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.S3SvcKey{}, tc.svc))
		AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.OptsKey{}, rootOpts))
		AddCmd.SetArgs(tc.args)

		err = AddCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}

	tagOpts.SetZeroValues()
}
