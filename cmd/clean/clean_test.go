package clean

import (
	"context"
	"fmt"
	"testing"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

func TestExecuteMissingRegion(t *testing.T) {
	//err = CleanCmd.PersistentFlags().Set("verbose", "true")
	//assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	fmt.Println(rootOpts)
	CleanCmd.SetContext(context.WithValue(context.Background(), options.OptsKey{}, rootOpts))
	err := CleanCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	cleanOpts.SetZeroValues()
}

func TestExecuteInvalidSortByOption(t *testing.T) {
	err := CleanCmd.Flags().Set("sortBy", "nonexistedsortbyflag")
	assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	CleanCmd.SetContext(context.WithValue(context.Background(), options.OptsKey{}, rootOpts))
	err = CleanCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	cleanOpts.SetZeroValues()
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

	rootOpts.SetZeroValues()
	cleanOpts.SetZeroValues()
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
	cleanOpts.SetZeroValues()
}
