package tags

import (
	"github.com/bilalcaliskan/s3-manager/cmd/tags/add"
	"github.com/bilalcaliskan/s3-manager/cmd/tags/remove"
	"github.com/bilalcaliskan/s3-manager/cmd/tags/show"
	"github.com/spf13/cobra"
)

func init() {
	TagsCmd.AddCommand(show.ShowCmd)
	TagsCmd.AddCommand(add.AddCmd)
	TagsCmd.AddCommand(remove.RemoveCmd)
}

var (
	TagsCmd = &cobra.Command{
		Use:           "tags",
		Short:         "shows/sets the tagging configuration of the target bucket",
		SilenceUsage:  true,
		SilenceErrors: true,
		// we should not define PreRunE since it will override the PreRunE which is inherited from RootCmd
	}
)
