package disabled

import (
	"github.com/bilalcaliskan/s3-manager/cmd/versioning/set/utils"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
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
	DisabledCmd    = &cobra.Command{
		Use:           "disabled",
		Short:         "disables the versioning configuration for the target bucket",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			svc, versioningOpts, logger = utils.PrepareConstants(cmd, options.GetVersioningOptions())

			if err := utils.CheckArgs(args); err != nil {
				logger.Error().
					Msg(err.Error())
				return err
			}

			versioningOpts.DesiredState = "disabled"

			return aws.SetBucketVersioning(svc, versioningOpts, logger)
		},
	}
)
