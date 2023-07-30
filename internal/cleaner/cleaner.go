package cleaner

import (
	"fmt"
	"strings"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	start "github.com/bilalcaliskan/s3-manager/cmd/clean/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/utils"
	"github.com/rs/zerolog"
)

// StartCleaning is a function that performs a deletion operation on AWS S3 files based on specified parameters.
//
// The function requires an S3 service, a prompt runner, clean options, and a logger as parameters.
// The function first retrieves the list of desired objects (files) from the specified AWS S3 bucket that match the
// provided regular expression. If an error occurs during retrieval, it immediately returns the error.
// The retrieved objects are sorted according to the configuration specified in the CleanOptions.
//
// The function then calculates the border index in the sorted array from which deletion should start, which is
// determined by subtracting the number of files to keep from the total number of retrieved objects. If the border
// index is less than or equal to 0, it means there aren't enough files to delete; it logs a warning message and
// the function returns without deleting any files.
//
// Next, it prepares a list of target objects (files) to delete based on the border index calculated previously.
// The file names (keys) of these target objects are extracted and logged for information.
//
// If the --dryRun flag is set to true in the CleanOptions, the function skips the actual deletion process,
// logs a message, and returns. This is a way to simulate a cleaning operation without making actual deletions.
//
// If the --autoApprove flag is set to false, the function prompts the user for approval before proceeding
// with deletion. If the user input is 'n' (denoting No), the function returns an error indicating that the
// user terminated the process. If the user input is neither 'y' nor 'n', it returns an error indicating invalid user input.
//
// Finally, the function performs the deletion of the target files from the S3 bucket. If an error occurs during
// deletion, it logs the error and returns it.
//
// The function returns nil if it completes without encountering any errors.
func StartCleaning(svc aws.S3ClientAPI, runner prompt.PromptRunner, cleanOpts *start.CleanOptions, logger zerolog.Logger) error {
	res, err := aws.GetDesiredObjects(svc, cleanOpts.BucketName, cleanOpts.Regex)
	if err != nil {
		return err
	}

	sortObjects(res, cleanOpts)

	border := len(res) - cleanOpts.KeepLastNFiles
	if border <= 0 {
		logger.Warn().
			Int("arrayLength", len(res)).
			Msg("not enough file, length of array is smaller than --keepLastNFiles flag")
		return nil
	}

	targetObjects := res[:len(res)-cleanOpts.KeepLastNFiles]

	keys := utils.GetKeysOnly(targetObjects)

	logger.Info().Msg("will attempt to delete these files")
	for _, key := range keys {
		fmt.Println(key)
	}

	if cleanOpts.DryRun {
		logger.Info().Msg("skipping object deletion since --dryRun flag is passed")
		return nil
	}

	if !cleanOpts.AutoApprove {
		logger.Info().Msg("above files will be removed if you approve")

		if res, err := runner.Run(); err != nil {
			if strings.ToLower(res) == "n" {
				return constants.ErrUserTerminated
			}

			return constants.ErrInvalidInput
		}
	}

	if err := aws.DeleteFiles(svc, cleanOpts.RootOptions.BucketName, targetObjects, cleanOpts.DryRun, logger); err != nil {
		logger.Error().Str("error", err.Error()).Msg("an error occurred while deleting target files")
		return err
	}

	return nil
}
