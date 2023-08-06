//go:build e2e

package remove

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/aws/types"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

func TestExecuteRemoveCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()

	ctx := context.Background()
	RemoveCmd.SetContext(ctx)

	cases := []struct {
		caseName                string
		args                    []string
		shouldPass              bool
		getBucketTaggingFunc    func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error)
		putBucketTaggingFunc    func(ctx context.Context, params *s3.PutBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.PutBucketTaggingOutput, error)
		deleteBucketTaggingFunc func(ctx context.Context, params *s3.DeleteBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketTaggingOutput, error)
		prompt.PromptRunner
		dryRun      bool
		autoApprove bool
	}{
		{
			"No arguments provided",
			[]string{},
			false,
			internalawstypes.DefaultGetBucketTaggingFunc,
			internalawstypes.DefaultPutBucketTaggingFunc,
			internalawstypes.DefaultDeleteBucketTaggingFunc,
			nil,
			false,
			false,
		},
		{
			"Success while has single tag",
			[]string{"foo=bar,foo2=bar2"},
			true,
			func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
				return &s3.GetBucketTaggingOutput{
					TagSet: []types.Tag{
						{
							Key:   aws.String("foo"),
							Value: aws.String("bar"),
						},
					},
				}, nil
			},
			internalawstypes.DefaultPutBucketTaggingFunc,
			internalawstypes.DefaultDeleteBucketTaggingFunc,
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
			func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
				return &s3.GetBucketTaggingOutput{
					TagSet: []types.Tag{
						{
							Key:   aws.String("foo"),
							Value: aws.String("bar"),
						},
					},
				}, nil
			},
			internalawstypes.DefaultPutBucketTaggingFunc,
			internalawstypes.DefaultDeleteBucketTaggingFunc,
			nil,
			true,
			false,
		},
		{
			"Success when auto-approve enabled",
			[]string{"foo=bar,foo2=bar2"},
			true,
			func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
				return &s3.GetBucketTaggingOutput{
					TagSet: []types.Tag{
						{
							Key:   aws.String("foo"),
							Value: aws.String("bar"),
						},
					},
				}, nil
			},
			internalawstypes.DefaultPutBucketTaggingFunc,
			internalawstypes.DefaultDeleteBucketTaggingFunc,
			nil,
			false,
			true,
		},
		{
			"Success while has multiple tags",
			[]string{"foo=bar,foo2=bar2"},
			true,
			func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
				return &s3.GetBucketTaggingOutput{
					TagSet: []types.Tag{
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
				}, nil
			},
			internalawstypes.DefaultPutBucketTaggingFunc,
			internalawstypes.DefaultDeleteBucketTaggingFunc,
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
			internalawstypes.DefaultGetBucketTaggingFunc,
			internalawstypes.DefaultPutBucketTaggingFunc,
			internalawstypes.DefaultDeleteBucketTaggingFunc,
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
			func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
				return &s3.GetBucketTaggingOutput{
					TagSet: []types.Tag{
						{
							Key:   aws.String("foo"),
							Value: aws.String("bar"),
						},
						{
							Key:   aws.String("foo2"),
							Value: aws.String("bar2"),
						},
					},
				}, nil
			},
			internalawstypes.DefaultPutBucketTaggingFunc,
			internalawstypes.DefaultDeleteBucketTaggingFunc,
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
			func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
				return nil, constants.ErrInjected
			},
			internalawstypes.DefaultPutBucketTaggingFunc,
			internalawstypes.DefaultDeleteBucketTaggingFunc,
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
			func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
				return &s3.GetBucketTaggingOutput{
					TagSet: []types.Tag{
						{
							Key:   aws.String("foo"),
							Value: aws.String("bar"),
						},
						{
							Key:   aws.String("foo2"),
							Value: aws.String("bar2"),
						},
					},
				}, nil
			},
			internalawstypes.DefaultPutBucketTaggingFunc,
			func(ctx context.Context, params *s3.DeleteBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketTaggingOutput, error) {
				return nil, constants.ErrInjected
			},
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
			func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
				return &s3.GetBucketTaggingOutput{
					TagSet: []types.Tag{
						{
							Key:   aws.String("foo"),
							Value: aws.String("bar"),
						},
						{
							Key:   aws.String("foo2"),
							Value: aws.String("bar2"),
						},
					},
				}, nil
			},
			internalawstypes.DefaultPutBucketTaggingFunc,
			internalawstypes.DefaultDeleteBucketTaggingFunc,
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
			func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
				return &s3.GetBucketTaggingOutput{
					TagSet: []types.Tag{
						{
							Key:   aws.String("foo"),
							Value: aws.String("bar"),
						},
						{
							Key:   aws.String("foo2"),
							Value: aws.String("bar2"),
						},
					},
				}, nil
			},
			internalawstypes.DefaultPutBucketTaggingFunc,
			internalawstypes.DefaultDeleteBucketTaggingFunc,
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
			func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
				return &s3.GetBucketTaggingOutput{
					TagSet: []types.Tag{
						{
							Key:   aws.String("foo"),
							Value: aws.String("bar"),
						},
						{
							Key:   aws.String("foo2"),
							Value: aws.String("bar2"),
						},
					},
				}, nil
			},
			func(ctx context.Context, params *s3.PutBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.PutBucketTaggingOutput, error) {
				return nil, constants.ErrInjected
			},
			internalawstypes.DefaultDeleteBucketTaggingFunc,
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

		mockS3 := new(internalawstypes.MockS3v2Client)
		mockS3.GetBucketTaggingAPI = tc.getBucketTaggingFunc
		mockS3.PutBucketTaggingAPI = tc.putBucketTaggingFunc
		mockS3.DeleteBucketTaggingAPI = tc.deleteBucketTaggingFunc

		rootOpts.DryRun = tc.dryRun
		rootOpts.AutoApprove = tc.autoApprove

		RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.S3ClientKey{}, mockS3))
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
