package transferacceleration

import (
	"github.com/bilalcaliskan/s3-manager/cmd/transferacceleration/set"
	"github.com/bilalcaliskan/s3-manager/cmd/transferacceleration/show"
	"github.com/spf13/cobra"
)

func init() {
	TransferAccelerationCmd.AddCommand(show.ShowCmd)
	TransferAccelerationCmd.AddCommand(set.SetCmd)
}

var (
	TransferAccelerationCmd = &cobra.Command{
		Use:           "transferacceleration",
		Short:         "shows/sets the transfer acceleration configuration of the target bucket",
		SilenceUsage:  true,
		SilenceErrors: true,
		// we should not define PreRunE since it will override the PreRunE which is inherited from RootCmd
	}
)
