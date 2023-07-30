package aws

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	internalutil "github.com/bilalcaliskan/s3-manager/internal/utils"
	"github.com/pkg/errors"

	options6 "github.com/bilalcaliskan/s3-manager/cmd/transferacceleration/options"

	options5 "github.com/bilalcaliskan/s3-manager/cmd/bucketpolicy/options"

	"github.com/bilalcaliskan/s3-manager/cmd/versioning/set/utils"

	options3 "github.com/bilalcaliskan/s3-manager/cmd/tags/options"

	options2 "github.com/bilalcaliskan/s3-manager/cmd/search/options"

	options4 "github.com/bilalcaliskan/s3-manager/cmd/versioning/options"

	v2aws "github.com/aws/aws-sdk-go-v2/aws"
	v2config "github.com/aws/aws-sdk-go-v2/config"
	v2creds "github.com/aws/aws-sdk-go-v2/credentials"
	v2s3 "github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/rs/zerolog"
)

// createSession initializes a new AWS SDK session using provided credentials.
//
// It accepts access key, secret key and region as arguments, then uses them
// to create an AWS Config, which is then used to initialize the session.
// It returns a pointer to session.Session along with any error encountered during
// session initialization process. If no error occurred during the process, the error is nil.
//func createSession(accessKey, secretKey, region string) (*session.Session, error) {
//	return session.NewSession(&aws.Config{
//		Region:      aws.String(region),
//		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
//	})
//}

func CreateClient(opts *options.RootOptions) (*v2s3.Client, error) {
	appCreds := v2aws.NewCredentialsCache(v2creds.NewStaticCredentialsProvider(opts.AccessKey, opts.SecretKey, ""))
	//value, err := appCreds.Retrieve(context.TODO())
	//if err != nil {
	//	// handle error
	//}

	config, err := v2config.LoadDefaultConfig(context.Background(),
		v2config.WithRegion(opts.Region),
		v2config.WithCredentialsProvider(appCreds),
	)

	if err != nil {
		return nil, err
	}

	return v2s3.NewFromConfig(config), nil
}

// GetBucketTags retrieves all tags attached to a specific S3 bucket.
//
// It accepts an S3API interface and pointer of TagOptions as arguments, and returns
// a GetBucketTaggingOutput, which contains all the bucket's tags, and any error encountered.
func GetBucketTags(svc S3ClientAPI, opts *options3.TagOptions) (res *v2s3.GetBucketTaggingOutput, err error) {
	return svc.GetBucketTagging(context.Background(), &v2s3.GetBucketTaggingInput{
		Bucket: aws.String(opts.BucketName),
	})
}

// SetBucketTags attaches a set of tags to a specific S3 bucket.
//
// It accepts an S3API interface and TagOptions as arguments.
// For each tag in the provided TagOptions, a new tag is created and added to a slice of tags.
// It then attaches these tags to the bucket and returns a PutBucketTaggingOutput and any error encountered.
func SetBucketTags(svc S3ClientAPI, opts *options3.TagOptions, runner prompt.PromptRunner, logger zerolog.Logger) error {
	if opts.DryRun {
		logger.Info().Msg(constants.InfDryRun)
		return nil
	}

	if !opts.AutoApprove {
		if err := prompt.AskForApproval(runner); err != nil {
			return err
		}
	}

	// fetch all the objects in target bucket
	var tagsSet []types.Tag
	for i, v := range opts.TagsToAdd {
		tag := types.Tag{
			Key:   aws.String(i),
			Value: aws.String(v),
		}
		tagsSet = append(tagsSet, tag)
	}

	_, err := svc.PutBucketTagging(context.Background(), &v2s3.PutBucketTaggingInput{
		Bucket:  aws.String(opts.BucketName),
		Tagging: &types.Tagging{TagSet: tagsSet},
	})

	if err != nil {
		return err
	}

	return nil
}

// DeleteAllBucketTags removes all tags attached to a specific S3 bucket.
//
// It accepts an S3API interface and TagOptions as arguments, and returns
// a DeleteBucketTaggingOutput and any error encountered.
func DeleteAllBucketTags(svc S3ClientAPI, opts *options3.TagOptions, runner prompt.PromptRunner, logger zerolog.Logger) (out *v2s3.DeleteBucketTaggingOutput, err error) {
	if opts.DryRun {
		logger.Info().Msg(constants.InfDryRun)
		return out, nil
	}

	if !opts.AutoApprove {
		if err := prompt.AskForApproval(runner); err != nil {
			return out, err
		}
	}

	return svc.DeleteBucketTagging(context.Background(), &v2s3.DeleteBucketTaggingInput{
		Bucket: aws.String(opts.BucketName),
	})
}

// GetTransferAcceleration retrieves the current transfer acceleration status of an S3 bucket.
//
// It accepts an S3API interface and TransferAccelerationOptions as arguments,
// and returns a GetBucketAccelerateConfigurationOutput and any error encountered.
func GetTransferAcceleration(svc S3ClientAPI, opts *options6.TransferAccelerationOptions) (res *v2s3.GetBucketAccelerateConfigurationOutput, err error) {
	return svc.GetBucketAccelerateConfiguration(context.Background(), &v2s3.GetBucketAccelerateConfigurationInput{
		Bucket: aws.String(opts.BucketName),
	})
}

// SetTransferAcceleration sets the transfer acceleration status of an S3 bucket.
//
// It accepts an S3API interface, TransferAccelerationOptions, a PromptRunner, and a Logger as arguments.
// If the provided 'DryRun' or 'AutoApprove' options are set, the function will return early.
// If not, it will set the bucket's transfer acceleration status based on the provided desired state.
// It logs any errors encountered and returns them.
func SetTransferAcceleration(svc S3ClientAPI, opts *options6.TransferAccelerationOptions, runner prompt.PromptRunner, logger zerolog.Logger) error {
	if opts.DryRun {
		logger.Info().Msg(constants.InfDryRun)
		return nil
	}

	if !opts.AutoApprove {
		if err := prompt.AskForApproval(runner); err != nil {
			return err
		}
	}

	res, err := GetTransferAcceleration(svc, opts)
	if err != nil {
		logger.Error().Msg(err.Error())
		return err
	}

	if res.Status == "Enabled" {
		opts.ActualState = "enabled"
	} else if res.Status == "Suspended" {
		opts.ActualState = "disabled"
	} else {
		err := fmt.Errorf("unknown status '%s' returned from AWS SDK", opts.ActualState)
		logger.Error().Msg(err.Error())
		return err
	}

	if opts.DesiredState == opts.ActualState {
		logger.Warn().Msg("transfer acceleration configuration is already at desired state")
		return nil
	}

	var status string
	switch opts.DesiredState {
	case "enabled":
		status = "Enabled"
	case "disabled":
		status = "Suspended"
	}

	_, err = svc.PutBucketAccelerateConfiguration(context.Background(), &v2s3.PutBucketAccelerateConfigurationInput{
		Bucket:                  aws.String(opts.BucketName),
		AccelerateConfiguration: &types.AccelerateConfiguration{Status: types.BucketAccelerateStatus(status)},
	})

	if err != nil {
		logger.Error().Msg(err.Error())
		return err
	}

	logger.Info().Msgf("successfully set transfer acceleration as %s", opts.DesiredState)

	return nil
}

// GetBucketPolicy retrieves the current policy of an S3 bucket.
//
// It accepts an S3API interface and BucketPolicyOptions as arguments,
// and returns a GetBucketPolicyOutput and any error encountered.
func GetBucketPolicy(svc S3ClientAPI, opts *options5.BucketPolicyOptions) (res *v2s3.GetBucketPolicyOutput, err error) {
	return svc.GetBucketPolicy(context.Background(), &v2s3.GetBucketPolicyInput{
		Bucket: aws.String(opts.BucketName),
	})

	//return svc.GetBucketPolicy()
}

// GetBucketPolicyString retrieves the current policy of an S3 bucket and beautifies it into a readable format.
//
// It accepts an S3API interface and BucketPolicyOptions as arguments,
// and returns the beautified policy as a string and any error encountered.
func GetBucketPolicyString(svc S3ClientAPI, opts *options5.BucketPolicyOptions) (out string, err error) {
	res, err := GetBucketPolicy(svc, opts)
	if err != nil {
		return out, errors.Wrap(err, "an error occurred while getting bucket policy")
	}

	return internalutil.BeautifyJSON(*res.Policy)
}

// SetBucketPolicy sets the policy of an S3 bucket.
//
// It accepts an S3API interface and BucketPolicyOptions as arguments,
// and returns a PutBucketPolicyOutput and any error encountered.
func SetBucketPolicy(svc S3ClientAPI, opts *options5.BucketPolicyOptions, runner prompt.PromptRunner, logger zerolog.Logger) (res *v2s3.PutBucketPolicyOutput, err error) {
	if opts.DryRun {
		logger.Info().Msg(constants.InfDryRun)
		return res, nil
	}

	if !opts.AutoApprove {
		if err := prompt.AskForApproval(runner); err != nil {
			return res, err
		}
	}

	return svc.PutBucketPolicy(context.Background(), &v2s3.PutBucketPolicyInput{
		Bucket: aws.String(opts.BucketName),
		Policy: aws.String(opts.BucketPolicyContent),
	})
}

// DeleteBucketPolicy removes the existing policy from a specified S3 bucket.
//
// It requires an S3API interface and a BucketPolicyOptions object, which should
// include the name of the target bucket. The function will use these to generate
// and execute a DeleteBucketPolicyInput request via the provided S3 service.
// It returns a DeleteBucketPolicyOutput, which acknowledges the operation,
// along with any error encountered during the process.
func DeleteBucketPolicy(svc S3ClientAPI, opts *options5.BucketPolicyOptions, runner prompt.PromptRunner, logger zerolog.Logger) (res *v2s3.DeleteBucketPolicyOutput, err error) {
	if opts.DryRun {
		logger.Info().Msg(constants.InfDryRun)
		return res, nil
	}

	if !opts.AutoApprove {
		if err := prompt.AskForApproval(runner); err != nil {
			return res, err
		}
	}

	return svc.DeleteBucketPolicy(context.Background(), &v2s3.DeleteBucketPolicyInput{
		Bucket: aws.String(opts.BucketName),
	})
}

// GetBucketVersioning retrieves the versioning state of the specified S3 bucket.
//
// The function accepts an S3API interface and RootOptions, which should include
// the name of the target bucket. It uses these to generate and execute a
// GetBucketVersioningInput request via the provided S3 service.
// The function returns a GetBucketVersioningOutput, which includes the bucket's
// versioning configuration, along with any error encountered during the process.
func GetBucketVersioning(svc S3ClientAPI, opts *options.RootOptions) (res *v2s3.GetBucketVersioningOutput, err error) {
	return svc.GetBucketVersioning(context.Background(), &v2s3.GetBucketVersioningInput{
		Bucket: aws.String(opts.BucketName),
	})
}

// SetBucketVersioning updates the versioning configuration of a specific S3 bucket.
//
// It accepts an S3API interface, VersioningOptions, a PromptRunner for user confirmations,
// and a Logger for logging events. The function uses these to check for dry-run or auto-approve
// flags, confirm versioning state changes with the user if needed, and execute a
// PutBucketVersioningInput request to set the bucket's versioning state.
// The function logs the process, including any errors encountered, and returns these errors.
func SetBucketVersioning(svc S3ClientAPI, versioningOpts *options4.VersioningOptions, runner prompt.PromptRunner, logger zerolog.Logger) (err error) {
	if versioningOpts.DryRun {
		logger.Info().Msg(constants.InfDryRun)
		return nil
	}

	if !versioningOpts.AutoApprove {
		if err := prompt.AskForApproval(runner); err != nil {
			return err
		}
	}

	var versioning *v2s3.GetBucketVersioningOutput
	versioning, err = GetBucketVersioning(svc, versioningOpts.RootOptions)
	if err != nil {
		logger.Error().Msg(err.Error())
		return err
	}

	if err = utils.DecideActualState(versioning, versioningOpts); err != nil {
		logger.Error().Msg(err.Error())
		return err
	}

	logger.Info().Msgf(utils.InfCurrentState, versioningOpts.ActualState)
	if versioningOpts.ActualState == versioningOpts.DesiredState {
		logger.Warn().
			Str("state", versioningOpts.ActualState).
			Msg(utils.WarnDesiredState)
		return nil
	}

	logger.Info().Msgf(utils.InfSettingVersioning, versioningOpts.DesiredState)

	var str string
	switch versioningOpts.DesiredState {
	case "enabled":
		str = "Enabled"
	case "disabled":
		str = "Suspended"
	}

	if _, err = svc.PutBucketVersioning(context.Background(), &v2s3.PutBucketVersioningInput{
		Bucket: aws.String(versioningOpts.BucketName),
		VersioningConfiguration: &types.VersioningConfiguration{
			Status: types.BucketVersioningStatus(str),
		},
	}); err != nil {
		logger.Error().Msg(err.Error())
		return err
	}

	logger.Info().Msgf(utils.InfSuccess, versioningOpts.DesiredState)
	return nil
}

// DeleteFiles removes a specific list of objects from a specified S3 bucket.
//
// The function accepts an S3API interface, the name of the target bucket, an array of
// S3 Objects to delete, a dryRun boolean flag, and a Logger. It iterates over the array of
// objects, logging each one, and unless dryRun is set, it sends a DeleteObjectInput request
// for each object to the S3 service. The function logs each successful deletion and returns
// any errors encountered during the process.
func DeleteFiles(svc S3ClientAPI, bucketName string, slice []*types.Object, dryRun bool, logger zerolog.Logger) error {
	for _, v := range slice {
		logger.Debug().Str("key", *v.Key).Time("lastModifiedDate", *v.LastModified).
			Float64("size", float64(v.Size)/1000000).Msg("will try to delete file")

		if !dryRun {
			if _, err := svc.DeleteObject(context.Background(), &v2s3.DeleteObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(*v.Key),
			}); err != nil {
				return err
			}

			logger.Info().Str("key", *v.Key).Msg("successfully deleted file")
		}
	}

	return nil
}

// GetDesiredObjects retrieves a list of objects in a specified S3 bucket that match a given
//
// regular expression. The function takes an S3API interface, the target bucket's name, and
// the regex as arguments. It fetches all objects in the target bucket and filters them using
// the regex. The function returns a list of matching S3 Object pointers and any error encountered.
func GetDesiredObjects(svc S3ClientAPI, bucketName, regex string) (objects []*types.Object, err error) {
	// fetch all the objects in target bucket
	listResult, err := svc.ListObjects(context.Background(), &v2s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return objects, err
	}

	pattern := regexp.MustCompile(regex)
	for _, v := range listResult.Contents {
		if match := pattern.FindStringSubmatch(*v.Key); len(match) > 0 {
			objects = append(objects, &v)
		}
	}

	return objects, err
}

// SearchString scans all objects in a specified S3 bucket for a given text string.
//
// The function accepts an S3API interface and SearchOptions, which include the bucket
// name, file name pattern, and search text. It first retrieves a list of objects that match
// the file name pattern, then concurrently checks each object's content for the search text.
// The function returns a list of object keys that contain the search text and a list of errors
// encountered during the search process.
func SearchString(svc S3ClientAPI, opts *options2.SearchOptions) (matchedFiles []string, errs []error) {
	var wg sync.WaitGroup
	mu := &sync.Mutex{}

	resultArr, err := GetDesiredObjects(svc, opts.BucketName, opts.FileName)
	if err != nil {
		errs = append(errs, err)
		return matchedFiles, errs
	}

	// check each txt file individually if it contains provided text
	for _, obj := range resultArr {
		wg.Add(1)
		go func(obj *types.Object, wg *sync.WaitGroup) {
			defer wg.Done()
			getResult, err := svc.GetObject(context.Background(), &v2s3.GetObjectInput{
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
			mu.Lock()
			if _, err := buf.ReadFrom(getResult.Body); err != nil {
				errs = append(errs, err)
				return
			}
			mu.Unlock()

			if strings.Contains(buf.String(), opts.Text) {
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

	wg.Wait()
	return matchedFiles, errs
}
