package options

import (
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/spf13/cobra"
)

var (
	searchOptions = &SearchOptions{}
)

// SearchOptions contains frequent command line and application options.
type SearchOptions struct {
	// Text is the target string to search in a bucket
	Text string
	// FileName is the regex or exact name of the target file to search for specific Text
	FileName string

	*options.RootOptions
}

func (opts *SearchOptions) InitFlags(cmd *cobra.Command) {
	if cmd.Name() == "text" {
		cmd.Flags().StringVarP(&opts.FileName, "file-name", "", "", "file-name is the regex "+
			"or exact name of the target file to search for specific text")
	}
}

func (opts *SearchOptions) SetZeroValues() {
	opts.Text = ""
	opts.FileName = ""
}

// GetSearchOptions returns the pointer of FindOptions
func GetSearchOptions() *SearchOptions {
	return searchOptions
}
