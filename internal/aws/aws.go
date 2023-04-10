package aws

import (
	"bytes"
	"log"
	"strings"
	"sync"

	options2 "github.com/bilalcaliskan/s3-manager/cmd/find/options"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/rs/zerolog"
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
func GetAllFiles(svc s3iface.S3API, opts *options.RootOptions, prefix string) (res *s3.ListObjectsOutput, err error) {
	// fetch all the objects in target bucket
	res, err = svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(opts.BucketName),
		Prefix: aws.String(prefix),
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

// Find does the heavy lifting, communicates with the S3 and finds the files
func Find(svc s3iface.S3API, opts *options2.FindOptions, logger zerolog.Logger) ([]string, []error) {
	var errors []error
	var matchedFiles []string
	mu := &sync.Mutex{}

	// fetch all the objects in target bucket
	listResult, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(opts.BucketName),
	})
	if err != nil {
		errors = append(errors, err)
		return matchedFiles, errors
	}

	var resultArr []*s3.Object
	var wg sync.WaitGroup

	extensions := strings.Split(opts.FileExtensions, ",")

	// separate the txt files from all of the fetched objects from bucket
	for _, v := range listResult.Contents {
		for _, y := range extensions {
			if strings.HasSuffix(*v.Key, y) {
				logger.Debug().Str("fileName", *v.Key).Msg("found file")
				resultArr = append(resultArr, v)
			}
		}
	}

	// check each txt file individually if it contains provided substring
	for _, obj := range resultArr {
		wg.Add(1)
		go func(obj *s3.Object, wg *sync.WaitGroup) {
			defer wg.Done()
			getResult, err := svc.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(opts.BucketName),
				Key:    obj.Key,
			})

			if err != nil {
				errors = append(errors, err)
				return
			}

			buf := new(bytes.Buffer)
			if _, err := buf.ReadFrom(getResult.Body); err != nil {
				errors = append(errors, err)
				return
			}

			if strings.Contains(buf.String(), opts.Substring) {
				mu.Lock()
				matchedFiles = append(matchedFiles, *obj.Key)
				mu.Unlock()
			}

			defer func() {
				if err := getResult.Body.Close(); err != nil {
					panic(err)
				}
			}()

		}(obj, &wg)
	}

	// wait for all the goroutines to complete
	wg.Wait()

	return matchedFiles, errors
}
