package show

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

func (m *mockS3Client) GetBucketTagging(input *s3.GetBucketTaggingInput) (*s3.GetBucketTaggingOutput, error) {
	return defaultGetBucketTaggingOutput, defaultGetBucketTaggingErr
}

func TestExecuteShowCmd(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	ShowCmd.SetContext(ctx)

	cases := []struct {
		caseName                      string
		args                          []string
		shouldPass                    bool
		shouldMock                    bool
		defaultGetBucketTaggingErr    error
		defaultGetBucketTaggingOutput *s3.GetBucketTaggingOutput
	}{
		{"Too many arguments", []string{"enabled", "foo"}, false, false, nil,
			&s3.GetBucketTaggingOutput{},
		},
		{"Success with empty TagSet", []string{}, true, true, nil,
			&s3.GetBucketTaggingOutput{},
		},
		{"Success with non-empty TagSet", []string{}, true, true, nil,
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
		},
		{"Failure", []string{}, false, true, errors.New("dummy error"),
			&s3.GetBucketTaggingOutput{},
		},
	}

	for _, tc := range cases {
		defaultGetBucketTaggingErr = tc.defaultGetBucketTaggingErr
		defaultGetBucketTaggingOutput = tc.defaultGetBucketTaggingOutput

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

		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3SvcKey{}, svc))
		ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))
		ShowCmd.SetArgs(tc.args)

		err = ShowCmd.Execute()

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}

	rootOpts.SetZeroValues()
	tagOpts.SetZeroValues()
}

/*func TestExecuteTooManyArguments(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	ShowCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3SvcKey{}, svc))
	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))

	args := []string{"enabled", "foo"}
	ShowCmd.SetArgs(args)

	err = ShowCmd.Execute()
	assert.NotNil(t, err)

	tagOpts.SetZeroValues()
}

func TestExecuteSuccessEnabled(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	ShowCmd.SetContext(ctx)
	mockSvc := &mockS3Client{}
	svc = mockSvc

	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3SvcKey{}, svc))
	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))

	ShowCmd.SetArgs([]string{})
	defaultGetBucketTaggingErr = nil
	err := ShowCmd.Execute()
	assert.Nil(t, err)

	tagOpts.SetZeroValues()
}

func TestExecuteSuccess(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	ShowCmd.SetContext(ctx)
	mockSvc := &mockS3Client{}
	svc = mockSvc

	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3SvcKey{}, svc))
	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))

	ShowCmd.SetArgs([]string{})
	defaultGetBucketTaggingErr = nil
	var tags []*s3.Tag
	tags = append(tags, &s3.Tag{Key: aws.String("foo"), Value: aws.String("bar")})

	defaultGetBucketTaggingOutput = &s3.GetBucketTaggingOutput{TagSet: tags}
	err := ShowCmd.Execute()
	assert.Nil(t, err)

	tagOpts.SetZeroValues()
}*/

/*func TestExecuteFailure(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	ShowCmd.SetContext(ctx)
	mockSvc := &mockS3Client{}
	svc = mockSvc

	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.S3SvcKey{}, svc))
	ShowCmd.SetContext(context.WithValue(ShowCmd.Context(), options.OptsKey{}, rootOpts))

	ShowCmd.SetArgs([]string{})
	defaultGetBucketTaggingErr = errors.New("dummy error")
	err := ShowCmd.Execute()
	assert.NotNil(t, err)

	tagOpts.SetZeroValues()
}*/
