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
	DryRun          bool
	AutoApprove     bool
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
	/*cmd.Flags().StringVarP(&opts.FileExtensions, "fileExtensions", "", "",
	"selects the files with defined extensions to clean from target bucket, \"\" means all files (default \"\")")*/
	cmd.Flags().IntVarP(&opts.KeepLastNFiles, "keep-last-n-files", "", 2,
		"defines how many of the files to skip deletion in specified criteria, 0 means clean them all")
	cmd.Flags().StringVarP(&opts.SortBy, "sort-by", "", "lastModificationDate",
		"defines the ascending order in the specified criteria, valid options are \"lastModificationDate\" and \"size\"")
	cmd.Flags().BoolVarP(&opts.AutoApprove, "auto-approve", "", false, "Skip interactive approval (default false)")
	cmd.Flags().BoolVarP(&opts.DryRun, "dry-run", "", false, "specifies that if you "+
		"just want to see what to delete or completely delete them all (default false)")
}

// GetCleanOptions returns the pointer of CleanOptions
func GetCleanOptions() *CleanOptions {
	return cleanOptions
}

func (opts *CleanOptions) SetZeroValues() {
	opts.MinFileSizeInMb = 0
	opts.MaxFileSizeInMb = 0
	opts.FileExtensions = ""
	opts.KeepLastNFiles = 2
	opts.DryRun = false
	opts.AutoApprove = false
	opts.SortBy = "lastModificationDate"
	opts.RootOptions.SetZeroValues()
}
