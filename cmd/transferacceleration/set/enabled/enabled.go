package enabled

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/transferacceleration/options"
	"github.com/bilalcaliskan/s3-manager/cmd/transferacceleration/utils"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	transferAccelerationOpts = options.GetTransferAccelerationOptions()
}

var (
	svc                      s3iface.S3API
	logger                   zerolog.Logger
	transferAccelerationOpts *options.TransferAccelerationOptions
	EnabledCmd               = &cobra.Command{
		Use:           "enabled",
		Short:         "enables the transfer acceleration configuration for the target bucket",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, transferAccelerationOpts, logger = utils.PrepareConstants(cmd, transferAccelerationOpts)

			if err := utils.CheckArgs(args); err != nil {
				logger.Error().
					Msg(err.Error())
				return err
			}

			transferAccelerationOpts.DesiredState = "enabled"

			return aws.SetTransferAcceleration(svc, transferAccelerationOpts, logger)
		},
	}
)
