package types

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	DefaultGetBucketTaggingFunc = func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
		return &s3.GetBucketTaggingOutput{}, nil
	}
	DefaultPutBucketTaggingFunc = func(ctx context.Context, params *s3.PutBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.PutBucketTaggingOutput, error) {
		return &s3.PutBucketTaggingOutput{}, nil
	}
	DefaultDeleteBucketTaggingFunc = func(ctx context.Context, params *s3.DeleteBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketTaggingOutput, error) {
		return &s3.DeleteBucketTaggingOutput{}, nil
	}
	DefaultPutBucketAccelerationFunc = func(ctx context.Context, params *s3.PutBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.PutBucketAccelerateConfigurationOutput, error) {
		return &s3.PutBucketAccelerateConfigurationOutput{}, nil
	}
	DefaultGetBucketVersioningFunc = func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
		return &s3.GetBucketVersioningOutput{}, nil
	}
	DefaultPutBucketVersioningFunc = func(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
		return &s3.PutBucketVersioningOutput{}, nil
	}
)

type S3ClientAPI interface {
	GetBucketPolicy(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error)
	PutBucketPolicy(ctx context.Context, params *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error)
	DeleteBucketPolicy(ctx context.Context, params *s3.DeleteBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketPolicyOutput, error)

	GetBucketAccelerateConfiguration(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error)
	PutBucketAccelerateConfiguration(ctx context.Context, params *s3.PutBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.PutBucketAccelerateConfigurationOutput, error)

	GetBucketVersioning(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error)
	PutBucketVersioning(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error)

	GetBucketTagging(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error)
	PutBucketTagging(ctx context.Context, params *s3.PutBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.PutBucketTaggingOutput, error)
	DeleteBucketTagging(ctx context.Context, params *s3.DeleteBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketTaggingOutput, error)

	ListObjects(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

type MockS3Client struct {
	GetBucketPolicyAPI                  func(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error)
	GetBucketAccelerateConfigurationAPI func(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error)
	PutBucketAccelerateConfigurationAPI func(ctx context.Context, params *s3.PutBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.PutBucketAccelerateConfigurationOutput, error)
	GetBucketVersioningAPI              func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error)
	PutBucketVersioningAPI              func(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error)
	GetBucketTaggingAPI                 func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error)
	PutBucketTaggingAPI                 func(ctx context.Context, params *s3.PutBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.PutBucketTaggingOutput, error)
	DeleteBucketTaggingAPI              func(ctx context.Context, params *s3.DeleteBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketTaggingOutput, error)
	ListObjectsAPI                      func(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error)
	GetObjectAPI                        func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	DeleteObjectAPI                     func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	PutBucketPolicyAPI                  func(ctx context.Context, params *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error)
	DeleteBucketPolicyAPI               func(ctx context.Context, params *s3.DeleteBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketPolicyOutput, error)
}

func (m *MockS3Client) PutBucketPolicy(ctx context.Context, params *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error) {
	return m.PutBucketPolicyAPI(ctx, params, optFns...)
}

func (m *MockS3Client) DeleteBucketPolicy(ctx context.Context, params *s3.DeleteBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketPolicyOutput, error) {
	return m.DeleteBucketPolicyAPI(ctx, params, optFns...)
}

func (m *MockS3Client) DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	return m.DeleteObjectAPI(ctx, params, optFns...)
}

func (m *MockS3Client) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return m.GetObjectAPI(ctx, params, optFns...)
}

func (m *MockS3Client) ListObjects(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
	return m.ListObjectsAPI(ctx, params, optFns...)
}

func (m *MockS3Client) GetBucketPolicy(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
	return m.GetBucketPolicyAPI(ctx, params, optFns...)
}

func (m *MockS3Client) GetBucketAccelerateConfiguration(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error) {
	return m.GetBucketAccelerateConfigurationAPI(ctx, params, optFns...)
}

func (m *MockS3Client) PutBucketAccelerateConfiguration(ctx context.Context, params *s3.PutBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.PutBucketAccelerateConfigurationOutput, error) {
	return m.PutBucketAccelerateConfigurationAPI(ctx, params, optFns...)
}

func (m *MockS3Client) GetBucketVersioning(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
	return m.GetBucketVersioningAPI(ctx, params, optFns...)
}

func (m *MockS3Client) PutBucketVersioning(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
	return m.PutBucketVersioningAPI(ctx, params, optFns...)
}

func (m *MockS3Client) GetBucketTagging(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
	return m.GetBucketTaggingAPI(ctx, params, optFns...)
}

func (m *MockS3Client) PutBucketTagging(ctx context.Context, params *s3.PutBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.PutBucketTaggingOutput, error) {
	return m.PutBucketTaggingAPI(ctx, params, optFns...)
}

func (m *MockS3Client) DeleteBucketTagging(ctx context.Context, params *s3.DeleteBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketTaggingOutput, error) {
	return m.DeleteBucketTaggingAPI(ctx, params, optFns...)
}
