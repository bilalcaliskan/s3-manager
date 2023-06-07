package utils

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	options2 "github.com/bilalcaliskan/s3-manager/cmd/versioning/options"
	"github.com/stretchr/testify/assert"
)

func TestCheckArgsSuccess(t *testing.T) {
	err := CheckArgs([]string{})
	assert.Nil(t, err)
}

func TestCheckArgsFailure(t *testing.T) {
	err := CheckArgs([]string{"foo"})
	assert.NotNil(t, err)
}

func TestDecideActualStateEnabled(t *testing.T) {
	res := &s3.GetBucketVersioningOutput{
		Status: aws.String("Enabled"),
	}

	rootOpts := options.GetRootOptions()
	opts := options2.GetVersioningOptions()
	opts.RootOptions = rootOpts

	err := DecideActualState(res, opts)
	assert.Nil(t, err)
}

func TestDecideActualStateSuspended(t *testing.T) {
	res := &s3.GetBucketVersioningOutput{
		Status: aws.String("Suspended"),
	}

	rootOpts := options.GetRootOptions()
	opts := options2.GetVersioningOptions()
	opts.RootOptions = rootOpts

	err := DecideActualState(res, opts)
	assert.Nil(t, err)
}

func TestDecideActualStateUndefined(t *testing.T) {
	res := &s3.GetBucketVersioningOutput{
		Status: aws.String("Suspendedddddd"),
	}

	rootOpts := options.GetRootOptions()
	opts := options2.GetVersioningOptions()
	opts.RootOptions = rootOpts

	err := DecideActualState(res, opts)
	assert.NotNil(t, err)
}
