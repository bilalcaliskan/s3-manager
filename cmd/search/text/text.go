package text

import (
	"fmt"

	"github.com/bilalcaliskan/s3-manager/cmd/search/utils"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/search/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	searchOpts = options.GetSearchOptions()
	searchOpts.InitFlags(TextCmd)
}

var (
	logger     zerolog.Logger
	searchOpts *options.SearchOptions
	svc        s3iface.S3API
	TextCmd    = &cobra.Command{
		Use:          "text",
		Short:        "",
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			svc, searchOpts, logger = utils.PrepareConstants(cmd, options.GetSearchOptions())

			if err := utils.CheckFlags(args); err != nil {
				logger.Error().Msg(err.Error())
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
					Str("text", searchOpts.Substring).
					Msg("no matched files on the bucket")
				return nil
			}

			logger.Info().
				Any("matchedFiles", matchedFiles).
				Str("text", searchOpts.Substring).
				Msg("fetched matching files")
			return nil
		},
	}
)