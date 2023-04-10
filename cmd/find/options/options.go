package options

import (
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/spf13/cobra"
)

var findOptions = &FindOptions{}

// FindOptions contains frequent command line and application options.
type FindOptions struct {
	// Substring is the target string to find in a bucket
	Substring string
	// FileExtensions is a comma separated list of file extensions to search on S3 bucket (txt, json etc)
	FileExtensions string
	*options.RootOptions
}

func (opts *FindOptions) InitFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&opts.Substring, "substring", "", "",
		"substring to find on txt files on target bucket (default \"\")")
	cmd.Flags().StringVarP(&opts.FileExtensions, "fileExtensions", "", "txt",
		"comma separated list of file extensions to search on S3 bucket")
}

// GetFindOptions returns the pointer of FindOptions
func GetFindOptions() *FindOptions {
	return findOptions
}

func (opts *FindOptions) SetZeroValues() {
	opts.Substring = ""
	opts.FileExtensions = "txt"

	opts.RootOptions.SetZeroValues()
}
