package search

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/bilalcaliskan/s3-manager/cmd/search/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/logging"
	"github.com/rs/zerolog"

	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"

	"github.com/spf13/cobra"
)

func init() {
	searchOpts = options.GetSearchOptions()
	searchOpts.InitFlags(SearchCmd)
}

var (
	logger     zerolog.Logger
	searchOpts *options.SearchOptions
	svc        s3iface.S3API
	SearchCmd  = &cobra.Command{
		Use:          "search",
		Short:        "search subcommand searches the files which has desired substrings in it",
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			svc = cmd.Context().Value(rootopts.S3SvcKey{}).(s3iface.S3API)
			rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
			searchOpts.RootOptions = rootOpts

			logger = logging.GetLogger(rootOpts)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if searchOpts.Interactive {
				if err := searchOpts.PromptInteractiveValues(); err != nil {
					logger.Error().Str("error", err.Error()).Msg("an error occurred while prompting values")
					return err
				}
			}

			if searchOpts.Substring == "" {
				logger.Warn().Msg("will list all files in specified file extensions since substring flag is empty")
			}

			logger.Info().
				Str("fileExtensions", searchOpts.FileExtensions).
				Msg("trying to search files on target bucket")

			matchedFiles, errs := aws.Find(svc, searchOpts, logger)
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
