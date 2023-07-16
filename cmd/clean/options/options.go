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
	//FileNamePrefix  string
	Regex          string
	KeepLastNFiles int
	SortBy         string
	Order          string
	*options.RootOptions
}

func (opts *CleanOptions) InitFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&opts.Regex, "regex", "", "",
		"regex is the regex of the target file/folder, as you guess, you can use it to specify a folder or file "+
			"extension also. empty string means all files")
	cmd.Flags().Int64VarP(&opts.MinFileSizeInMb, "min-size-mb", "", 0,
		"minimum size in mb to clean from target bucket, 0 means no lower limit")
	cmd.Flags().Int64VarP(&opts.MaxFileSizeInMb, "max-size-mb", "", 0,
		"maximum size in mb to clean from target bucket, 0 means no upper limit")
	cmd.Flags().IntVarP(&opts.KeepLastNFiles, "keep-last-n-files", "", 2,
		"defines how many of the files to skip deletion in specified criteria, 0 means clean them all")
	cmd.Flags().StringVarP(&opts.SortBy, "sort-by", "", "lastModificationDate",
		"defines the ascending or descending order in the specified criteria, strongly adviced to be used with the "+
			"flag \"--order\", valid options are \"lastModificationDate\" and \"size\"")
	cmd.Flags().StringVarP(&opts.Order, "order", "", "descending",
		"specifies the ordering strategy to sort objects in the \"--sort-by\" flag, valid options are \"ascending\" and \"descending\"")
}

func (opts *CleanOptions) SetZeroValues() {
	opts.Regex = ""
	opts.MinFileSizeInMb = 0
	opts.MaxFileSizeInMb = 0
	opts.KeepLastNFiles = 2
	opts.SortBy = "lastModificationDate"
	opts.Order = "descending"
}

// GetCleanOptions returns the pointer of CleanOptions
func GetCleanOptions() *CleanOptions {
	return cleanOptions
}
