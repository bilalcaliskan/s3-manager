package search

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
	SearchCmd.SetContext(context.WithValue(context.Background(), options.OptsKey{}, rootOpts))
	err := SearchCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	searchOpts.SetZeroValues()
}

func TestExecute(t *testing.T) {
	//err = CleanCmd.PersistentFlags().Set("verbose", "true")
	//assert.Nil(t, err)

	rootOpts := options.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"
	SearchCmd.SetContext(context.WithValue(context.Background(), options.OptsKey{}, rootOpts))
	err := SearchCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	searchOpts.SetZeroValues()
}
