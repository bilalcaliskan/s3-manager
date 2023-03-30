package utils

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/stretchr/testify/assert"
)

var mockObjects = []*s3.Object{{
	ChecksumAlgorithm: nil,
	ETag:              aws.String("233b4ce689c7086b7958eb31d8f8b811"),
	Key:               aws.String("bar-service/233b4ce689c7086b7958eb31d8f8b811.template"),
	LastModified:      aws.Time(time.Now()),
	Owner: &s3.Owner{
		DisplayName: aws.String("developer1"),
		ID:          aws.String("3becc289963dfc26fe632e4d2fc78d2c7875fc4f030813629e28db2c1fbba4b7"),
	},
	Size:         aws.Int64(2129),
	StorageClass: aws.String("STANDARD"),
}, {
	ChecksumAlgorithm: nil,
	ETag:              aws.String("233b4ce689c7086b7958eb31d8f8b811"),
	Key:               aws.String("foo-service/233b4ce689c7086b7958eb31d8f8b811.template"),
	LastModified:      aws.Time(time.Now()),
	Owner: &s3.Owner{
		DisplayName: aws.String("developer1"),
		ID:          aws.String("3becc289963dfc26fe632e4d2fc78d2c7875fc4f030813629e28db2c1fbba4b7"),
	},
	Size:         aws.Int64(2129),
	StorageClass: aws.String("STANDARD"),
}}

func TestContains(t *testing.T) {
	res := Contains([]string{"size", "lastModificationDate"}, "size")
	assert.True(t, res)
}

func TestNotContains(t *testing.T) {
	res := Contains([]string{"size", "lastModificationDate"}, "sizee")
	assert.False(t, res)
}

func TestGetKeysOnly(t *testing.T) {
	keys := GetKeysOnly(mockObjects)
	assert.NotEmpty(t, keys)
}
