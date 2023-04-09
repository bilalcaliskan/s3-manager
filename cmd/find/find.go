package clean

import (
	"github.com/bilalcaliskan/s3-manager/cmd/find/options"
	"github.com/rs/zerolog"

	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"

	"github.com/spf13/cobra"
)

func init() {
	findOpts = options.GetFindOptions()
	findOpts.InitFlags(CleanCmd)
}

var (
	logger   zerolog.Logger
	findOpts *options.FindOptions
	// CleanCmd represents the clean command
	CleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "clean subcommand cleans the app, finds and clears desired files",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			logger = cmd.Context().Value(rootopts.LoggerKey{}).(zerolog.Logger)
			rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
			findOpts.RootOptions = rootOpts

			logger.Info().Msg("dummy log")

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
)
