package file

import (
	"fmt"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/aws/types"

	"github.com/bilalcaliskan/s3-manager/internal/aws"

	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"

	"github.com/bilalcaliskan/s3-manager/internal/utils"

	"github.com/bilalcaliskan/s3-manager/cmd/search/options"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	searchOpts = options.GetSearchOptions()
}

var (
	logger     zerolog.Logger
	searchOpts *options.SearchOptions
	svc        internalawstypes.S3ClientAPI
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
			svc, rootOpts, logger, _ = utils.PrepareConstants(cmd)
			searchOpts.RootOptions = rootOpts

			if err := utils.CheckArgs(args, 1); err != nil {
				logger.Error().Msg(err.Error())
				return err
			}

			searchOpts.FileName = args[0]

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// TODO: recover from panic if something is broken with regex
			files, err := aws.GetDesiredObjects(svc, searchOpts.BucketName, searchOpts.FileName)
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
				fmt.Println(*v.Key)
			}

			return nil
		},
	}
)
