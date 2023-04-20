package clean

import (
	"fmt"

	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/logging"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/clean/options"
	"github.com/bilalcaliskan/s3-manager/internal/cleaner"
	"github.com/bilalcaliskan/s3-manager/internal/utils"

	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	cleanOpts = options.GetCleanOptions()
	cleanOpts.InitFlags(CleanCmd)
}

var (
	logger          zerolog.Logger
	ValidSortByOpts = []string{"size", "lastModificationDate"}
	cleanOpts       *options.CleanOptions
	svc             *s3.S3
	// CleanCmd represents the clean command
	CleanCmd = &cobra.Command{
		Use:          "clean",
		Short:        "clean subcommand cleans the app, finds and clears desired files",
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
			cleanOpts.RootOptions = rootOpts
			logger = logging.GetLogger(rootOpts)

			svc, err = aws.CreateAwsService(rootOpts)
			if err != nil {
				logger.Error().
					Str("error", err.Error()).
					Msg("an error occurred while creating aws service")
				return err
			}

			logger.Info().Msg("aws service successfully created with provided AWS credentials")

			if cleanOpts.MinFileSizeInMb > cleanOpts.MaxFileSizeInMb && (cleanOpts.MinFileSizeInMb != 0 && cleanOpts.MaxFileSizeInMb != 0) {
				err := fmt.Errorf("minFileSizeInMb should be lower than maxFileSizeInMb")
				logger.Error().Str("error", err.Error()).Msg("an error occured while validating flags")
				return err
			}

			if !utils.Contains(ValidSortByOpts, cleanOpts.SortBy) {
				err := fmt.Errorf("no such sortBy option called %s, valid options are %v", cleanOpts.SortBy,
					ValidSortByOpts)
				logger.Error().Str("error", err.Error()).Msg("an error occurred while validating flags")
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info().Msg("trying to search files on target bucket")

			return cleaner.StartCleaning(svc, cleanOpts, logger)
		},
	}
)
