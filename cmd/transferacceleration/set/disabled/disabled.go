package disabled

import (
	options2 "github.com/bilalcaliskan/s3-manager/cmd/transferacceleration/options"
	"github.com/bilalcaliskan/s3-manager/cmd/transferacceleration/utils"
	"github.com/bilalcaliskan/s3-manager/internal/aws"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	transferAccelerationOpts = options2.GetTransferAccelerationOptions()
}

var (
	svc                      s3iface.S3API
	logger                   zerolog.Logger
	transferAccelerationOpts *options2.TransferAccelerationOptions
	DisabledCmd              = &cobra.Command{
		Use:           "disabled",
		Short:         "disables the transfer acceleration configuration for the target bucket",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: `# set the transfer acceleration configuration for bucket as disabled
s3-manager transferacceleration set disabled
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, transferAccelerationOpts, logger = utils.PrepareConstants(cmd, transferAccelerationOpts)

			if err := utils.CheckArgs(args); err != nil {
				logger.Error().
					Msg(err.Error())
				return err
			}

			transferAccelerationOpts.DesiredState = "disabled"

			return aws.SetTransferAcceleration(svc, transferAccelerationOpts, logger)
		},
	}
)
