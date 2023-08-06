package root

import (
	"context"
	"os"
	"strings"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	"github.com/bilalcaliskan/s3-manager/cmd/transferacceleration"

	"github.com/bilalcaliskan/s3-manager/cmd/bucketpolicy"

	"github.com/bilalcaliskan/s3-manager/cmd/tags"

	"github.com/bilalcaliskan/s3-manager/cmd/versioning"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/logging"
	"github.com/dimiro1/banner"
	"github.com/rs/zerolog"

	"github.com/bilalcaliskan/s3-manager/cmd/clean"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/cmd/search"

	"github.com/bilalcaliskan/s3-manager/internal/version"
	"github.com/spf13/cobra"
)

func init() {
	opts = options.GetRootOptions()
	opts.InitFlags(rootCmd)

	if err := opts.SetAccessCredentialsFromEnv(); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(clean.CleanCmd)
	rootCmd.AddCommand(search.SearchCmd)
	rootCmd.AddCommand(versioning.VersioningCmd)
	rootCmd.AddCommand(tags.TagsCmd)
	rootCmd.AddCommand(bucketpolicy.BucketPolicyCmd)
	rootCmd.AddCommand(transferacceleration.TransferAccelerationCmd)
}

var (
	opts    *options.RootOptions
	ver     = version.Get()
	logger  zerolog.Logger
	rootCmd = &cobra.Command{
		Use:           "s3-manager",
		Short:         "configure subcommand configures the bucket level settings",
		Long:          ``,
		Version:       ver.GitVersion,
		SilenceUsage:  false,
		SilenceErrors: false,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			opts.SetAccessFlagsRequired(cmd)

			client, err := aws.CreateClient(opts)
			if err != nil {
				logger.Error().
					Str("error", err.Error()).
					Msg("an error occurred while creating s3 client")
				return err
			}

			if _, err := os.Stat(opts.BannerFilePath); err == nil {
				bannerBytes, _ := os.ReadFile(opts.BannerFilePath)
				banner.Init(os.Stdout, true, false, strings.NewReader(string(bannerBytes)))
			}

			if opts.VerboseLog {
				logging.EnableDebugLogging()
			}

			logger = logging.GetLogger(opts)
			logger.Info().Str("appVersion", ver.GitVersion).Str("goVersion", ver.GoVersion).Str("goOS", ver.GoOs).
				Str("goArch", ver.GoArch).Str("gitCommit", ver.GitCommit).Str("buildDate", ver.BuildDate).
				Msg("s3-manager is started!")

			cmd.SetContext(context.WithValue(cmd.Context(), options.LoggerKey{}, logger))
			cmd.SetContext(context.WithValue(cmd.Context(), options.OptsKey{}, opts))
			cmd.SetContext(context.WithValue(cmd.Context(), options.S3ClientKey{}, client))
			cmd.SetContext(context.WithValue(cmd.Context(), options.ConfirmRunnerKey{}, prompt.GetConfirmRunner()))

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}
