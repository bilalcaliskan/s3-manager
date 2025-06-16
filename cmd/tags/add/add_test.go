//go:build e2e

package add

import (
	"context"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/pkg/aws/types"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/constants"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/prompt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

func TestExecuteAddCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()

	ctx := context.Background()
	AddCmd.SetContext(ctx)

	cases := []struct {
		caseName             string
		args                 []string
		shouldPass           bool
		getBucketTaggingFunc func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error)
		putBucketTaggingFunc func(ctx context.Context, params *s3.PutBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.PutBucketTaggingOutput, error)
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
			nil,
			false,
			false,
		},
		{
			"Success when auto-approve disabled",
			[]string{"foo=bar,foo2=bar2"},
			true,
			func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
				return &s3.GetBucketTaggingOutput{
					TagSet: []types.Tag{
						{
							Key:   aws.String("hasan1"),
							Value: aws.String("huseyin1"),
						},
						{
							Key:   aws.String("hasan2"),
							Value: aws.String("huseyin2"),
						},
					},
				}, nil
			},
			internalawstypes.DefaultPutBucketTaggingFunc,
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
			func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
				return &s3.GetBucketTaggingOutput{
					TagSet: []types.Tag{
						{
							Key:   aws.String("hasan1"),
							Value: aws.String("huseyin1"),
						},
						{
							Key:   aws.String("hasan2"),
							Value: aws.String("huseyin2"),
						},
					},
				}, nil
			},
			internalawstypes.DefaultPutBucketTaggingFunc,
			nil,
			false,
			true,
		},
		{
			"Success when dry run enabled",
			[]string{"foo=bar,foo2=bar2"},
			true,
			func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
				return &s3.GetBucketTaggingOutput{
					TagSet: []types.Tag{
						{
							Key:   aws.String("hasan1"),
							Value: aws.String("huseyin1"),
						},
						{
							Key:   aws.String("hasan2"),
							Value: aws.String("huseyin2"),
						},
					},
				}, nil
			},
			internalawstypes.DefaultPutBucketTaggingFunc,
			nil,
			true,
			false,
		},
		{
			"Success",
			[]string{"foo=bar,foo2=bar2"},
			true,
			func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
				return &s3.GetBucketTaggingOutput{
					TagSet: []types.Tag{
						{
							Key:   aws.String("hasan1"),
							Value: aws.String("huseyin1"),
						},
						{
							Key:   aws.String("hasan2"),
							Value: aws.String("huseyin2"),
						},
					},
				}, nil
			},
			internalawstypes.DefaultPutBucketTaggingFunc,
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
			internalawstypes.DefaultGetBucketTaggingFunc,
			internalawstypes.DefaultPutBucketTaggingFunc,
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
			internalawstypes.DefaultGetBucketTaggingFunc,
			internalawstypes.DefaultPutBucketTaggingFunc,
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
			internalawstypes.DefaultGetBucketTaggingFunc,
			internalawstypes.DefaultPutBucketTaggingFunc,
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
			internalawstypes.DefaultGetBucketTaggingFunc,
			func(ctx context.Context, params *s3.PutBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.PutBucketTaggingOutput, error) {
				return nil, constants.ErrInjected
			},
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

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.GetBucketTaggingAPI = tc.getBucketTaggingFunc
		mockS3.PutBucketTaggingAPI = tc.putBucketTaggingFunc

		AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.S3ClientKey{}, mockS3))
		AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.ConfirmRunnerKey{}, tc.PromptRunner))
		AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.OptsKey{}, rootOpts))
		AddCmd.SetArgs(tc.args)

		err := AddCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		tagOpts.SetZeroValues()
	}
}
