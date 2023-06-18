package add

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/stretchr/testify/assert"
)

var (
	defaultPutBucketPolicyOutput = &s3.PutBucketPolicyOutput{}
	defaultPutBucketPolicyErr    error
)

func createSvc(rootOpts *options.RootOptions) (*s3.S3, error) {
	return internalaws.CreateAwsService(rootOpts)
}

// Define a mock struct to be used in your unit tests
type mockS3Client struct {
	s3iface.S3API
}

func (m *mockS3Client) PutBucketPolicy(input *s3.PutBucketPolicyInput) (*s3.PutBucketPolicyOutput, error) {
	return defaultPutBucketPolicyOutput, defaultPutBucketPolicyErr
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
		caseName              string
		args                  []string
		shouldPass            bool
		shouldMock            bool
		putBucketPolicyErr    error
		putBucketPolicyOutput *s3.PutBucketPolicyOutput
	}{
		{"Too many arguments", []string{"enabled", "foo"}, false, false,
			nil, &s3.PutBucketPolicyOutput{},
		},
		{"No arguments", []string{}, false, false,
			nil, &s3.PutBucketPolicyOutput{},
		},
		{"Success", []string{"../../../mock/bucketpolicy.json"}, true, true,
			nil, &s3.PutBucketPolicyOutput{},
		},
		{"Failure", []string{"../../../mock/bucketpolicy.json"}, false, true,
			errors.New("dummy error"), &s3.PutBucketPolicyOutput{},
		},
		{"Failure target file not found", []string{"../../../mock/bucketpolicy.jsonjjnnn"}, false,
			true, nil, &s3.PutBucketPolicyOutput{},
		},
	}

	for _, tc := range cases {
		defaultPutBucketPolicyErr = tc.putBucketPolicyErr
		defaultPutBucketPolicyOutput = tc.putBucketPolicyOutput

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
	bucketPolicyOpts.SetZeroValues()
}

/*func TestExecuteTooManyArguments(t *testing.T) {
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
*/
/*func TestExecuteNoArgument(t *testing.T) {
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
*/
/*func TestExecuteFailureFileNotFound(t *testing.T) {
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

	AddCmd.SetArgs([]string{"dummy.json"})
	err := AddCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	bucketPolicyOpts.SetZeroValues()
}
*/
/*func TestExecuteSuccess(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	AddCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultPutBucketPolicyErr = nil

	AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.S3SvcKey{}, svc))
	AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.OptsKey{}, rootOpts))

	AddCmd.SetArgs([]string{"../../../mock/bucketpolicy.json"})
	err := AddCmd.Execute()
	assert.Nil(t, err)

	rootOpts.SetZeroValues()
	bucketPolicyOpts.SetZeroValues()
}
*/
/*func TestExecutePutError(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	ctx := context.Background()
	AddCmd.SetContext(ctx)

	mockSvc := &mockS3Client{}
	svc = mockSvc

	defaultPutBucketPolicyErr = errors.New("dummy error")

	AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.S3SvcKey{}, svc))
	AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.OptsKey{}, rootOpts))

	AddCmd.SetArgs([]string{"../../../mock/bucketpolicy.json"})
	err := AddCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	bucketPolicyOpts.SetZeroValues()
}
*/
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
