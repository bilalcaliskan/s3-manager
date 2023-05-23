package root

import (
	"context"
	"os"
	"strings"

	"github.com/bilalcaliskan/s3-manager/cmd/configure"
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
	rootCmd.AddCommand(configure.ConfigureCmd)
	rootCmd.AddCommand(versioning.VersioningCmd)
}

var (
	opts    *options.RootOptions
	ver     = version.Get()
	logger  zerolog.Logger
	rootCmd = &cobra.Command{
		Use:     "s3-manager",
		Short:   "configure subcommand configures the bucket level settings",
		Long:    ``,
		Version: ver.GitVersion,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !opts.Interactive {
				opts.SetAccessFlagsRequired(cmd)
			}

			// TODO: fail if credentials are expired (meaning wrong credentials provided)

			if opts.Interactive {
				if err := opts.PromptAccessCredentials(options.AccessKeyRunner, options.SecretKeyRunner,
					options.BucketRunner, options.RegionRunner); err != nil {
					logger.Error().
						Str("error", err.Error()).
						Msg("an error occurred while creating aws service")
					return err
				}
			}

			svc, err := aws.CreateAwsService(opts)
			if err != nil {
				logger.Error().
					Str("error", err.Error()).
					Msg("an error occurred while creating aws service")
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
			cmd.SetContext(context.WithValue(cmd.Context(), options.S3SvcKey{}, svc))

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if !opts.Interactive {
				return nil
			}

			cmd.SetContext(context.WithValue(cmd.Context(), options.OptsKey{}, opts))

			_, result, err := options.SelectRunner.Run()
			if err != nil {
				logger.Error().Str("error", err.Error()).Msg("unknown error occurred while prompting user")
				return err
			}

			switch result {
			case "search":
				if err := search.SearchCmd.RunE(cmd, args); err != nil {
					logger.Error().Str("error", err.Error()).Msg("an error occurred while running search subcommand")
					return err
				}
			case "clean":
				if err := clean.CleanCmd.RunE(cmd, args); err != nil {
					logger.Error().Str("error", err.Error()).Msg("an error occurred while running clean subcommand")
					return err
				}
			}

			return nil
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}
