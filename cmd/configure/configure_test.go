package configure

import (
	"context"
	"testing"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

func TestExecuteMissingRegion(t *testing.T) {
	//err = CleanCmd.PersistentFlags().Set("verbose", "true")
	//assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = ""
	rootOpts.BucketName = "thisisbucketname"
	rootOpts.Interactive = false
	ConfigureCmd.SetContext(context.WithValue(context.Background(), options.OptsKey{}, rootOpts))
	err := ConfigureCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	configureOpts.SetZeroValues()
}

func TestExecute(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"
	rootOpts.Interactive = false
	ConfigureCmd.SetContext(context.WithValue(context.Background(), options.OptsKey{}, rootOpts))

	err := ConfigureCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	configureOpts.SetZeroValues()
}

func TestExecuteFailingPutRequest(t *testing.T) {
	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"
	rootOpts.Interactive = true
	ConfigureCmd.SetContext(context.WithValue(context.Background(), options.OptsKey{}, rootOpts))

	err := ConfigureCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	configureOpts.SetZeroValues()
}
