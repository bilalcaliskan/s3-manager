package cleaner

import (
	"bytes"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	start "github.com/bilalcaliskan/s3-manager/cmd/clean/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/utils"
	"github.com/rs/zerolog"
)

func StartCleaning(svc s3iface.S3API, startOpts *start.CleanOptions, logger zerolog.Logger) error {
	allFiles, err := aws.GetAllFiles(svc, startOpts.RootOptions)
	if err != nil {
		return err
	}

	res := getProperObjects(startOpts, allFiles, logger)
	sortObjects(res, startOpts)

	border := len(res) - startOpts.KeepLastNFiles
	if border < 0 {
		logger.Warn().
			Int("arrayLength", len(res)).
			Int("keepLastNFiles", startOpts.KeepLastNFiles).
			Msg("not enough file, length of array is smaller than --keepLastNFiles flag")
		return nil
	}

	targetObjects := res[:len(res)-startOpts.KeepLastNFiles]
	if err := checkLength(targetObjects); err != nil {
		logger.Warn().Msg(err.Error())
		return nil
	}

	keys := utils.GetKeysOnly(targetObjects)
	var buffer bytes.Buffer
	for _, v := range keys {
		buffer.WriteString(v)
	}

	logger.Info().Any("files", keys).Msg("will attempt to delete these files")
	if startOpts.DryRun {
		logger.Info().Msg("skipping object deletion since --dryRun flag is passed")
		return nil
	}

	if err := promptDeletion(startOpts, logger, keys); err != nil {
		logger.Warn().Str("error", err.Error()).Msg("an error occurred while prompting file deletion")
		return err
	}

	if err := aws.DeleteFiles(svc, startOpts.RootOptions.BucketName, targetObjects, startOpts.DryRun, logger); err != nil {
		logger.Error().Str("error", err.Error()).Msg("an error occurred while deleting target files")
		return err
	}

	return nil
}
