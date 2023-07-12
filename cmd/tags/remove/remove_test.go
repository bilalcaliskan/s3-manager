//go:build e2e

package remove

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

	defaultDeleteBucketTaggingErr    error
	defaultDeleteBucketTaggingOutput = &s3.DeleteBucketTaggingOutput{}
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

func (m *mockS3Client) DeleteBucketTagging(input *s3.DeleteBucketTaggingInput) (*s3.DeleteBucketTaggingOutput, error) {
	return defaultDeleteBucketTaggingOutput, defaultDeleteBucketTaggingErr
}

func TestExecuteRemoveCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()
	svc, err := internalaws.CreateAwsService(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	ctx := context.Background()
	RemoveCmd.SetContext(ctx)

	cases := []struct {
		caseName                  string
		args                      []string
		shouldPass                bool
		svc                       s3iface.S3API
		getBucketTaggingErr       error
		getBucketTaggingOutput    *s3.GetBucketTaggingOutput
		putBucketTaggingErr       error
		putBucketTaggingOutput    *s3.PutBucketTaggingOutput
		deleteBucketTaggingErr    error
		deleteBucketTaggingOutput *s3.DeleteBucketTaggingOutput
		promptMock                *promptMock
		dryRun                    bool
		autoApprove               bool
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
			&s3.DeleteBucketTaggingOutput{},
			nil,
			false,
			false,
		},
		{
			"Success while has single tag",
			[]string{"foo=bar,foo2=bar2"},
			true,
			&mockS3Client{},
			nil,
			&s3.GetBucketTaggingOutput{
				TagSet: []*s3.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
				},
			},
			nil,
			&s3.PutBucketTaggingOutput{},
			nil,
			&s3.DeleteBucketTaggingOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
		{
			"Success when dry-run enabled",
			[]string{"foo=bar,foo2=bar2"},
			true, &mockS3Client{},
			nil,
			&s3.GetBucketTaggingOutput{
				TagSet: []*s3.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
				},
			},
			nil,
			&s3.PutBucketTaggingOutput{},
			nil,
			&s3.DeleteBucketTaggingOutput{},
			nil,
			true,
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
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
				},
			},
			nil,
			&s3.PutBucketTaggingOutput{},
			nil,
			&s3.DeleteBucketTaggingOutput{},
			nil,
			false,
			true,
		},
		{
			"Success while has multiple tags",
			[]string{"foo=bar,foo2=bar2"},
			true, &mockS3Client{},
			nil,
			&s3.GetBucketTaggingOutput{
				TagSet: []*s3.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
					{
						Key:   aws.String("foo2"),
						Value: aws.String("bar2"),
					},
					{
						Key:   aws.String("foo3"),
						Value: aws.String("bar3"),
					},
					{
						Key:   aws.String("foo4"),
						Value: aws.String("bar4"),
					},
				},
			},
			nil,
			&s3.PutBucketTaggingOutput{},
			nil,
			&s3.DeleteBucketTaggingOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by wrong argument provided",
			[]string{"foo=bar=bar3,foo2=bar2"},
			false,
			&mockS3Client{},
			nil,
			&s3.GetBucketTaggingOutput{},
			nil,
			&s3.PutBucketTaggingOutput{},
			nil,
			&s3.DeleteBucketTaggingOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
		{
			"Warn while has no tags to remove",
			[]string{"foo3=bar3,foo4=bar4"},
			true,
			&mockS3Client{},
			nil,
			&s3.GetBucketTaggingOutput{
				TagSet: []*s3.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
					{
						Key:   aws.String("foo2"),
						Value: aws.String("bar2"),
					},
				},
			},
			nil,
			&s3.PutBucketTaggingOutput{},
			nil,
			&s3.DeleteBucketTaggingOutput{},
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
			constants.ErrInjected,
			&s3.DeleteBucketTaggingOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by DeleteAllBucketTags error",
			[]string{"foo=bar,foo2=bar2"},
			false,
			&mockS3Client{},
			nil,
			&s3.GetBucketTaggingOutput{
				TagSet: []*s3.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
					{
						Key:   aws.String("foo2"),
						Value: aws.String("bar2"),
					},
				},
			},
			nil,
			&s3.PutBucketTaggingOutput{},
			constants.ErrInjected,
			&s3.DeleteBucketTaggingOutput{},
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
			&s3.GetBucketTaggingOutput{
				TagSet: []*s3.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
					{
						Key:   aws.String("foo2"),
						Value: aws.String("bar2"),
					},
				},
			},
			nil,
			&s3.PutBucketTaggingOutput{},
			nil,
			&s3.DeleteBucketTaggingOutput{},
			&promptMock{
				msg: "yasdfas",
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
			&s3.GetBucketTaggingOutput{
				TagSet: []*s3.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
					{
						Key:   aws.String("foo2"),
						Value: aws.String("bar2"),
					},
				},
			},
			nil,
			&s3.PutBucketTaggingOutput{},
			nil,
			&s3.DeleteBucketTaggingOutput{},
			&promptMock{
				msg: "n",
				err: constants.ErrInjected,
			},
			false,
			false,
		},
		{
			"Failure caused by DeleteAllBucketTags error",
			[]string{"foo=bar,foo2=bar2"},
			false,
			&mockS3Client{},
			nil,
			&s3.GetBucketTaggingOutput{
				TagSet: []*s3.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
					{
						Key:   aws.String("foo2"),
						Value: aws.String("bar2"),
					},
				},
			},
			constants.ErrInjected,
			&s3.PutBucketTaggingOutput{},
			nil,
			&s3.DeleteBucketTaggingOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		rootOpts.DryRun = tc.dryRun
		rootOpts.AutoApprove = tc.autoApprove

		defaultGetBucketTaggingErr = tc.getBucketTaggingErr
		defaultPutBucketTaggingErr = tc.putBucketTaggingErr
		defaultGetBucketTaggingOutput = tc.getBucketTaggingOutput
		defaultPutBucketTaggingOutput = tc.putBucketTaggingOutput
		defaultDeleteBucketTaggingErr = tc.deleteBucketTaggingErr
		defaultDeleteBucketTaggingOutput = tc.deleteBucketTaggingOutput

		if tc.promptMock != nil {
			confirmRunner = tc.promptMock
		}

		RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.S3SvcKey{}, tc.svc))
		RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.OptsKey{}, rootOpts))
		RemoveCmd.SetArgs(tc.args)

		err = RemoveCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}

	tagOpts.SetZeroValues()
}
