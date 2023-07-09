package set

import (
	"github.com/bilalcaliskan/s3-manager/cmd/versioning/set/disabled"
	"github.com/bilalcaliskan/s3-manager/cmd/versioning/set/enabled"
	"github.com/spf13/cobra"
)

func init() {
	SetCmd.AddCommand(enabled.EnabledCmd)
	SetCmd.AddCommand(disabled.DisabledCmd)
}

var (
	SetCmd = &cobra.Command{
		Use:           "set",
		Short:         "sets the versioning configuration for the target bucket (enabled/disabled)",
		SilenceUsage:  false,
		SilenceErrors: false,
	}
)
