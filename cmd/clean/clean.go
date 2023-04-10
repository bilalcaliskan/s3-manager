package clean

import (
	"fmt"

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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			logger = cmd.Context().Value(rootopts.LoggerKey{}).(zerolog.Logger)
			rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
			cleanOpts.RootOptions = rootOpts
			svc = cmd.Context().Value(rootopts.S3SvcKey{}).(*s3.S3)

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
			logger.Info().Msg("trying to find files on target bucket")

			return cleaner.StartCleaning(svc, cleanOpts, logger)
		},
	}
)
