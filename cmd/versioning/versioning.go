package versioning

import (
	"github.com/bilalcaliskan/s3-manager/cmd/versioning/set"
	"github.com/bilalcaliskan/s3-manager/cmd/versioning/show"
	"github.com/spf13/cobra"
)

func init() {
	VersioningCmd.AddCommand(show.ShowCmd)
	VersioningCmd.AddCommand(set.SetCmd)
}

var (
	VersioningCmd = &cobra.Command{
		Use:           "versioning",
		Short:         "shows/sets the versioning configuration of the target bucket",
		SilenceUsage:  false,
		SilenceErrors: false,
	}
)
