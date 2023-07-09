package text

import (
	"fmt"

	"github.com/bilalcaliskan/s3-manager/internal/utils"

	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"

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
		Use:           "text",
		Short:         "searches the texts in files which has desired file name pattern and string pattern in it (supports regex)",
		SilenceUsage:  false,
		SilenceErrors: true,
		Example: `# search a text on target bucket by specifying regex for files
s3-manager search text "catch me if you can" --file-name=".*.txt"
		`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var rootOpts *rootopts.RootOptions
			svc, rootOpts, logger = utils.PrepareConstants(cmd)
			searchOpts.RootOptions = rootOpts

			if err := utils.CheckArgs(args, 1); err != nil {
				logger.Error().Msg(err.Error())
				return err
			}

			searchOpts.Text = args[0]

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
				Str("fileName", searchOpts.FileName).
				Msg("trying to search files on target bucket")

			matchedFiles, errs := aws.SearchString(svc, searchOpts)
			if len(errs) != 0 {
				err := fmt.Errorf("multiple errors occurred while searching files, try to target individual files %s", errs)
				logger.Error().Str("error", err.Error())
				return err
			}

			if len(matchedFiles) == 0 {
				logger.Info().
					Any("matchedFiles", matchedFiles).
					Str("text", searchOpts.Text).
					Msg("no matched files on the bucket")
				return nil
			}

			logger.Info().
				Str("text", searchOpts.Text).
				Msg("fetched below matching files")

			for _, v := range matchedFiles {
				fmt.Println(v)
			}

			return nil
		},
	}
)
