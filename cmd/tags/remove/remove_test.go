//go:build e2e

package remove

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

func TestExecuteRemoveCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()

	ctx := context.Background()
	RemoveCmd.SetContext(ctx)

	cases := []struct {
		caseName                  string
		args                      []string
		shouldPass                bool
		getBucketTaggingErr       error
		getBucketTaggingOutput    *s3.GetBucketTaggingOutput
		putBucketTaggingErr       error
		putBucketTaggingOutput    *s3.PutBucketTaggingOutput
		deleteBucketTaggingErr    error
		deleteBucketTaggingOutput *s3.DeleteBucketTaggingOutput
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
			&s3.DeleteBucketTaggingOutput{},
			nil,
			false,
			false,
		},
		{
			"Success while has single tag",
			[]string{"foo=bar,foo2=bar2"},
			true,
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
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Success when dry-run enabled",
			[]string{"foo=bar,foo2=bar2"},
			true,
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
			true,
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
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by wrong argument provided",
			[]string{"foo=bar=bar3,foo2=bar2"},
			false,
			nil,
			&s3.GetBucketTaggingOutput{},
			nil,
			&s3.PutBucketTaggingOutput{},
			nil,
			&s3.DeleteBucketTaggingOutput{},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Warn while has no tags to remove",
			[]string{"foo3=bar3,foo4=bar4"},
			true,
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
			constants.ErrInjected,
			&s3.DeleteBucketTaggingOutput{},
			prompt.PromptMock{
				Msg: "y",
				Err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by DeleteAllBucketTags error",
			[]string{"foo=bar,foo2=bar2"},
			false,
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
			prompt.PromptMock{
				Msg: "yasdfas",
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
			prompt.PromptMock{
				Msg: "n",
				Err: constants.ErrInjected,
			},
			false,
			false,
		},
		{
			"Failure caused by DeleteAllBucketTags error",
			[]string{"foo=bar,foo2=bar2"},
			false,
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

		mockS3 := new(internalaws.MockS3Client)
		mockS3.On("GetBucketTagging", mock.AnythingOfType("*s3.GetBucketTaggingInput")).Return(tc.getBucketTaggingOutput, tc.getBucketTaggingErr)
		mockS3.On("PutBucketTagging", mock.AnythingOfType("*s3.PutBucketTaggingInput")).Return(tc.putBucketTaggingOutput, tc.putBucketTaggingErr)
		mockS3.On("DeleteBucketTagging", mock.AnythingOfType("*s3.DeleteBucketTaggingInput")).Return(tc.deleteBucketTaggingOutput, tc.deleteBucketTaggingErr)

		rootOpts.DryRun = tc.dryRun
		rootOpts.AutoApprove = tc.autoApprove

		RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.S3SvcKey{}, mockS3))
		RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.OptsKey{}, rootOpts))
		RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.ConfirmRunnerKey{}, tc.PromptRunner))
		RemoveCmd.SetArgs(tc.args)

		err := RemoveCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		tagOpts.SetZeroValues()
	}
}
