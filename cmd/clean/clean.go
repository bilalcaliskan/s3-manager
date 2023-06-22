package clean

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	"github.com/bilalcaliskan/s3-manager/cmd/clean/options"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/internal/cleaner"
	"github.com/bilalcaliskan/s3-manager/internal/logging"
	"github.com/bilalcaliskan/s3-manager/internal/utils"
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
	svc             s3iface.S3API
	promptRunner    prompt.PromptRunner = prompt.GetPromptRunner("Delete Files? (y/N)", true, func(s string) error {
		if len(s) == 1 {
			return nil
		}

		return errors.New("invalid input")
	})
	// CleanCmd represents the clean command
	CleanCmd = &cobra.Command{
		Use:          "clean",
		Short:        "finds and clears desired files by a pre-configured rule set",
		SilenceUsage: true,
		// we should not define PreRunE since it will override the PreRunE which is inherited from RootCmd
		RunE: func(cmd *cobra.Command, args []string) error {
			svc = cmd.Context().Value(rootopts.S3SvcKey{}).(s3iface.S3API)
			rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
			cleanOpts.RootOptions = rootOpts
			logger = logging.GetLogger(rootOpts)

			if cleanOpts.MinFileSizeInMb > cleanOpts.MaxFileSizeInMb && (cleanOpts.MinFileSizeInMb != 0 && cleanOpts.MaxFileSizeInMb != 0) {
				err := fmt.Errorf("flag '--min-size-mb' must be equal or lower than '--max-size-mb'")
				logger.Error().Str("error", err.Error()).Msg("an error occured while validating flags")
				return err
			}

			if !utils.Contains(ValidSortByOpts, cleanOpts.SortBy) {
				err := fmt.Errorf("no such '--sort-by' option called %s, valid options are %v", cleanOpts.SortBy,
					ValidSortByOpts)
				logger.Error().Str("error", err.Error()).Msg("an error occurred while validating flags")
				return err
			}

			logger.Info().Msg("trying to search files on target bucket")

			if err := cleaner.StartCleaning(svc, promptRunner, cleanOpts, logger); err != nil {
				logger.Error().Str("error", err.Error()).Msg("an error occurred while cleaning")
				return err
			}

			return nil
		},
	}
)
