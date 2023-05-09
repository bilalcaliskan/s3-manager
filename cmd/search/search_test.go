package search

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/internal/aws"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"
)

func createSvc(rootOpts *options.RootOptions) (*s3.S3, error) {
	return aws.CreateAwsService(rootOpts)
}

func TestExecuteMissingRegion(t *testing.T) {
	//err = CleanCmd.PersistentFlags().Set("verbose", "true")
	//assert.Nil(t, err)

	rootOpts := options.GetRootOptions()

	ctx := context.Background()
	SearchCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.Nil(t, svc)
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

	ctx := context.Background()
	SearchCmd.SetContext(ctx)
	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	SearchCmd.SetContext(context.WithValue(SearchCmd.Context(), options.S3SvcKey{}, svc))
	SearchCmd.SetContext(context.WithValue(SearchCmd.Context(), options.OptsKey{}, rootOpts))

	err = SearchCmd.Execute()
	assert.NotNil(t, err)

	rootOpts.SetZeroValues()
	searchOpts.SetZeroValues()
}
