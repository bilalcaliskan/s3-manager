package cleaner

import (
	"fmt"
	"strings"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	start "github.com/bilalcaliskan/s3-manager/cmd/clean/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/utils"
	"github.com/rs/zerolog"
)

func StartCleaning(svc s3iface.S3API, runner prompt.PromptRunner, cleanOpts *start.CleanOptions, logger zerolog.Logger) error {
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
