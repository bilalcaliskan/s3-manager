package configure

import (
	"github.com/bilalcaliskan/s3-manager/cmd/configure/versioning"
	"github.com/spf13/cobra"
)

func init() {
	ConfigureCmd.AddCommand(versioning.VersioningCmd)
}

var (
	ConfigureCmd = &cobra.Command{
		Use:          "configure",
		Short:        "configure subcommand configures the bucket level settings",
		SilenceUsage: true,
	}
)
