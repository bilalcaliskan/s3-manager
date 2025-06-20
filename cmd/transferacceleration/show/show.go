package show

import (
	"fmt"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/aws"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/pkg/aws/types"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/utils"

	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/cmd/transferacceleration/options"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	transferAccelerationOpts = options.GetTransferAccelerationOptions()
}

var (
	svc                      internalawstypes.S3ClientAPI
	logger                   zerolog.Logger
	transferAccelerationOpts *options.TransferAccelerationOptions
	ShowCmd                  = &cobra.Command{
		Use:           "show",
		Short:         "shows the transfer acceleration configuration for the target bucket",
		SilenceUsage:  false,
		SilenceErrors: true,
		Example: `# show the current transfer acceleration configuration for bucket
s3-manager transferacceleration show
		`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var rootOpts *rootopts.RootOptions
			svc, rootOpts, logger, _ = utils.PrepareConstants(cmd)
			transferAccelerationOpts.RootOptions = rootOpts

			if err := utils.CheckArgs(args, 0); err != nil {
				logger.Error().
					Msg(err.Error())
				return err
			}

			res, err := aws.GetTransferAcceleration(svc, transferAccelerationOpts)
			if err != nil {
				logger.Error().Msg(err.Error())
				return err
			}

			if res.Status == "Enabled" {
				transferAccelerationOpts.ActualState = "enabled"
			} else if res.Status == "Suspended" {
				transferAccelerationOpts.ActualState = "disabled"
			} else {
				err := fmt.Errorf("unknown status '%s' returned from AWS SDK", transferAccelerationOpts.ActualState)
				logger.Error().Msg(err.Error())
				return err
			}

			logger.Info().Msgf("current transfer acceleration configuration is %s", transferAccelerationOpts.ActualState)

			return nil
		},
	}
)
