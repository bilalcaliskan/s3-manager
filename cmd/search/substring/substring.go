package substring

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	options2 "github.com/bilalcaliskan/s3-manager/cmd/search/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/logging"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	searchOpts = options2.GetSearchOptions()
	searchOpts.InitFlags(SubstringCmd)
}

var (
	logger       zerolog.Logger
	searchOpts   *options2.SearchOptions
	svc          s3iface.S3API
	SubstringCmd = &cobra.Command{
		Use:          "substring",
		Short:        "",
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			svc = cmd.Context().Value(rootopts.S3SvcKey{}).(s3iface.S3API)
			rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
			searchOpts.RootOptions = rootOpts

			logger = logging.GetLogger(searchOpts.RootOptions)

			if err := checkFlags(logger, args); err != nil {
				return err
			}

			searchOpts.Substring = args[0]

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// TODO: uncomment when interactivity enabled again
			/*if searchOpts.Interactive {
				if err := searchOpts.PromptInteractiveValues(); err != nil {
					logger.Error().Str("error", err.Error()).Msg("an error occurred while prompting values")
					return err
				}
			}*/

			logger.Info().
				Str("fileExtensions", searchOpts.FileExtensions).
				Msg("trying to search files on target bucket")

			matchedFiles, errs := aws.SearchString(svc, searchOpts, logger)
			if len(errs) != 0 {
				err := fmt.Errorf("multiple errors occurred while searching files, try to target individual files %s", errs)
				logger.Error().Str("error", err.Error())
				return err
			}

			if len(matchedFiles) == 0 {
				logger.Info().
					Any("matchedFiles", matchedFiles).
					Str("substring", searchOpts.Substring).
					Msg("no matched files on the bucket")
				return nil
			}

			logger.Info().
				Any("matchedFiles", matchedFiles).
				Str("substring", searchOpts.Substring).
				Msg("fetched matching files")
			return nil
		},
	}
)
