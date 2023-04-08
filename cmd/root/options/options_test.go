package options

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

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
	cmd := cobra.Command{}
	opts := GetRootOptions()
	err := opts.SetAccessCredentialsFromEnv(&cmd)
	assert.Nil(t, err)
}

func TestRootOptions_SetAccessCredentialsFromEnv_Filled(t *testing.T) {
	cmd := cobra.Command{}
	opts := GetRootOptions()

	err := os.Setenv("AWS_REGION", "us-east-1")
	assert.Nil(t, err)
	err = opts.SetAccessCredentialsFromEnv(&cmd)
	assert.Nil(t, err)
	err = os.Setenv("AWS_REGION", "")
	assert.Nil(t, err)

	err = os.Setenv("AWS_ACCESS_KEY", "xxxxx")
	assert.Nil(t, err)
	err = opts.SetAccessCredentialsFromEnv(&cmd)
	assert.Nil(t, err)
	err = os.Setenv("AWS_ACCESS_KEY", "")
	assert.Nil(t, err)

	err = os.Setenv("AWS_SECRET_KEY", "xxxxx")
	assert.Nil(t, err)
	err = opts.SetAccessCredentialsFromEnv(&cmd)
	assert.Nil(t, err)
	err = os.Setenv("AWS_SECRET_KEY", "")
	assert.Nil(t, err)

	err = os.Setenv("AWS_BUCKET_NAME", "xxxxx")
	assert.Nil(t, err)
	err = opts.SetAccessCredentialsFromEnv(&cmd)
	assert.Nil(t, err)
	err = os.Setenv("AWS_BUCKET_NAME", "")
	assert.Nil(t, err)
}
