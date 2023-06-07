package enabled

import (
	"github.com/bilalcaliskan/s3-manager/cmd/versioning/set/utils"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/cmd/versioning/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/logging"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	versioningOpts = options.GetVersioningOptions()
}

var (
	svc            s3iface.S3API
	logger         zerolog.Logger
	versioningOpts *options.VersioningOptions
	EnabledCmd     = &cobra.Command{
		Use:           "enabled",
		Short:         "",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
			svc = cmd.Context().Value(rootopts.S3SvcKey{}).(s3iface.S3API)

			versioningOpts.RootOptions = rootOpts
			logger = logging.GetLogger(rootOpts)

			if err := utils.CheckArgs(args); err != nil {
				logger.Error().
					Msg(err.Error())
				return err
			}

			versioningOpts.DesiredState = "enabled"
			versioning, err := aws.GetBucketVersioning(svc, versioningOpts.RootOptions)
			if err != nil {
				logger.Error().Msg(err.Error())
				return err
			}

			if err := utils.DecideActualState(versioning, versioningOpts); err != nil {
				logger.Error().Msg(err.Error())
				return err
			}

			logger.Info().Msgf(utils.InfCurrentState, versioningOpts.ActualState)
			if versioningOpts.ActualState == "enabled" {
				logger.Warn().
					Str("state", versioningOpts.ActualState).
					Msg(utils.WarnDesiredState)
				return nil
			}

			logger.Info().Msgf(utils.InfSettingVersioning, versioningOpts.DesiredState)
			_, err = aws.SetBucketVersioning(svc, versioningOpts)
			if err != nil {
				return err
			}

			logger.Info().Msgf(utils.InfSuccess, versioningOpts.DesiredState)

			return nil
		},
	}
)
