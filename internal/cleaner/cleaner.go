package cleaner

import (
	"bytes"
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
	allFiles, err := aws.GetAllFiles(svc, cleanOpts.RootOptions, cleanOpts.FileNamePrefix)
	if err != nil {
		return err
	}

	res := getProperObjects(cleanOpts, allFiles, logger)
	sortObjects(res, cleanOpts)

	border := len(res) - cleanOpts.KeepLastNFiles
	if border <= 0 {
		logger.Warn().
			Int("arrayLength", len(res)).
			Int("keepLastNFiles", cleanOpts.KeepLastNFiles).
			Msg("not enough file, length of array is smaller than --keepLastNFiles flag")
		return nil
	}

	targetObjects := res[:len(res)-cleanOpts.KeepLastNFiles]

	keys := utils.GetKeysOnly(targetObjects)
	var buffer bytes.Buffer
	for _, v := range keys {
		buffer.WriteString(v)
	}

	logger.Info().Any("files", keys).Msg("will attempt to delete these files")
	if cleanOpts.DryRun {
		logger.Info().Msg("skipping object deletion since --dryRun flag is passed")
		return nil
	}

	if !cleanOpts.AutoApprove {
		logger.Info().Any("files", keys).Msg("these files will be removed if you approve:")

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
