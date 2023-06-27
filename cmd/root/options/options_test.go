//go:build unit

package options

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// TODO: uncomment when interactivity enabled again
/*type promptMock struct {
	msg string
	err error
}

func (p promptMock) Run() (string, error) {
	// return expected result
	return p.msg, p.err
}*/

func TestGetRootOptions(t *testing.T) {
	opts := GetRootOptions()
	assert.NotNil(t, opts)
}

func TestRootOptions_InitFlags(t *testing.T) {
	cmd := cobra.Command{}
	opts := GetRootOptions()
	opts.InitFlags(&cmd)
}

func TestRootOptions_SetAccessCredentialsFromEnv(t *testing.T) {
	opts := GetRootOptions()
	err := opts.SetAccessCredentialsFromEnv()
	assert.Nil(t, err)
}

func TestRootOptions_SetAccessFlagsRequired(t *testing.T) {
	cmd := &cobra.Command{}
	opts := GetRootOptions()
	opts.SetZeroValues()

	opts.SetAccessFlagsRequired(cmd)
}

// TODO: uncomment when interactivity enabled again
/*func TestRootOptions_PromptAccessCredentials_AllSuccess(t *testing.T) {
	accessKeyRunnerOrg := AccessKeyRunner
	AccessKeyRunner = promptMock{
		msg: "adlskfjasldkf",
		err: nil,
	}

	secretKeyRunnerOrg := SecretKeyRunner
	SecretKeyRunner = promptMock{
		msg: "adlskfjasldkf",
		err: nil,
	}

	bucketRunnerOrg := BucketRunner
	BucketRunner = promptMock{
		msg: "adlskfjasldkf",
		err: nil,
	}

	regionRunnerOrg := RegionRunner
	RegionRunner = promptMock{
		msg: "adlskfjasldkf",
		err: nil,
	}

	opts := GetRootOptions()

	err := opts.PromptAccessCredentials(AccessKeyRunner, SecretKeyRunner, BucketRunner, RegionRunner)
	assert.Nil(t, err)

	opts.SetZeroValues()
	AccessKeyRunner = accessKeyRunnerOrg
	SecretKeyRunner = secretKeyRunnerOrg
	RegionRunner = regionRunnerOrg
	BucketRunner = bucketRunnerOrg
}*/

// TODO: uncomment when interactivity enabled again
/*func TestRootOptions_PromptAccessCredentials_AccessKeyFailure(t *testing.T) {
	accessKeyRunnerOrg := AccessKeyRunner
	AccessKeyRunner = promptMock{
		msg: "",
		err: errors.New("asdlfkjasdf"),
	}

	secretKeyRunnerOrg := SecretKeyRunner
	SecretKeyRunner = promptMock{
		msg: "adlskfjasldkf",
		err: nil,
	}

	bucketRunnerOrg := BucketRunner
	BucketRunner = promptMock{
		msg: "adlskfjasldkf",
		err: nil,
	}

	regionRunnerOrg := RegionRunner
	RegionRunner = promptMock{
		msg: "adlskfjasldkf",
		err: nil,
	}

	opts := GetRootOptions()

	err := opts.PromptAccessCredentials(AccessKeyRunner, SecretKeyRunner, BucketRunner, RegionRunner)
	assert.NotNil(t, err)

	opts.SetZeroValues()
	AccessKeyRunner = accessKeyRunnerOrg
	SecretKeyRunner = secretKeyRunnerOrg
	RegionRunner = regionRunnerOrg
	BucketRunner = bucketRunnerOrg
}*/

// TODO: uncomment when interactivity enabled again
/*func TestRootOptions_PromptAccessCredentials_SecretKeyFailure(t *testing.T) {
	accessKeyRunnerOrg := AccessKeyRunner
	AccessKeyRunner = promptMock{
		msg: "dsafasdfdfs",
		err: nil,
	}

	secretKeyRunnerOrg := SecretKeyRunner
	SecretKeyRunner = promptMock{
		msg: "",
		err: errors.New("adlskfjasldkf"),
	}

	bucketRunnerOrg := BucketRunner
	BucketRunner = promptMock{
		msg: "adlskfjasldkf",
		err: nil,
	}

	regionRunnerOrg := RegionRunner
	RegionRunner = promptMock{
		msg: "adlskfjasldkf",
		err: nil,
	}

	opts := GetRootOptions()

	err := opts.PromptAccessCredentials(AccessKeyRunner, SecretKeyRunner, BucketRunner, RegionRunner)
	assert.NotNil(t, err)

	opts.SetZeroValues()
	AccessKeyRunner = accessKeyRunnerOrg
	SecretKeyRunner = secretKeyRunnerOrg
	RegionRunner = regionRunnerOrg
	BucketRunner = bucketRunnerOrg
}*/

// TODO: uncomment when interactivity enabled again
/*func TestRootOptions_PromptAccessCredentials_BucketFailure(t *testing.T) {
	accessKeyRunnerOrg := AccessKeyRunner
	AccessKeyRunner = promptMock{
		msg: "dsafasdfdfs",
		err: nil,
	}

	secretKeyRunnerOrg := SecretKeyRunner
	SecretKeyRunner = promptMock{
		msg: "asdfasdf",
		err: nil,
	}

	bucketRunnerOrg := BucketRunner
	BucketRunner = promptMock{
		msg: "",
		err: errors.New("adlskfjasldkf"),
	}

	regionRunnerOrg := RegionRunner
	RegionRunner = promptMock{
		msg: "adlskfjasldkf",
		err: nil,
	}

	opts := GetRootOptions()

	err := opts.PromptAccessCredentials(AccessKeyRunner, SecretKeyRunner, BucketRunner, RegionRunner)
	assert.NotNil(t, err)

	opts.SetZeroValues()
	AccessKeyRunner = accessKeyRunnerOrg
	SecretKeyRunner = secretKeyRunnerOrg
	RegionRunner = regionRunnerOrg
	BucketRunner = bucketRunnerOrg
}*/

// TODO: uncomment when interactivity enabled again
/*func TestRootOptions_PromptAccessCredentials_RegionFailure(t *testing.T) {
	accessKeyRunnerOrg := AccessKeyRunner
	AccessKeyRunner = promptMock{
		msg: "dsafasdfdfs",
		err: nil,
	}

	secretKeyRunnerOrg := SecretKeyRunner
	SecretKeyRunner = promptMock{
		msg: "asdfasdf",
		err: nil,
	}

	bucketRunnerOrg := BucketRunner
	BucketRunner = promptMock{
		msg: "asdfasdfasfd",
		err: nil,
	}

	regionRunnerOrg := RegionRunner
	RegionRunner = promptMock{
		msg: "",
		err: errors.New("adlskfjasldkf"),
	}

	opts := GetRootOptions()

	err := opts.PromptAccessCredentials(AccessKeyRunner, SecretKeyRunner, BucketRunner, RegionRunner)
	assert.NotNil(t, err)

	opts.SetZeroValues()
	AccessKeyRunner = accessKeyRunnerOrg
	SecretKeyRunner = secretKeyRunnerOrg
	RegionRunner = regionRunnerOrg
	BucketRunner = bucketRunnerOrg
}*/

func TestRootOptions_SetAccessCredentialsFromEnv_Filled(t *testing.T) {
	opts := GetRootOptions()

	err := os.Setenv("AWS_REGION", "us-east-1")
	assert.Nil(t, err)
	err = opts.SetAccessCredentialsFromEnv()
	assert.Nil(t, err)
	err = os.Setenv("AWS_REGION", "")
	assert.Nil(t, err)

	err = os.Setenv("AWS_ACCESS_KEY", "xxxxx")
	assert.Nil(t, err)
	err = opts.SetAccessCredentialsFromEnv()
	assert.Nil(t, err)
	err = os.Setenv("AWS_ACCESS_KEY", "")
	assert.Nil(t, err)

	err = os.Setenv("AWS_SECRET_KEY", "xxxxx")
	assert.Nil(t, err)
	err = opts.SetAccessCredentialsFromEnv()
	assert.Nil(t, err)
	err = os.Setenv("AWS_SECRET_KEY", "")
	assert.Nil(t, err)

	err = os.Setenv("AWS_BUCKET_NAME", "xxxxx")
	assert.Nil(t, err)
	err = opts.SetAccessCredentialsFromEnv()
	assert.Nil(t, err)
	err = os.Setenv("AWS_BUCKET_NAME", "")
	assert.Nil(t, err)
}
