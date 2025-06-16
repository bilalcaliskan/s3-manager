package options

import (
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/spf13/cobra"
)

type ListOptsKey struct{}

var listOpts = &ListOptions{}

// ListOptions is the struct that holds the options for the tags command
type ListOptions struct {
	*options.RootOptions

	StorageClass    string
	MinFileSizeInKB int64
}

func (opts *ListOptions) InitFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&opts.StorageClass, "storageclass", "", "STANDARD",
		"storage class of the objects to list, valid options are \"STANDARD\", \"INTELLIGENT_TIERING\" etc.")
	cmd.Flags().Int64VarP(&opts.MinFileSizeInKB, "min-file-size-in-kb", "", 0,
		"minimum file size in KB to list, 0 means no lower limit")
}

// GetListOptions returns the pointer of ListOptions
func GetListOptions() *ListOptions {
	return listOpts
}

func (opts *ListOptions) SetZeroValues() {
	opts.StorageClass = ""
	opts.MinFileSizeInKB = 0
}
