package disabled

import (
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/cmd/versioning/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/aws/types"
	"github.com/bilalcaliskan/s3-manager/internal/prompt"
	"github.com/bilalcaliskan/s3-manager/internal/utils"
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
	DisabledCmd    = &cobra.Command{
		Use:           "disabled",
		Short:         "disables the versioning configuration for the target bucket",
		SilenceUsage:  false,
		SilenceErrors: true,
		Example: `# set the versioning configuration for bucket as disabled
s3-manager versioning set disabled
		`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var rootOpts *rootopts.RootOptions
			svc, rootOpts, logger, confirmRunner = utils.PrepareConstants(cmd)
			versioningOpts.RootOptions = rootOpts

			if err := utils.CheckArgs(args, 0); err != nil {
				logger.Error().
					Msg(err.Error())
				return err
			}

			versioningOpts.DesiredState = "disabled"

			return aws.SetBucketVersioning(svc, versioningOpts, confirmRunner, logger)
		},
	}
)
