package enabled

import (
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/cmd/versioning/options"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/aws"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/pkg/aws/types"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/prompt"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	versioningOpts = options.GetVersioningOptions()
}

var (
	svc            internalawstypes.S3ClientAPI
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
