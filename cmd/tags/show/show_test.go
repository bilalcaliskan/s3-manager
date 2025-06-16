//go:build e2e

package show

import (
	"context"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/pkg/aws/types"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/constants"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

func TestExecuteShowCmd(t *testing.T) {
	rootOpts := options.GetMockedRootOptions()

	ctx := context.Background()
	ShowCmd.SetContext(ctx)

	cases := []struct {
		caseName             string
		args                 []string
		shouldPass           bool
		getBucketTaggingFunc func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error)
	}{
		{
			"Too many arguments",
			[]string{"enabled", "foo"},
			false,
			internalawstypes.DefaultGetBucketTaggingFunc,
		},
		{
			"Success with empty TagSet",
			[]string{},
			true,
			internalawstypes.DefaultGetBucketTaggingFunc,
		},
		{
			"Success with non-empty TagSet",
			[]string{},
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
		},
		{
			"Failure",
			[]string{},
			false,
			func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
				return nil, constants.ErrInjected
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(internalawstypes.MockS3Client)
		mockS3.GetBucketTaggingAPI = tc.getBucketTaggingFunc

		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3ClientKey{}, mockS3))
		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))
		ShowCmd.SetArgs(tc.args)

		err := ShowCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		tagOpts.SetZeroValues()
	}
}
