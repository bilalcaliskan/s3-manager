//go:build e2e

package remove

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
	defaultPutBucketTaggingErr    error
	defaultPutBucketTaggingOutput = &s3.PutBucketTaggingOutput{}

	defaultGetBucketTaggingErr    error
	defaultGetBucketTaggingOutput = &s3.GetBucketTaggingOutput{}

	defaultDeleteBucketTaggingErr    error
	defaultDeleteBucketTaggingOutput = &s3.DeleteBucketTaggingOutput{}
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
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	RemoveCmd.SetContext(ctx)

	cases := []struct {
		caseName                  string
		args                      []string
		shouldPass                bool
		shouldMock                bool
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
		{"No arguments provided", []string{}, false, false,
			nil,
			&s3.GetBucketTaggingOutput{},
			nil,
			&s3.PutBucketTaggingOutput{},
			nil,
			&s3.DeleteBucketTaggingOutput{},
			nil, false, false,
		},
		{"Success while has single tag", []string{"foo=bar,foo2=bar2"}, true, true,
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
			}, false, false,
		},
		{"Success when dry-run enabled", []string{"foo=bar,foo2=bar2"}, true, true,
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
			nil, true, false,
		},
		{"Success when auto-approve enabled", []string{"foo=bar,foo2=bar2"}, true, true,
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
			nil, false, true,
		},
		{"Success while has multiple tags", []string{"foo=bar,foo2=bar2"}, true, true,
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
			}, false, false,
		},
		{"Failure caused by wrong argument provided", []string{"foo=bar=bar3,foo2=bar2"}, false, true,
			nil,
			&s3.GetBucketTaggingOutput{},
			nil,
			&s3.PutBucketTaggingOutput{},
			nil,
			&s3.DeleteBucketTaggingOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			}, false, false,
		},
		{"Warn while has no tags to remove", []string{"foo3=bar3,foo4=bar4"}, true, true,
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
			}, false, false,
		},
		{"Failure caused by GetBucketTags error", []string{"foo=bar,foo2=bar2"}, false, true,
			errors.New("injected error"),
			&s3.GetBucketTaggingOutput{},
			nil,
			&s3.PutBucketTaggingOutput{},
			errors.New("injected error"),
			&s3.DeleteBucketTaggingOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			}, false, false,
		},
		{"Failure caused by DeleteAllBucketTags error", []string{"foo=bar,foo2=bar2"}, false, true,
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
			errors.New("injected error"),
			&s3.DeleteBucketTaggingOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			}, false, false,
		},
		{"Failure caused by prompt error", []string{"foo=bar,foo2=bar2"}, false, true,
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
				err: errors.New("injected error"),
			}, false, false,
		},
		{"Failure caused by user terminated the process", []string{"foo=bar,foo2=bar2"}, false, true,
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
				err: nil,
			}, false, false,
		},
		{"Failure caused by DeleteAllBucketTags error", []string{"foo=bar,foo2=bar2"}, false, true,
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
			errors.New("injected error"),
			&s3.PutBucketTaggingOutput{},
			nil,
			&s3.DeleteBucketTaggingOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			}, false, false,
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

	rootOpts.SetZeroValues()
	tagOpts.SetZeroValues()
}
