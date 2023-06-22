package enabled

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/versioning/options"
	"github.com/bilalcaliskan/s3-manager/cmd/versioning/set/utils"
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
	EnabledCmd     = &cobra.Command{
		Use:           "enabled",
		Short:         "enables the versioning configuration for the target bucket",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: `# set the versioning configuration for bucket as enabled
s3-manager versioning set enabled
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, versioningOpts, logger = utils.PrepareConstants(cmd, options.GetVersioningOptions())

			if err := utils.CheckArgs(args); err != nil {
				logger.Error().
					Msg(err.Error())
				return err
			}

			versioningOpts.DesiredState = "enabled"

			return aws.SetBucketVersioning(svc, versioningOpts, logger)
		},
	}
)
