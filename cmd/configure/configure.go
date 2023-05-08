package configure

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/configure/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/logging"
	"github.com/rs/zerolog"

	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"

	"github.com/spf13/cobra"
)

func init() {
	configureOpts = options.GetConfigureOptions()
	configureOpts.InitFlags(ConfigureCmd)
}

var (
	logger        zerolog.Logger
	configureOpts *options.ConfigureOptions
	svc           *s3.S3
	ConfigureCmd  = &cobra.Command{
		Use:          "configure",
		Short:        "configure subcommand configures the bucket level settings",
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
			configureOpts.RootOptions = rootOpts
			logger = logging.GetLogger(rootOpts)

			svc, err = aws.CreateAwsService(rootOpts)
			if err != nil {
				logger.Error().
					Str("error", err.Error()).
					Msg("an error occurred while creating aws service")
				return err
			}

			logger.Info().Msg("aws service successfully created with provided AWS credentials")

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			versioning, err := aws.GetBucketVersioning(svc, configureOpts)
			if err != nil {
				return err
			}

			if *versioning.Status == "Enabled" && configureOpts.Versioning || *versioning.Status == "Suspended" && !configureOpts.Versioning {
				logger.Info().Msg("versioning is already at the desired state, skipping")
				return nil
			}

			logger.Info().Msgf("setting versioning as %v", configureOpts.Versioning)
			_, err = aws.SetBucketVersioning(svc, configureOpts, configureOpts.Versioning)
			if err != nil {
				return err
			}

			versioning, err = aws.GetBucketVersioning(svc, configureOpts)
			if err != nil {
				return err
			}
			logger.Info().Msg(*versioning.Status)

			return nil
		},
	}
)
