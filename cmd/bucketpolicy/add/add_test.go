package add

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/stretchr/testify/assert"
)

/*var (
	defaultGetBucketPolicyOutput = &s3.GetBucketPolicyOutput{
		Policy: aws.String("{}"),
	}
	defaultGetBucketPolicyErr error
)*/

func createSvc(rootOpts *options.RootOptions) (*s3.S3, error) {
	return internalaws.CreateAwsService(rootOpts)
}

// Define a mock struct to be used in your unit tests
/*type mockS3Client struct {
	s3iface.S3API
}*/

/*func (m *mockS3Client) GetBucketPolicy(input *s3.GetBucketPolicyInput) (*s3.GetBucketPolicyOutput, error) {
	return defaultGetBucketPolicyOutput, defaultGetBucketPolicyErr
}*/

func TestExecuteTooManyArguments(t *testing.T) {
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

	args := []string{"enabled", "foo"}
	AddCmd.SetArgs(args)

	err = AddCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	bucketPolicyOpts.SetZeroValues()
}

func TestExecuteNoArgument(t *testing.T) {
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

	AddCmd.SetArgs([]string{})
	err = AddCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	bucketPolicyOpts.SetZeroValues()
}

/*func TestExecuteSuccessEnabled(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	AddCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketPolicyErr = nil

	AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.S3SvcKey{}, svc))
	AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.OptsKey{}, rootOpts))

	AddCmd.SetArgs([]string{"dummy.json"})
	err := AddCmd.Execute()
	assert.Nil(t, err)

	rootOpts.SetZeroValues()
	bucketPolicyOpts.SetZeroValues()
}*/

/*func TestExecuteSuccessEnabled2(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	AddCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultGetBucketPolicyErr = nil

	AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.S3SvcKey{}, svc))
	AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.OptsKey{}, rootOpts))

	AddCmd.SetArgs([]string{})
	err := AddCmd.Execute()
	assert.Nil(t, err)

	rootOpts.SetZeroValues()
	bucketPolicyOpts.SetZeroValues()
}
*/
