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
		Use:           "search",
		Short:         "searches the files which has desired substrings in it",
		SilenceUsage:  false,
		SilenceErrors: false,
	}
)
