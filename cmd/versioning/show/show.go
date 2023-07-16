package show

import (
	"fmt"

	"github.com/bilalcaliskan/s3-manager/internal/utils"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/cmd/versioning/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
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
	ShowCmd        = &cobra.Command{
		Use:           "show",
		Short:         "shows the versioning configuration for the target bucket",
		SilenceUsage:  false,
		SilenceErrors: true,
		Example: `# show the current versioning configuration for bucket
s3-manager versioning show
		`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var rootOpts *rootopts.RootOptions
			svc, rootOpts, logger, _ = utils.PrepareConstants(cmd)
			versioningOpts.RootOptions = rootOpts

			if err := utils.CheckArgs(args, 0); err != nil {
				logger.Error().Msg(err.Error())
				return err
			}

			versioning, err := aws.GetBucketVersioning(svc, versioningOpts.RootOptions)
			if err != nil {
				return err
			}

			switch *versioning.Status {
			case "Enabled":
				versioningOpts.ActualState = "enabled"
			case "Suspended":
				versioningOpts.ActualState = "disabled"
			default:
				err := fmt.Errorf("unknown versioning status %s returned from S3 SDK", *versioning.Status)
				logger.Error().Msg(err.Error())
				return err
			}

			logger.Info().Msgf("current versioning configuration is %s", versioningOpts.ActualState)

			return nil
		},
	}
)
