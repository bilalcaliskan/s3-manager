package clean

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
	CleanCmd.SetContext(context.WithValue(context.Background(), options.OptsKey{}, rootOpts))
	err := CleanCmd.Execute()
	assert.NotNil(t, err)
}

func TestExecuteInvalidSortByOption(t *testing.T) {
	err := CleanCmd.Flags().Set("sortBy", "nonexistedsortbyflag")
	assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	t.Log(rootOpts.Region)
	CleanCmd.SetContext(context.WithValue(context.Background(), options.OptsKey{}, rootOpts))
	err = CleanCmd.Execute()
	assert.NotNil(t, err)
}

func TestExecuteInvalidMinMaxValues(t *testing.T) {
	err := CleanCmd.Flags().Set("minFileSizeInMb", "20")
	assert.Nil(t, err)

	err = CleanCmd.Flags().Set("maxFileSizeInMb", "10")
	assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	CleanCmd.SetContext(context.WithValue(context.Background(), options.OptsKey{}, rootOpts))
	err = CleanCmd.Execute()
	assert.NotNil(t, err)
}

func TestExecute(t *testing.T) {
	//err = CleanCmd.PersistentFlags().Set("verbose", "true")
	//assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"
	CleanCmd.SetContext(context.WithValue(context.Background(), options.OptsKey{}, rootOpts))
	err := CleanCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
}
