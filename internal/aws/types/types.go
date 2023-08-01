package types

import (
	"context"

	v2s3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3ClientAPI interface {
	GetBucketPolicy(ctx context.Context, params *v2s3.GetBucketPolicyInput, optFns ...func(*v2s3.Options)) (*v2s3.GetBucketPolicyOutput, error)
	GetBucketAccelerateConfiguration(ctx context.Context, params *v2s3.GetBucketAccelerateConfigurationInput, optFns ...func(*v2s3.Options)) (*v2s3.GetBucketAccelerateConfigurationOutput, error)
	PutBucketAccelerateConfiguration(ctx context.Context, params *v2s3.PutBucketAccelerateConfigurationInput, optFns ...func(*v2s3.Options)) (*v2s3.PutBucketAccelerateConfigurationOutput, error)
	GetBucketVersioning(ctx context.Context, params *v2s3.GetBucketVersioningInput, optFns ...func(*v2s3.Options)) (*v2s3.GetBucketVersioningOutput, error)
	PutBucketVersioning(ctx context.Context, params *v2s3.PutBucketVersioningInput, optFns ...func(*v2s3.Options)) (*v2s3.PutBucketVersioningOutput, error)
	GetBucketTagging(ctx context.Context, params *v2s3.GetBucketTaggingInput, optFns ...func(*v2s3.Options)) (*v2s3.GetBucketTaggingOutput, error)
	PutBucketTagging(ctx context.Context, params *v2s3.PutBucketTaggingInput, optFns ...func(*v2s3.Options)) (*v2s3.PutBucketTaggingOutput, error)
	DeleteBucketTagging(ctx context.Context, params *v2s3.DeleteBucketTaggingInput, optFns ...func(*v2s3.Options)) (*v2s3.DeleteBucketTaggingOutput, error)
	ListObjects(ctx context.Context, params *v2s3.ListObjectsInput, optFns ...func(*v2s3.Options)) (*v2s3.ListObjectsOutput, error)
	GetObject(ctx context.Context, params *v2s3.GetObjectInput, optFns ...func(*v2s3.Options)) (*v2s3.GetObjectOutput, error)
	DeleteObject(ctx context.Context, params *v2s3.DeleteObjectInput, optFns ...func(*v2s3.Options)) (*v2s3.DeleteObjectOutput, error)
	PutBucketPolicy(ctx context.Context, params *v2s3.PutBucketPolicyInput, optFns ...func(*v2s3.Options)) (*v2s3.PutBucketPolicyOutput, error)
	DeleteBucketPolicy(ctx context.Context, params *v2s3.DeleteBucketPolicyInput, optFns ...func(*v2s3.Options)) (*v2s3.DeleteBucketPolicyOutput, error)
}

type MockS3v2Client struct {
	GetBucketPolicyAPI                  func(ctx context.Context, params *v2s3.GetBucketPolicyInput, optFns ...func(*v2s3.Options)) (*v2s3.GetBucketPolicyOutput, error)
	GetBucketAccelerateConfigurationAPI func(ctx context.Context, params *v2s3.GetBucketAccelerateConfigurationInput, optFns ...func(*v2s3.Options)) (*v2s3.GetBucketAccelerateConfigurationOutput, error)
	PutBucketAccelerateConfigurationAPI func(ctx context.Context, params *v2s3.PutBucketAccelerateConfigurationInput, optFns ...func(*v2s3.Options)) (*v2s3.PutBucketAccelerateConfigurationOutput, error)
	GetBucketVersioningAPI              func(ctx context.Context, params *v2s3.GetBucketVersioningInput, optFns ...func(*v2s3.Options)) (*v2s3.GetBucketVersioningOutput, error)
	PutBucketVersioningAPI              func(ctx context.Context, params *v2s3.PutBucketVersioningInput, optFns ...func(*v2s3.Options)) (*v2s3.PutBucketVersioningOutput, error)
	GetBucketTaggingAPI                 func(ctx context.Context, params *v2s3.GetBucketTaggingInput, optFns ...func(*v2s3.Options)) (*v2s3.GetBucketTaggingOutput, error)
	PutBucketTaggingAPI                 func(ctx context.Context, params *v2s3.PutBucketTaggingInput, optFns ...func(*v2s3.Options)) (*v2s3.PutBucketTaggingOutput, error)
	DeleteBucketTaggingAPI              func(ctx context.Context, params *v2s3.DeleteBucketTaggingInput, optFns ...func(*v2s3.Options)) (*v2s3.DeleteBucketTaggingOutput, error)
	ListObjectsAPI                      func(ctx context.Context, params *v2s3.ListObjectsInput, optFns ...func(*v2s3.Options)) (*v2s3.ListObjectsOutput, error)
	GetObjectAPI                        func(ctx context.Context, params *v2s3.GetObjectInput, optFns ...func(*v2s3.Options)) (*v2s3.GetObjectOutput, error)
	DeleteObjectAPI                     func(ctx context.Context, params *v2s3.DeleteObjectInput, optFns ...func(*v2s3.Options)) (*v2s3.DeleteObjectOutput, error)
	PutBucketPolicyAPI                  func(ctx context.Context, params *v2s3.PutBucketPolicyInput, optFns ...func(*v2s3.Options)) (*v2s3.PutBucketPolicyOutput, error)
	DeleteBucketPolicyAPI               func(ctx context.Context, params *v2s3.DeleteBucketPolicyInput, optFns ...func(*v2s3.Options)) (*v2s3.DeleteBucketPolicyOutput, error)
}

func (m *MockS3v2Client) PutBucketPolicy(ctx context.Context, params *v2s3.PutBucketPolicyInput, optFns ...func(*v2s3.Options)) (*v2s3.PutBucketPolicyOutput, error) {
	return m.PutBucketPolicyAPI(ctx, params, optFns...)
}

func (m *MockS3v2Client) DeleteBucketPolicy(ctx context.Context, params *v2s3.DeleteBucketPolicyInput, optFns ...func(*v2s3.Options)) (*v2s3.DeleteBucketPolicyOutput, error) {
	return m.DeleteBucketPolicyAPI(ctx, params, optFns...)
}

func (m *MockS3v2Client) DeleteObject(ctx context.Context, params *v2s3.DeleteObjectInput, optFns ...func(*v2s3.Options)) (*v2s3.DeleteObjectOutput, error) {
	return m.DeleteObjectAPI(ctx, params, optFns...)
}

func (m *MockS3v2Client) GetObject(ctx context.Context, params *v2s3.GetObjectInput, optFns ...func(*v2s3.Options)) (*v2s3.GetObjectOutput, error) {
	return m.GetObjectAPI(ctx, params, optFns...)
}

func (m *MockS3v2Client) ListObjects(ctx context.Context, params *v2s3.ListObjectsInput, optFns ...func(*v2s3.Options)) (*v2s3.ListObjectsOutput, error) {
	return m.ListObjectsAPI(ctx, params, optFns...)
}

func (m *MockS3v2Client) GetBucketPolicy(ctx context.Context, params *v2s3.GetBucketPolicyInput, optFns ...func(*v2s3.Options)) (*v2s3.GetBucketPolicyOutput, error) {
	return m.GetBucketPolicyAPI(ctx, params, optFns...)
}

func (m *MockS3v2Client) GetBucketAccelerateConfiguration(ctx context.Context, params *v2s3.GetBucketAccelerateConfigurationInput, optFns ...func(*v2s3.Options)) (*v2s3.GetBucketAccelerateConfigurationOutput, error) {
	return m.GetBucketAccelerateConfigurationAPI(ctx, params, optFns...)
}

func (m *MockS3v2Client) PutBucketAccelerateConfiguration(ctx context.Context, params *v2s3.PutBucketAccelerateConfigurationInput, optFns ...func(*v2s3.Options)) (*v2s3.PutBucketAccelerateConfigurationOutput, error) {
	return m.PutBucketAccelerateConfigurationAPI(ctx, params, optFns...)
}

func (m *MockS3v2Client) GetBucketVersioning(ctx context.Context, params *v2s3.GetBucketVersioningInput, optFns ...func(*v2s3.Options)) (*v2s3.GetBucketVersioningOutput, error) {
	return m.GetBucketVersioningAPI(ctx, params, optFns...)
}

func (m *MockS3v2Client) PutBucketVersioning(ctx context.Context, params *v2s3.PutBucketVersioningInput, optFns ...func(*v2s3.Options)) (*v2s3.PutBucketVersioningOutput, error) {
	return m.PutBucketVersioningAPI(ctx, params, optFns...)
}

func (m *MockS3v2Client) GetBucketTagging(ctx context.Context, params *v2s3.GetBucketTaggingInput, optFns ...func(*v2s3.Options)) (*v2s3.GetBucketTaggingOutput, error) {
	return m.GetBucketTaggingAPI(ctx, params, optFns...)
}

func (m *MockS3v2Client) PutBucketTagging(ctx context.Context, params *v2s3.PutBucketTaggingInput, optFns ...func(*v2s3.Options)) (*v2s3.PutBucketTaggingOutput, error) {
	return m.PutBucketTaggingAPI(ctx, params, optFns...)
}

func (m *MockS3v2Client) DeleteBucketTagging(ctx context.Context, params *v2s3.DeleteBucketTaggingInput, optFns ...func(*v2s3.Options)) (*v2s3.DeleteBucketTaggingOutput, error) {
	return m.DeleteBucketTaggingAPI(ctx, params, optFns...)
}

//type MockS3Client struct {
//	mock.Mock
//	s3iface.S3API
//}
//
//// GetBucketPolicy mocks the GetBucketPolicy method of s3iface.S3API
//func (m *MockS3Client) GetBucketPolicy(input *s3.GetBucketPolicyInput) (*s3.GetBucketPolicyOutput, error) {
//	// Return the mocked output values using the `On` method of testify/mock
//	args := m.Called(input)
//	return args.Get(0).(*s3.GetBucketPolicyOutput), args.Error(1)
//}
//
//func (m *MockS3Client) DeleteBucketPolicy(input *s3.DeleteBucketPolicyInput) (*s3.DeleteBucketPolicyOutput, error) {
//	// Return the mocked output values using the `On` method of testify/mock
//	args := m.Called(input)
//	return args.Get(0).(*s3.DeleteBucketPolicyOutput), args.Error(1)
//}
//
//func (m *MockS3Client) PutBucketPolicy(input *s3.PutBucketPolicyInput) (*s3.PutBucketPolicyOutput, error) {
//	// Return the mocked output values using the `On` method of testify/mock
//	args := m.Called(input)
//	return args.Get(0).(*s3.PutBucketPolicyOutput), args.Error(1)
//}
//
//func (m *MockS3Client) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
//	// Return the mocked output values using the `On` method of testify/mock
//	args := m.Called(input)
//	return args.Get(0).(*s3.DeleteObjectOutput), args.Error(1)
//}
//
//func (m *MockS3Client) ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
//	// Return the mocked output values using the `On` method of testify/mock
//	args := m.Called(input)
//	return args.Get(0).(*s3.ListObjectsOutput), args.Error(1)
//}
//
//// GetObject mocks the GetObject method of s3iface.S3API
//func (m *MockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
//	bytes, err := os.Open(*input.Key)
//	if err != nil {
//		return nil, err
//	}
//
//	args := m.Called(input)
//
//	return &s3.GetObjectOutput{
//		AcceptRanges:  aws.String("bytes"),
//		Body:          bytes,
//		ContentLength: aws.Int64(1000),
//		ContentType:   aws.String("text/plain"),
//		ETag:          aws.String("d73a503d212d9279e6b2ed8ac6bb81f3"),
//	}, args.Error(1)
//}
//
//// GetBucketTagging mocks the GetBucketTagging method of s3iface.S3API
//func (m *MockS3Client) GetBucketTagging(input *s3.GetBucketTaggingInput) (*s3.GetBucketTaggingOutput, error) {
//	// Return the mocked output values using the `On` method of testify/mock
//	args := m.Called(input)
//	return args.Get(0).(*s3.GetBucketTaggingOutput), args.Error(1)
//}
//
//// PutBucketTagging mocks the PutBucketTagging method of s3iface.S3API
//func (m *MockS3Client) PutBucketTagging(input *s3.PutBucketTaggingInput) (*s3.PutBucketTaggingOutput, error) {
//	// Return the mocked output values using the `On` method of testify/mock
//	args := m.Called(input)
//	return args.Get(0).(*s3.PutBucketTaggingOutput), args.Error(1)
//}
//
//// DeleteBucketTagging mocks the DeleteBucketTagging method of s3iface.S3API
//func (m *MockS3Client) DeleteBucketTagging(input *s3.DeleteBucketTaggingInput) (*s3.DeleteBucketTaggingOutput, error) {
//	// Return the mocked output values using the `On` method of testify/mock
//	args := m.Called(input)
//	return args.Get(0).(*s3.DeleteBucketTaggingOutput), args.Error(1)
//}
//
//// GetBucketAccelerateConfiguration mocks the GetBucketAccelerateConfiguration method of s3iface.S3API
//func (m *MockS3Client) GetBucketAccelerateConfiguration(input *s3.GetBucketAccelerateConfigurationInput) (*s3.GetBucketAccelerateConfigurationOutput, error) {
//	// Return the mocked output values using the `On` method of testify/mock
//	args := m.Called(input)
//	return args.Get(0).(*s3.GetBucketAccelerateConfigurationOutput), args.Error(1)
//}
//
//// PutBucketAccelerateConfiguration mocks the PutBucketAccelerateConfiguration method of s3iface.S3API
//func (m *MockS3Client) PutBucketAccelerateConfiguration(input *s3.PutBucketAccelerateConfigurationInput) (*s3.PutBucketAccelerateConfigurationOutput, error) {
//	// Return the mocked output values using the `On` method of testify/mock
//	args := m.Called(input)
//	return args.Get(0).(*s3.PutBucketAccelerateConfigurationOutput), args.Error(1)
//}
//
//// GetBucketVersioning mocks the GetBucketVersioning method of s3iface.S3API
//func (m *MockS3Client) GetBucketVersioning(input *s3.GetBucketVersioningInput) (*s3.GetBucketVersioningOutput, error) {
//	// Return the mocked output values using the `On` method of testify/mock
//	args := m.Called(input)
//	return args.Get(0).(*s3.GetBucketVersioningOutput), args.Error(1)
//}
//
//// PutBucketVersioning mocks the PutBucketVersioning method of s3iface.S3API
//func (m *MockS3Client) PutBucketVersioning(input *s3.PutBucketVersioningInput) (*s3.PutBucketVersioningOutput, error) {
//	// Return the mocked output values using the `On` method of testify/mock
//	args := m.Called(input)
//	return args.Get(0).(*s3.PutBucketVersioningOutput), args.Error(1)
//}
