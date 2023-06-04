package aws

import (
	"bytes"
	"errors"
	"log"
	"regexp"
	"strings"
	"sync"

	options3 "github.com/bilalcaliskan/s3-manager/cmd/tags/options"

	options2 "github.com/bilalcaliskan/s3-manager/cmd/search/options"

	options4 "github.com/bilalcaliskan/s3-manager/cmd/versioning/options"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/rs/zerolog"
)

// createSession initializes session with provided credentials
func createSession(accessKey, secretKey, region string) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})

	return sess, err
}

func CreateAwsService(opts *options.RootOptions) (svc *s3.S3, err error) {
	var sess *session.Session

	if opts.AccessKey == "" || opts.SecretKey == "" || opts.Region == "" {
		return svc, errors.New("missing required fields")
	}

	sess, err = createSession(opts.AccessKey, opts.SecretKey, opts.Region)
	if err != nil {
		return svc, err
	}

	return s3.New(sess), err
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

func GetBucketTags(svc s3iface.S3API, opts *options3.TagOptions) (res *s3.GetBucketTaggingOutput, err error) {
	// fetch all the objects in target bucket
	res, err = svc.GetBucketTagging(&s3.GetBucketTaggingInput{
		Bucket: aws.String(opts.BucketName),
	})

	if err != nil {
		return res, err
	}

	return res, nil
}

func SetBucketTags(svc s3iface.S3API, opts *options3.TagOptions) (res *s3.PutBucketTaggingOutput, err error) {
	// fetch all the objects in target bucket
	var tagsSet []*s3.Tag
	for i, v := range opts.TagsToAdd {
		tag := &s3.Tag{
			Key:   aws.String(i),
			Value: aws.String(v),
		}
		tagsSet = append(tagsSet, tag)
	}

	res, err = svc.PutBucketTagging(&s3.PutBucketTaggingInput{
		Bucket:  aws.String(opts.BucketName),
		Tagging: &s3.Tagging{TagSet: tagsSet},
	})

	if err != nil {
		return res, err
	}

	return res, nil
}

func DeleteAllBucketTags(svc s3iface.S3API, opts *options3.TagOptions) (res *s3.DeleteBucketTaggingOutput, err error) {
	return svc.DeleteBucketTagging(&s3.DeleteBucketTaggingInput{
		Bucket: aws.String(opts.BucketName),
	})
}

// GetBucketVersioning gets the target bucket
func GetBucketVersioning(svc s3iface.S3API, opts *options.RootOptions) (res *s3.GetBucketVersioningOutput, err error) {
	// fetch all the objects in target bucket
	res, err = svc.GetBucketVersioning(&s3.GetBucketVersioningInput{
		Bucket: aws.String(opts.BucketName),
	})

	if err != nil {
		return res, err
	}

	return res, nil
}

// SetBucketVersioning sets the target bucket
func SetBucketVersioning(svc s3iface.S3API, opts *options4.VersioningOptions) (res *s3.PutBucketVersioningOutput, err error) {
	var str string
	if opts.DesiredState == "enabled" {
		str = "Enabled"
	} else if opts.DesiredState == "disabled" {
		str = "Suspended"
	}

	res, err = svc.PutBucketVersioning(&s3.PutBucketVersioningInput{
		Bucket: aws.String(opts.BucketName),
		VersioningConfiguration: &s3.VersioningConfiguration{
			Status: aws.String(str),
		},
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

func GetDesiredFiles(svc s3iface.S3API, opts *options2.SearchOptions) (matchedFiles []string, err error) {
	// fetch all the objects in target bucket
	listResult, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(opts.BucketName),
	})
	if err != nil {
		return matchedFiles, err
	}

	pattern := regexp.MustCompile(opts.FileName)
	for _, v := range listResult.Contents {
		if match := pattern.FindStringSubmatch(*v.Key); len(match) > 0 {
			matchedFiles = append(matchedFiles, *v.Key)
		}
	}

	return matchedFiles, err
}

// SearchString does the heavy lifting, communicates with the S3 and finds the files
func SearchString(svc s3iface.S3API, opts *options2.SearchOptions, logger zerolog.Logger) ([]string, []error) {
	var errs []error
	var matchedFiles []string
	mu := &sync.Mutex{}

	// fetch all the objects in target bucket
	listResult, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(opts.BucketName),
	})
	if err != nil {
		errs = append(errs, err)
		return matchedFiles, errs
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
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}

			buf := new(bytes.Buffer)
			if _, err := buf.ReadFrom(getResult.Body); err != nil {
				errs = append(errs, err)
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

	return matchedFiles, errs
}
