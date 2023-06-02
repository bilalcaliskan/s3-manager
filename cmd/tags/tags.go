package tags

import (
	"github.com/bilalcaliskan/s3-manager/cmd/tags/add"
	"github.com/bilalcaliskan/s3-manager/cmd/tags/show"
	"github.com/spf13/cobra"
)

func init() {
	/*TagsCmd.AddCommand(show.ShowCmd)
	TagsCmd.AddCommand(set.SetCmd)*/
	TagsCmd.AddCommand(show.ShowCmd)
	TagsCmd.AddCommand(add.AddCmd)
}

var (
	TagsCmd = &cobra.Command{
		Use:           "tags",
		Short:         "",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
)
