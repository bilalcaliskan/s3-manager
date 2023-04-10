package root

import (
	"context"
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/find"
	"github.com/bilalcaliskan/s3-manager/internal/aws"

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

	if err := opts.SetAccessCredentialsFromEnv(rootCmd); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(clean.CleanCmd)
	rootCmd.AddCommand(find.FindCmd)
}

var (
	opts           *options.RootOptions
	ver            = version.Get()
	logger         zerolog.Logger
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

			logger = logging.GetLogger(opts)
			logger.Info().Str("appVersion", ver.GitVersion).Str("goVersion", ver.GoVersion).Str("goOS", ver.GoOs).
				Str("goArch", ver.GoArch).Str("gitCommit", ver.GitCommit).Str("buildDate", ver.BuildDate).
				Msg("s3-manager is started!")

			sess, err := aws.CreateSession(opts)
			if err != nil {
				logger.Error().
					Str("error", err.Error()).
					Msg("an error occurred while creating session")
				return err
			}

			logger.Info().Msg("session successfully created with provided AWS credentials")

			cmd.SetContext(context.WithValue(cmd.Context(), options.LoggerKey{}, logger))
			cmd.SetContext(context.WithValue(cmd.Context(), options.OptsKey{}, opts))
			cmd.SetContext(context.WithValue(cmd.Context(), options.S3SvcKey{}, s3.New(sess)))

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Interactive {
				prompt := promptui.Select{
					Label: "Select operation",
					Items: []string{"find", "clean"},
				}

				_, result, err := prompt.Run()

				if err != nil {
					logger.Error().Str("error", err.Error()).Msg("unknown error occurred while prompting user")
					return err
				}

				switch result {
				case "find":
					if err := find.FindCmd.PreRunE(cmd, args); err != nil {
						panic(err)
					}

					if err := find.FindCmd.RunE(cmd, args); err != nil {
						panic(err)
					}
				case "clean":
					if err := clean.CleanCmd.RunE(cmd, args); err != nil {
						panic(err)
					}
				}

				/*if res, err := prompt.Run(); err != nil {
					if strings.ToLower(res) == "n" {
						return errors.New("user terminated the process")
					}

					return errors.New("invalid input")
				}*/
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
