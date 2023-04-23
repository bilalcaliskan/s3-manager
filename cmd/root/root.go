package root

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	"github.com/rs/zerolog"

	"github.com/bilalcaliskan/s3-manager/cmd/clean"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/cmd/search"

	"github.com/dimiro1/banner"

	"github.com/bilalcaliskan/s3-manager/internal/logging"
	"github.com/bilalcaliskan/s3-manager/internal/version"
	"github.com/spf13/cobra"
)

func init() {
	opts = options.GetRootOptions()
	opts.InitFlags(rootCmd)

	if err := opts.SetAccessCredentialsFromEnv(rootCmd); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(clean.CleanCmd)
	rootCmd.AddCommand(search.SearchCmd)
}

/*
- persistentprerun, prerun dan daha once calisiyor
- required olan bir flagi gecmesen bile persistentprerun ve prerun calisiyor (root command icin)
*/

var (
	selectRunner prompt.SelectRunner
	//promptRunner prompt.PromptRunner
	opts   *options.RootOptions
	ver    = version.Get()
	logger zerolog.Logger
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:     "s3-manager",
		Short:   "",
		Long:    ``,
		Version: ver.GitVersion,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !opts.Interactive {
				opts.SetAccessFlagsRequired(cmd)
			}

			if _, err := os.Stat(opts.BannerFilePath); err == nil {
				bannerBytes, _ := os.ReadFile(opts.BannerFilePath)
				banner.Init(os.Stdout, true, false, strings.NewReader(string(bannerBytes)))
			}

			if opts.VerboseLog {
				fmt.Println("here")
				logging.EnableDebugLogging()
			}

			logger = logging.GetLogger(opts)
			logger.Info().Str("appVersion", ver.GitVersion).Str("goVersion", ver.GoVersion).Str("goOS", ver.GoOs).
				Str("goArch", ver.GoArch).Str("gitCommit", ver.GitCommit).Str("buildDate", ver.BuildDate).
				Msg("s3-manager is started!")

			cmd.SetContext(context.WithValue(cmd.Context(), options.LoggerKey{}, logger))
			cmd.SetContext(context.WithValue(cmd.Context(), options.OptsKey{}, opts))

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Interactive {
				if err := prompt.PromptAccessCreds(opts, logger); err != nil {
					return err
				}

				cmd.SetContext(context.WithValue(cmd.Context(), options.OptsKey{}, opts))

				selectRunner = prompt.GetSelectRunner("Select operation", []string{"search", "clean"})
				_, result, err := selectRunner.Run()
				if err != nil {
					logger.Error().Str("error", err.Error()).Msg("unknown error occurred while prompting user")
					return err
				}

				switch result {
				case "search":
					if err := search.SearchCmd.PreRunE(cmd, args); err != nil {
						logger.Error().Str("error", err.Error()).Msg("an error occurred while running search subcommand")
						return err
					}

					if err := search.SearchCmd.RunE(cmd, args); err != nil {
						logger.Error().Str("error", err.Error()).Msg("an error occurred while running search subcommand")
						return err
					}
				case "clean":
					if err := clean.CleanCmd.PreRunE(cmd, args); err != nil {
						logger.Error().Str("error", err.Error()).Msg("an error occurred while running clean subcommand")
						return err
					}

					if err := clean.CleanCmd.RunE(cmd, args); err != nil {
						logger.Error().Str("error", err.Error()).Msg("an error occurred while running clean subcommand")
						return err
					}
				}
			}

			return nil
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
