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

// Define a mock struct to be used in your unit tests
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

func TestExecuteNoArgumentsProvided(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	RemoveCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.S3SvcKey{}, svc))
	RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.OptsKey{}, rootOpts))

	RemoveCmd.SetArgs([]string{})

	err = RemoveCmd.Execute()
	assert.NotNil(t, err)

	tagOpts.SetZeroValues()
}

func TestExecuteTooManyArgumentsProvided(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	RemoveCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.S3SvcKey{}, svc))
	RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.OptsKey{}, rootOpts))

	RemoveCmd.SetArgs([]string{"foo", "bar"})

	err = RemoveCmd.Execute()
	assert.NotNil(t, err)

	tagOpts.SetZeroValues()
}

func TestExecuteWrongFormattedArgumentProvided(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	RemoveCmd.SetContext(ctx)
	mockSvc := &mockS3Client{}
	svc = mockSvc

	RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.S3SvcKey{}, svc))
	RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.OptsKey{}, rootOpts))

	RemoveCmd.SetArgs([]string{"foo=bar,foo2rtwegrg"})

	err := RemoveCmd.Execute()
	assert.NotNil(t, err)

	tagOpts.SetZeroValues()
}

func TestExecuteGetTagsFailure(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	RemoveCmd.SetContext(ctx)
	mockSvc := &mockS3Client{}
	svc = mockSvc

	RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.S3SvcKey{}, svc))
	RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.OptsKey{}, rootOpts))

	defaultGetBucketTaggingErr = errors.New("dummy error")

	RemoveCmd.SetArgs([]string{"foo=bar"})

	err := RemoveCmd.Execute()
	assert.NotNil(t, err)

	tagOpts.SetZeroValues()
}

func TestExecuteSuccess(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	RemoveCmd.SetContext(ctx)
	mockSvc := &mockS3Client{}
	svc = mockSvc

	RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.S3SvcKey{}, svc))
	RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.OptsKey{}, rootOpts))

	var tags []*s3.Tag
	tags = append(tags, &s3.Tag{Key: aws.String("foo"), Value: aws.String("bar")})
	tags = append(tags, &s3.Tag{Key: aws.String("hasan"), Value: aws.String("huseyin")})
	defaultGetBucketTaggingOutput = &s3.GetBucketTaggingOutput{TagSet: tags}
	defaultGetBucketTaggingErr = nil
	defaultPutBucketTaggingErr = nil
	defaultDeleteBucketTaggingErr = nil

	RemoveCmd.SetArgs([]string{"foo=bar"})
	err := RemoveCmd.Execute()
	assert.Nil(t, err)

	tagOpts.SetZeroValues()
}

func TestExecuteNoTagsToDelete(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	RemoveCmd.SetContext(ctx)
	mockSvc := &mockS3Client{}
	svc = mockSvc

	RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.S3SvcKey{}, svc))
	RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.OptsKey{}, rootOpts))

	var tags []*s3.Tag
	tags = append(tags, &s3.Tag{Key: aws.String("foo"), Value: aws.String("bar")})
	defaultGetBucketTaggingOutput = &s3.GetBucketTaggingOutput{TagSet: tags}
	defaultGetBucketTaggingErr = nil
	defaultPutBucketTaggingErr = nil
	defaultDeleteBucketTaggingErr = nil

	RemoveCmd.SetArgs([]string{"foo2=bar"})
	err := RemoveCmd.Execute()
	assert.Nil(t, err)

	tagOpts.SetZeroValues()
}

func TestExecuteDeleteFailure(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	RemoveCmd.SetContext(ctx)
	mockSvc := &mockS3Client{}
	svc = mockSvc

	RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.S3SvcKey{}, svc))
	RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.OptsKey{}, rootOpts))

	var tags []*s3.Tag
	tags = append(tags, &s3.Tag{Key: aws.String("foo"), Value: aws.String("bar")})
	defaultGetBucketTaggingOutput = &s3.GetBucketTaggingOutput{TagSet: tags}
	defaultGetBucketTaggingErr = nil
	defaultPutBucketTaggingErr = nil
	defaultDeleteBucketTaggingErr = errors.New("dummy error")

	RemoveCmd.SetArgs([]string{"foo=bar"})
	err := RemoveCmd.Execute()
	assert.NotNil(t, err)

	tagOpts.SetZeroValues()
}

func TestExecuteUpdateFailure(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	RemoveCmd.SetContext(ctx)
	mockSvc := &mockS3Client{}
	svc = mockSvc

	RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.S3SvcKey{}, svc))
	RemoveCmd.SetContext(context.WithValue(RemoveCmd.Context(), options.OptsKey{}, rootOpts))

	var tags []*s3.Tag
	tags = append(tags, &s3.Tag{Key: aws.String("foo"), Value: aws.String("bar")})
	defaultGetBucketTaggingOutput = &s3.GetBucketTaggingOutput{TagSet: tags}
	defaultGetBucketTaggingErr = nil
	defaultPutBucketTaggingErr = errors.New("dummy error")
	defaultDeleteBucketTaggingErr = nil

	RemoveCmd.SetArgs([]string{"foo=bar"})
	err := RemoveCmd.Execute()
	assert.NotNil(t, err)

	tagOpts.SetZeroValues()
}

/*func TestExecuteTooManyArgumentsProvided(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	AddCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.S3SvcKey{}, svc))
	AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.OptsKey{}, rootOpts))

	args := []string{"hello", "bar"}
	AddCmd.SetArgs(args)

	err = AddCmd.Execute()
	assert.NotNil(t, err)

	tagOpts.SetZeroValues()
}

func TestExecuteSuccess(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	AddCmd.SetContext(ctx)
	mockSvc := &mockS3Client{}
	svc = mockSvc

	AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.S3SvcKey{}, svc))
	AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.OptsKey{}, rootOpts))

	args := []string{"foo=bar"}
	AddCmd.SetArgs(args)

	defaultGetBucketTaggingErr = nil
	var tags []*s3.Tag
	tags = append(tags, &s3.Tag{Key: aws.String("foo"), Value: aws.String("bar")})

	defaultGetBucketTaggingOutput = &s3.GetBucketTaggingOutput{TagSet: tags}
	defaultPutBucketTaggingErr = nil

	err := AddCmd.Execute()
	assert.Nil(t, err)

	tagOpts.SetZeroValues()
}

func TestExecuteGetBucketTagsFailure(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	AddCmd.SetContext(ctx)
	mockSvc := &mockS3Client{}
	svc = mockSvc

	AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.S3SvcKey{}, svc))
	AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.OptsKey{}, rootOpts))

	args := []string{"foo=bar"}
	AddCmd.SetArgs(args)

	defaultGetBucketTaggingErr = errors.New("dummy error")

	err := AddCmd.Execute()
	assert.NotNil(t, err)

	tagOpts.SetZeroValues()
}

func TestExecuteSetBucketTagsFailure(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	AddCmd.SetContext(ctx)
	mockSvc := &mockS3Client{}
	svc = mockSvc

	AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.S3SvcKey{}, svc))
	AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.OptsKey{}, rootOpts))

	args := []string{"foo=bar"}
	AddCmd.SetArgs(args)

	defaultGetBucketTaggingErr = nil
	defaultPutBucketTaggingErr = errors.New("dummy error")

	err := AddCmd.Execute()
	assert.NotNil(t, err)

	tagOpts.SetZeroValues()
}*/

/*func TestExecuteWrongArguments(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	AddCmd.SetContext(ctx)
	mockSvc := &mockS3Client{}
	svc = mockSvc

	AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.S3SvcKey{}, svc))
	AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.OptsKey{}, rootOpts))

	args := []string{"foo=bar=asdfasdf"}
	AddCmd.SetArgs(args)

	defaultGetBucketTaggingErr = nil

	err := AddCmd.Execute()
	assert.NotNil(t, err)

	tagOpts.SetZeroValues()
}
*/
