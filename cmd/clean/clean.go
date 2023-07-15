package clean

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	"github.com/bilalcaliskan/s3-manager/cmd/clean/options"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/internal/cleaner"
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
	confirmRunner   prompt.PromptRunner
	CleanCmd        = &cobra.Command{
		Use:           "clean",
		Short:         "finds and clears desired files by a pre-configured rule set",
		SilenceUsage:  false,
		SilenceErrors: true,
		Example: `# clean the desired files on target bucket
s3-manager clean --min-size-mb=1 --max-size-mb=1000 --keep-last-n-files=2 --sort-by=lastModificationDate
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var rootOpts *rootopts.RootOptions
			svc, rootOpts, logger, confirmRunner = utils.PrepareConstants(cmd)
			cleanOpts.RootOptions = rootOpts

			if err := utils.CheckArgs(args, 0); err != nil {
				logger.Error().Msg(err.Error())
				return err
			}

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

			if err := cleaner.StartCleaning(svc, confirmRunner, cleanOpts, logger); err != nil {
				logger.Error().Str("error", err.Error()).Msg("an error occurred while cleaning")
				return err
			}

			return nil
		},
	}
)
