//go:build e2e

package add

import (
	"context"
	"testing"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"

	"github.com/stretchr/testify/mock"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

func TestExecuteAddCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()

	ctx := context.Background()
	AddCmd.SetContext(ctx)

	cases := []struct {
		caseName               string
		args                   []string
		shouldPass             bool
		getBucketTaggingErr    error
		getBucketTaggingOutput *s3.GetBucketTaggingOutput
		putBucketTaggingErr    error
		putBucketTaggingOutput *s3.PutBucketTaggingOutput
		prompt.PromptRunner
		dryRun      bool
		autoApprove bool
	}{
		{
			"No arguments provided",
			[]string{},
			false,
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
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Success when auto-approve enabled",
			[]string{"foo=bar,foo2=bar2"},
			true,
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
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by GetBucketTags error",
			[]string{"foo=bar,foo2=bar2"},
			false,
			constants.ErrInjected,
			&s3.GetBucketTaggingOutput{},
			nil,
			&s3.PutBucketTaggingOutput{},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by prompt error",
			[]string{"foo=bar,foo2=bar2"},
			false,
			nil,
			&s3.GetBucketTaggingOutput{},
			nil,
			&s3.PutBucketTaggingOutput{},
			prompt.PromptMock{
				Msg: "asdfasdf",
				Err: constants.ErrInjected,
			},
			false,
			false,
		},
		{
			"Failure caused by user terminated the process",
			[]string{"foo=bar,foo2=bar2"},
			false,
			nil,
			&s3.GetBucketTaggingOutput{},
			nil,
			&s3.PutBucketTaggingOutput{},
			prompt.PromptMock{
				Msg: "n",
				Err: constants.ErrInjected,
			},
			false,
			false,
		},
		{
			"Failure caused by wrong provided arg",
			[]string{"foo=bar=barX,foo2=bar2"},
			false,
			nil,
			&s3.GetBucketTaggingOutput{},
			nil,
			&s3.PutBucketTaggingOutput{},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by SetBucketTags error",
			[]string{"foo=bar,foo2=bar2"},
			false,
			nil,
			&s3.GetBucketTaggingOutput{},
			constants.ErrInjected,
			&s3.PutBucketTaggingOutput{},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		rootOpts.DryRun = tc.dryRun
		rootOpts.AutoApprove = tc.autoApprove

		mockS3 := new(internalaws.MockS3Client)
		mockS3.On("GetBucketTagging", mock.AnythingOfType("*s3.GetBucketTaggingInput")).Return(tc.getBucketTaggingOutput, tc.getBucketTaggingErr)
		mockS3.On("PutBucketTagging", mock.AnythingOfType("*s3.PutBucketTaggingInput")).Return(tc.putBucketTaggingOutput, tc.putBucketTaggingErr)

		AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.S3SvcKey{}, mockS3))
		AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.ConfirmRunnerKey{}, tc.PromptRunner))
		AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.OptsKey{}, rootOpts))
		AddCmd.SetArgs(tc.args)

		err := AddCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}
