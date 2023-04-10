package find

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/find/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/rs/zerolog"

	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"

	"github.com/spf13/cobra"
)

func init() {
	findOpts = options.GetFindOptions()
	findOpts.InitFlags(FindCmd)
}

var (
	logger   zerolog.Logger
	findOpts *options.FindOptions
	svc      *s3.S3
	// FindCmd represents the find subcommand
	FindCmd = &cobra.Command{
		Use:   "find",
		Short: "find subcommand finds the files which has desired substrings in it",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			logger = cmd.Context().Value(rootopts.LoggerKey{}).(zerolog.Logger)
			rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
			findOpts.RootOptions = rootOpts
			svc = cmd.Context().Value(rootopts.S3SvcKey{}).(*s3.S3)

			// TODO: add validation needed logic here

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			logger.Debug().
				Str("fileExtensions", findOpts.FileExtensions).
				Msg("trying to find files on target bucket")

			matchedFiles, errors := aws.Find(svc, findOpts, logger)
			if len(errors) != 0 {
				logger.Error().Str("error", err.Error()).Msg("an error occurred while finding target files on target bucket")
				return err
			}

			if len(matchedFiles) == 0 {
				logger.Info().
					Any("matchedFiles", matchedFiles).
					Str("substring", findOpts.Substring).
					Msg("no matched files on the bucket")
				return nil
			}

			logger.Info().
				Any("matchedFiles", matchedFiles).
				Str("substring", findOpts.Substring).
				Msg("fetched matching files")
			return nil
		},
	}
)
