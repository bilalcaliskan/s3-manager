package types

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestMockS3Client_DeleteObject(t *testing.T) {
	f := func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
		return &s3.DeleteObjectOutput{}, nil
	}

	mock := new(MockS3Client)
	mock.DeleteObjectAPI = f

	res, err := mock.DeleteObject(context.Background(), &s3.DeleteObjectInput{})
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestMockS3Client_GetObject(t *testing.T) {
	f := func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
		return &s3.GetObjectOutput{}, nil
	}

	mock := new(MockS3Client)
	mock.GetObjectAPI = f

	res, err := mock.GetObject(context.Background(), &s3.GetObjectInput{})
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestMockS3Client_ListObjects(t *testing.T) {
	f := func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
		return &s3.ListObjectsOutput{}, nil
	}

	mock := new(MockS3Client)
	mock.ListObjectsAPI = f

	res, err := mock.ListObjects(context.Background(), &s3.ListObjectsInput{})
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestMockS3Client_GetBucketPolicy(t *testing.T) {
	f := func(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
		return &s3.GetBucketPolicyOutput{}, nil
	}

	mock := new(MockS3Client)
	mock.GetBucketPolicyAPI = f

	res, err := mock.GetBucketPolicy(context.Background(), &s3.GetBucketPolicyInput{})
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestMockS3Client_PutBucketPolicy(t *testing.T) {
	f := func(ctx context.Context, params *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error) {
		return &s3.PutBucketPolicyOutput{}, nil
	}

	mock := new(MockS3Client)
	mock.PutBucketPolicyAPI = f

	res, err := mock.PutBucketPolicy(context.Background(), &s3.PutBucketPolicyInput{})
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestMockS3Client_DeleteBucketPolicy(t *testing.T) {
	f := func(ctx context.Context, params *s3.DeleteBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketPolicyOutput, error) {
		return &s3.DeleteBucketPolicyOutput{}, nil
	}

	mock := new(MockS3Client)
	mock.DeleteBucketPolicyAPI = f

	res, err := mock.DeleteBucketPolicy(context.Background(), &s3.DeleteBucketPolicyInput{})
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestMockS3Client_GetBucketAccelerateConfiguration(t *testing.T) {
	f := func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
		return &s3.GetBucketAccelerateConfigurationOutput{}, nil
	}

	mock := new(MockS3Client)
	mock.GetBucketAccelerateConfigurationAPI = f

	res, err := mock.GetBucketAccelerateConfiguration(context.Background(), &s3.GetBucketAccelerateConfigurationInput{})
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestMockS3Client_PutBucketAccelerateConfiguration(t *testing.T) {
	mock := new(MockS3Client)
	mock.PutBucketAccelerateConfigurationAPI = DefaultPutBucketAccelerationFunc

	res, err := mock.PutBucketAccelerateConfiguration(context.Background(), &s3.PutBucketAccelerateConfigurationInput{})
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestMockS3Client_GetBucketVersioning(t *testing.T) {
	mock := new(MockS3Client)
	mock.GetBucketVersioningAPI = DefaultGetBucketVersioningFunc

	res, err := mock.GetBucketVersioning(context.Background(), &s3.GetBucketVersioningInput{})
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestMockS3Client_PutBucketVersioning(t *testing.T) {
	mock := new(MockS3Client)
	mock.PutBucketVersioningAPI = DefaultPutBucketVersioningFunc

	res, err := mock.PutBucketVersioning(context.Background(), &s3.PutBucketVersioningInput{})
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestMockS3Client_GetBucketTagging(t *testing.T) {
	mock := new(MockS3Client)
	mock.GetBucketTaggingAPI = DefaultGetBucketTaggingFunc

	res, err := mock.GetBucketTagging(context.Background(), &s3.GetBucketTaggingInput{})
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestMockS3Client_PutBucketTagging(t *testing.T) {
	mock := new(MockS3Client)
	mock.PutBucketTaggingAPI = DefaultPutBucketTaggingFunc

	res, err := mock.PutBucketTagging(context.Background(), &s3.PutBucketTaggingInput{})
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestMockS3Client_DeleteBucketTagging(t *testing.T) {
	mock := new(MockS3Client)
	mock.DeleteBucketTaggingAPI = DefaultDeleteBucketTaggingFunc

	res, err := mock.DeleteBucketTagging(context.Background(), &s3.DeleteBucketTaggingInput{})
	assert.NotNil(t, res)
	assert.Nil(t, err)
}
