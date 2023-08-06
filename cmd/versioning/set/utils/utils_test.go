//go:build unit

package utils

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"testing"

	options2 "github.com/bilalcaliskan/s3-manager/cmd/versioning/options"
	"github.com/stretchr/testify/assert"
)

func TestDecideActualState(t *testing.T) {
	tests := []struct {
		name     string
		res      *s3.GetBucketVersioningOutput
		expected error
	}{
		{
			"Success enabled",
			&s3.GetBucketVersioningOutput{
				Status: types.BucketVersioningStatusEnabled,
			},
			nil,
		},
		{
			"Success suspended",
			&s3.GetBucketVersioningOutput{
				Status: types.BucketVersioningStatusSuspended,
			},
			nil,
		},
		{
			"Failure caused by unknown state",
			&s3.GetBucketVersioningOutput{
				Status: "lasdkfjlkasdfjsldf",
			},
			fmt.Errorf(ErrUnknownStatus, "lasdkfjlkasdfjsldf"),
		},
	}

	opts := options2.GetVersioningOptions()
	for _, test := range tests {
		err := DecideActualState(test.res, opts)
		assert.Equal(t, test.expected, err)
	}
}
