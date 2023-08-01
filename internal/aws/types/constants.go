package types

import (
	"context"
	v2s3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	DefaultGetBucketTaggingFunc = func(ctx context.Context, params *v2s3.GetBucketTaggingInput, optFns ...func(*v2s3.Options)) (*v2s3.GetBucketTaggingOutput, error) {
		return &v2s3.GetBucketTaggingOutput{}, nil
	}
	DefaultPutBucketTaggingFunc = func(ctx context.Context, params *v2s3.PutBucketTaggingInput, optFns ...func(*v2s3.Options)) (*v2s3.PutBucketTaggingOutput, error) {
		return &v2s3.PutBucketTaggingOutput{}, nil
	}
	DefaultDeleteBucketTaggingFunc = func(ctx context.Context, params *v2s3.DeleteBucketTaggingInput, optFns ...func(*v2s3.Options)) (*v2s3.DeleteBucketTaggingOutput, error) {
		return &v2s3.DeleteBucketTaggingOutput{}, nil
	}
	DefaultPutBucketAccelerationFunc = func(ctx context.Context, params *v2s3.PutBucketAccelerateConfigurationInput, optFns ...func(*v2s3.Options)) (*v2s3.PutBucketAccelerateConfigurationOutput, error) {
		return &v2s3.PutBucketAccelerateConfigurationOutput{}, nil
	}
	DefaultGetBucketVersioningFunc = func(ctx context.Context, params *v2s3.GetBucketVersioningInput, optFns ...func(*v2s3.Options)) (*v2s3.GetBucketVersioningOutput, error) {
		return &v2s3.GetBucketVersioningOutput{}, nil
	}
	DefaultPutBucketVersioningFunc = func(ctx context.Context, params *v2s3.PutBucketVersioningInput, optFns ...func(*v2s3.Options)) (*v2s3.PutBucketVersioningOutput, error) {
		return &v2s3.PutBucketVersioningOutput{}, nil
	}
)
