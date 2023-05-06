package root

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"
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
}

/*
- persistentprerun, prerun dan daha once calisiyor
- required olan bir flagi gecmesen bile persistentprerun ve prerun calisiyor (root command icin)
*/

var (
	selectRunner    prompt.SelectRunner = prompt.GetSelectRunner("Select operation", []string{"search", "clean"})
	accessKeyRunner prompt.PromptRunner = prompt.GetPromptRunner("Provide AWS Access Key", nil)
	secretKeyRunner prompt.PromptRunner = prompt.GetPromptRunner("Provide AWS Secret Key", nil)
	regionRunner    prompt.PromptRunner = prompt.GetPromptRunner("Provide AWS Region", nil)
	bucketRunner    prompt.PromptRunner = prompt.GetPromptRunner("Provide AWS Bucket Name", nil)
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
			//if !opts.Interactive {
			//	opts.SetAccessFlagsRequired(cmd)
			//}

			// https://sonarcloud.io/component_measures?id=bilalcaliskan_s3-manager&metric=coverage&view=list
			// TODO: create svc here instead of each subcommand
			// TODO: fail if credentials are expired (meaning wrong credentials provided)

			//if _, err := os.Stat(opts.BannerFilePath); err == nil {
			//	bannerBytes, _ := os.ReadFile(opts.BannerFilePath)
			//	banner.Init(os.Stdout, true, false, strings.NewReader(string(bannerBytes)))
			//}
			//
			//if opts.VerboseLog {
			//	logging.EnableDebugLogging()
			//}
			//
			//logger = logging.GetLogger(opts)
			//logger.Info().Str("appVersion", ver.GitVersion).Str("goVersion", ver.GoVersion).Str("goOS", ver.GoOs).
			//	Str("goArch", ver.GoArch).Str("gitCommit", ver.GitCommit).Str("buildDate", ver.BuildDate).
			//	Msg("s3-manager is started!")

			//cmd.SetContext(context.WithValue(cmd.Context(), options.LoggerKey{}, logger))
			//cmd.SetContext(context.WithValue(cmd.Context(), options.OptsKey{}, opts))
			//ctx, cancel := context.WithTimeout(cmd.Context(), 10*time.Second)
			//
			//cmd.SetContext(context.WithValue(ctx, options.OptsKey{}, cancel))
			//cmd.SetContext(ctx)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			for i := 0; i < 500; i++ {
				log.Println("hello")
				time.Sleep(1 * time.Second)
			}

			return nil

			/*if !opts.Interactive {
				return nil
			}

			if opts.AccessKey == "" {
				res, err := accessKeyRunner.Run()
				if err != nil {
					return err
				}

				opts.AccessKey = res
			}

			if opts.SecretKey == "" {
				res, err := secretKeyRunner.Run()
				if err != nil {
					return err
				}

				opts.SecretKey = res
			}

			if opts.Region == "" {
				res, err := regionRunner.Run()
				if err != nil {
					return err
				}

				opts.Region = res
			}

			if opts.BucketName == "" {
				res, err := bucketRunner.Run()
				if err != nil {
					return err
				}

				opts.BucketName = res
			}

			cmd.SetContext(context.WithValue(cmd.Context(), options.OptsKey{}, opts))
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

			return nil*/
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	fmt.Println(opts.Timeout)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(opts.Timeout)*time.Second)
	defer cancel()

	go func() {
		if err := rootCmd.Execute(); err != nil {
			os.Exit(1)
		}
	}()

	select {
	//case <-time.After(1 * time.Second):
	//	fmt.Println("overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err()) // prints "context deadline exceeded"
	}
}
