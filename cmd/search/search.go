package search

import (
	"github.com/bilalcaliskan/s3-manager/cmd/search/file"
	"github.com/bilalcaliskan/s3-manager/cmd/search/text"

	"github.com/spf13/cobra"
)

func init() {
	SearchCmd.AddCommand(text.TextCmd)
	SearchCmd.AddCommand(file.FileCmd)
}

var (
	SearchCmd = &cobra.Command{
		Use:   "search",
		Short: "searches the files which has desired substrings in it",
		// we should not define PreRunE since it will override the PreRunE which is inherited from RootCmd
		SilenceUsage: true,
	}
)
