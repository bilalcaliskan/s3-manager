package search

import (
	"github.com/bilalcaliskan/s3-manager/cmd/search/file"
	"github.com/bilalcaliskan/s3-manager/cmd/search/substring"

	"github.com/spf13/cobra"
)

func init() {
	SearchCmd.AddCommand(substring.SubstringCmd)
	SearchCmd.AddCommand(file.FileCmd)
}

var (
	SearchCmd = &cobra.Command{
		Use:          "search",
		Short:        "search subcommand searches the files which has desired substrings in it",
		SilenceUsage: true,
	}
)
