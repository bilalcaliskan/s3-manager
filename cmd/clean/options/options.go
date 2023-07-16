package options

import (
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/spf13/cobra"
)

var cleanOptions = &CleanOptions{}

// CleanOptions contains frequent command line and application options.
type CleanOptions struct {
	MinFileSizeInMb int64
	MaxFileSizeInMb int64
	FileExtensions  string
	FileNamePrefix  string
	KeepLastNFiles  int
	SortBy          string
	*options.RootOptions
}

func (opts *CleanOptions) InitFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&opts.FileNamePrefix, "file-name-prefix", "", "",
		"folder name of target bucket objects, means it can be used for folder-based object grouping buckets (default \"\")")
	cmd.Flags().Int64VarP(&opts.MinFileSizeInMb, "min-size-mb", "", 0,
		"minimum size in mb to clean from target bucket, 0 means no lower limit")
	cmd.Flags().Int64VarP(&opts.MaxFileSizeInMb, "max-size-mb", "", 0,
		"maximum size in mb to clean from target bucket, 0 means no upper limit")
	cmd.Flags().IntVarP(&opts.KeepLastNFiles, "keep-last-n-files", "", 2,
		"defines how many of the files to skip deletion in specified criteria, 0 means clean them all")
	cmd.Flags().StringVarP(&opts.SortBy, "sort-by", "", "lastModificationDate",
		"defines the ascending order in the specified criteria, valid options are \"lastModificationDate\" and \"size\"")
}

// GetCleanOptions returns the pointer of CleanOptions
func GetCleanOptions() *CleanOptions {
	return cleanOptions
}
