//go:build unit

package utils

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	options2 "github.com/bilalcaliskan/s3-manager/cmd/versioning/options"
	"github.com/stretchr/testify/assert"
)

// Define a testdata struct to be used in your unit tests
type mockS3Client struct {
	s3iface.S3API
}

func TestDecideActualState(t *testing.T) {
	tests := []struct {
		name     string
		res      *s3.GetBucketVersioningOutput
		expected error
	}{
		{
			name: "Success enabled",
			res: &s3.GetBucketVersioningOutput{
				Status: aws.String("Enabled"),
			},
			expected: nil,
		},
		{
			name: "Success suspended",
			res: &s3.GetBucketVersioningOutput{
				Status: aws.String("Suspended"),
			},
			expected: nil,
		},
		{
			name: "Failure caused by unknown state",
			res: &s3.GetBucketVersioningOutput{
				Status: aws.String("Suspendeddd"),
			},
			expected: fmt.Errorf(ErrUnknownStatus, "Suspendeddd"),
		},
	}

	opts := options2.GetVersioningOptions()
	for _, test := range tests {
		err := DecideActualState(test.res, opts)
		assert.Equal(t, test.expected, err)
	}
}
