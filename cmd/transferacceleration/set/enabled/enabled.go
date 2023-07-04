package enabled

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/transferacceleration/options"
	"github.com/bilalcaliskan/s3-manager/cmd/transferacceleration/utils"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/prompt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	transferAccelerationOpts = options.GetTransferAccelerationOptions()
}

var (
	svc                      s3iface.S3API
	logger                   zerolog.Logger
	confirmRunner            prompt.PromptRunner = prompt.GetConfirmRunner()
	transferAccelerationOpts *options.TransferAccelerationOptions
	EnabledCmd               = &cobra.Command{
		Use:           "enabled",
		Short:         "enables the transfer acceleration configuration for the target bucket",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: `# set the transfer acceleration configuration for bucket as enabled
s3-manager transferacceleration set enabled
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, transferAccelerationOpts, logger = utils.PrepareConstants(cmd, transferAccelerationOpts)

			if err := utils.CheckArgs(args); err != nil {
				logger.Error().
					Msg(err.Error())
				return err
			}

			transferAccelerationOpts.DesiredState = "enabled"

			if transferAccelerationOpts.DryRun {
				logger.Info().Msg("skipping operation since '--dry-run' flag is passed")
				return nil
			}

			var err error
			if !transferAccelerationOpts.AutoApprove {
				var res string
				if res, err = confirmRunner.Run(); err != nil {
					return err
				}

				if strings.ToLower(res) == "n" {
					return errors.New("user terminated the process")
				}
			}

			return aws.SetTransferAcceleration(svc, transferAccelerationOpts, logger)
		},
	}
)
