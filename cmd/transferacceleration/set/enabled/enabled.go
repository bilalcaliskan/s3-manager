package enabled

import (
	v2s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/cmd/transferacceleration/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/prompt"
	"github.com/bilalcaliskan/s3-manager/internal/utils"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	transferAccelerationOpts = options.GetTransferAccelerationOptions()
}

var (
	svc                      *v2s3.Client
	logger                   zerolog.Logger
	confirmRunner            prompt.PromptRunner
	transferAccelerationOpts *options.TransferAccelerationOptions
	EnabledCmd               = &cobra.Command{
		Use:           "enabled",
		Short:         "enables the transfer acceleration configuration for the target bucket",
		SilenceUsage:  false,
		SilenceErrors: true,
		Example: `# set the transfer acceleration configuration for bucket as enabled
s3-manager transferacceleration set enabled
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var rootOpts *rootopts.RootOptions
			svc, rootOpts, logger, confirmRunner = utils.PrepareConstants(cmd)
			transferAccelerationOpts.RootOptions = rootOpts

			if err := utils.CheckArgs(args, 0); err != nil {
				logger.Error().
					Msg(err.Error())
				return err
			}

			transferAccelerationOpts.DesiredState = "enabled"

			return aws.SetTransferAcceleration(svc, transferAccelerationOpts, confirmRunner, logger)
		},
	}
)
