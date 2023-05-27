package file

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/cmd/search/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/logging"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	searchOpts = options.GetSearchOptions()
	//searchOpts.InitFlags(FileCmd)
}

var (
	logger     zerolog.Logger
	searchOpts *options.SearchOptions
	svc        s3iface.S3API
	FileCmd    = &cobra.Command{
		Use:          "file",
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
