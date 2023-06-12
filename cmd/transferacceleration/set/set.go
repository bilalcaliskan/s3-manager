package set

import (
	"github.com/bilalcaliskan/s3-manager/cmd/transferacceleration/set/disabled"
	"github.com/bilalcaliskan/s3-manager/cmd/transferacceleration/set/enabled"
	"github.com/spf13/cobra"
)

func init() {
	SetCmd.AddCommand(enabled.EnabledCmd)
	SetCmd.AddCommand(disabled.DisabledCmd)
}

var (
	SetCmd = &cobra.Command{
		Use:           "set",
		Short:         "sets the transfer acceleration configuration for the target bucket (enabled/disabled)",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
)
