package aws

import (
	"log"

	"github.com/rs/zerolog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
)

// CreateSession initializes session with provided credentials
func CreateSession(opts *options.RootOptions) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(opts.Region),
		Credentials: credentials.NewStaticCredentials(opts.AccessKey, opts.SecretKey, ""),
	})

	return sess, err
}

// GetAllFiles gets all of the files in the target bucket as the function name indicates
func GetAllFiles(svc s3iface.S3API, opts *options.RootOptions) (*s3.ListObjectsOutput, error) {
	var err error
	var res *s3.ListObjectsOutput

	// fetch all the objects in target bucket
	res, err = svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(opts.BucketName),
		Prefix: aws.String(opts.FileNamePrefix),
	})
	if err != nil {
		return res, err
	}

	return res, nil
}

// DeleteFiles deletes the slice of []*s3.Object objects in the target bucket
func DeleteFiles(svc s3iface.S3API, bucketName string, slice []*s3.Object, dryRun bool, logger zerolog.Logger) error {
	for _, v := range slice {
		logger.Debug().Str("key", *v.Key).Time("lastModifiedDate", *v.LastModified).
			Float64("size", float64(*v.Size)/1000000).Msg("will try to delete file")

		if !dryRun {
			_, err := svc.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(*v.Key),
			})

			if err != nil {
				return err
			}

			log.Printf("successfully deleted file %s", *v.Key)
			logger.Info().Str("key", *v.Key).Msg("successfully deleted file")
		}
	}

	return nil
}
