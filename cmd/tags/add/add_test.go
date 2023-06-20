package add

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
)

func createSvc(rootOpts *options.RootOptions) (*s3.S3, error) {
	return internalaws.CreateAwsService(rootOpts)
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
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	AddCmd.SetContext(ctx)

	cases := []struct {
		caseName               string
		args                   []string
		shouldPass             bool
		shouldMock             bool
		getBucketTaggingErr    error
		getBucketTaggingOutput *s3.GetBucketTaggingOutput
		putBucketTaggingErr    error
		putBucketTaggingOutput *s3.PutBucketTaggingOutput
	}{
		{"No arguments provided", []string{}, false, false,
			nil,
			&s3.GetBucketTaggingOutput{},
			nil,
			&s3.PutBucketTaggingOutput{},
		},
		{"Success", []string{"foo=bar,foo2=bar2"}, true, true,
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
		},
		{"Failure caused by GetBucketTags error", []string{"foo=bar,foo2=bar2"}, false, true,
			errors.New("injected error"),
			&s3.GetBucketTaggingOutput{},
			nil,
			&s3.PutBucketTaggingOutput{},
		},
		{"Failure caused by wrong provided arg", []string{"foo=bar=barX,foo2=bar2"}, false, true,
			nil,
			&s3.GetBucketTaggingOutput{},
			nil,
			&s3.PutBucketTaggingOutput{},
		},
		{"Failure caused by SetBucketTags error", []string{"foo=bar,foo2=bar2"}, false, true,
			nil,
			&s3.GetBucketTaggingOutput{},
			errors.New("injected error"),
			&s3.PutBucketTaggingOutput{},
		},
	}

	for _, tc := range cases {
		defaultGetBucketTaggingErr = tc.getBucketTaggingErr
		defaultPutBucketTaggingErr = tc.putBucketTaggingErr
		defaultGetBucketTaggingOutput = tc.getBucketTaggingOutput
		defaultPutBucketTaggingOutput = tc.putBucketTaggingOutput

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

	rootOpts.SetZeroValues()
	tagOpts.SetZeroValues()
}
