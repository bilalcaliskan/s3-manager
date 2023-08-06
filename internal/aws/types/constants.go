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
