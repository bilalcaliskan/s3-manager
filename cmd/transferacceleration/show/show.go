package show

import (
	"errors"
	"fmt"

	"github.com/bilalcaliskan/s3-manager/internal/aws"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/cmd/transferacceleration/options"
	"github.com/bilalcaliskan/s3-manager/internal/logging"
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
	ShowCmd                  = &cobra.Command{
		Use:           "show",
		Short:         "",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
			svc = cmd.Context().Value(rootopts.S3SvcKey{}).(s3iface.S3API)

			transferAccelerationOpts.RootOptions = rootOpts
			logger = logging.GetLogger(rootOpts)

			if len(args) > 0 {
				err = errors.New("too many arguments provided")
				logger.Error().
					Msg(err.Error())
				return err
			}

			/*ta, err := aws.GetTransferAcceleration(svc, transferAccelerationOpts)
			if err != nil {
				return err
			}*/

			res, err := aws.GetTransferAcceleration(svc, transferAccelerationOpts)
			if err != nil {
				logger.Error().Msg(err.Error())
				return err
			}

			if *res.Status == "Enabled" {
				transferAccelerationOpts.ActualState = "enabled"
			} else if *res.Status == "Suspended" {
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
