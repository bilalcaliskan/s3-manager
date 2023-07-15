package enabled

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/cmd/versioning/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/prompt"
	"github.com/bilalcaliskan/s3-manager/internal/utils"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	versioningOpts = options.GetVersioningOptions()
}

var (
	svc            s3iface.S3API
	logger         zerolog.Logger
	confirmRunner  prompt.PromptRunner = prompt.GetConfirmRunner()
	versioningOpts *options.VersioningOptions
	EnabledCmd     = &cobra.Command{
		Use:           "enabled",
		Short:         "enables the versioning configuration for the target bucket",
		SilenceUsage:  false,
		SilenceErrors: true,
		Example: `# set the versioning configuration for bucket as enabled
s3-manager versioning set enabled
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var rootOpts *rootopts.RootOptions
			svc, rootOpts, logger, confirmRunner = utils.PrepareConstants(cmd)
			versioningOpts.RootOptions = rootOpts

			if err := utils.CheckArgs(args, 0); err != nil {
				logger.Error().
					Msg(err.Error())
				return err
			}

			versioningOpts.DesiredState = "enabled"

			return aws.SetBucketVersioning(svc, versioningOpts, confirmRunner, logger)
		},
	}
)
