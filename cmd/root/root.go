package root

import (
	"context"
	"os"
	"strings"

	"github.com/bilalcaliskan/s3-manager/cmd/clean"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"

	"github.com/dimiro1/banner"

	"github.com/bilalcaliskan/s3-manager/internal/logging"
	"github.com/bilalcaliskan/s3-manager/internal/version"
	"github.com/spf13/cobra"
)

func init() {
	opts = options.GetRootOptions()
	opts.InitFlags(rootCmd)

	rootCmd.AddCommand(clean.CleanCmd)
}

var (
	opts           *options.RootOptions
	ver            = version.Get()
	bannerFilePath = "build/ci/banner.txt"
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:     "s3-manager",
		Short:   "",
		Long:    ``,
		Version: ver.GitVersion,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat("build/ci/banner.txt"); err == nil {
				bannerBytes, _ := os.ReadFile(bannerFilePath)
				banner.Init(os.Stdout, true, false, strings.NewReader(string(bannerBytes)))
			}

			if opts.VerboseLog {
				logging.EnableDebugLogging()
			}

			logger := logging.GetLogger(opts)
			logger.Info().Str("appVersion", ver.GitVersion).Str("goVersion", ver.GoVersion).Str("goOS", ver.GoOs).
				Str("goArch", ver.GoArch).Str("gitCommit", ver.GitCommit).Str("buildDate", ver.BuildDate).
				Msg("s3-manager is started!")

			cmd.SetContext(context.WithValue(cmd.Context(), options.LoggerKey{}, logger))
			cmd.SetContext(context.WithValue(cmd.Context(), options.OptsKey{}, opts))

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
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
