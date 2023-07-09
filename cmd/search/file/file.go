package file

import (
	"fmt"

	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"

	"github.com/bilalcaliskan/s3-manager/internal/utils"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/search/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	searchOpts = options.GetSearchOptions()
}

var (
	logger     zerolog.Logger
	searchOpts *options.SearchOptions
	svc        s3iface.S3API
	FileCmd    = &cobra.Command{
		Use:           "file",
		Short:         "searches the files which has desired file name pattern in it (supports regex)",
		SilenceUsage:  false,
		SilenceErrors: true,
		Example: `# search a file on target bucket by specifying regex for files
s3-manager search file ".*.json"
		`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var rootOpts *rootopts.RootOptions
			svc, rootOpts, logger = utils.PrepareConstants(cmd)
			searchOpts.RootOptions = rootOpts

			if err := utils.CheckArgs(args, 1); err != nil {
				logger.Error().Msg(err.Error())
				return err
			}

			searchOpts.FileName = args[0]

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

			files, err := aws.GetDesiredFiles(svc, searchOpts)
			if err != nil {
				logger.Error().
					Str("fileName", searchOpts.FileName).
					Str("error", err.Error()).
					Msg("an error occurred while fetching desired files")
				return err
			}

			if len(files) == 0 {
				logger.Warn().
					Str("fileName", searchOpts.FileName).
					Msg("no file found with the specified fileName or pattern")
				return nil
			}

			for _, v := range files {
				fmt.Println(v)
			}

			return nil
		},
	}
)
