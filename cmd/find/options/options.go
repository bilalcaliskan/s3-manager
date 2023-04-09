package options

import (
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/spf13/cobra"
)

var findOptions = &FindOptions{}

// FindOptions contains frequent command line and application options.
type FindOptions struct {
	*options.RootOptions
}

func (opts *FindOptions) InitFlags(cmd *cobra.Command) {
	/*cmd.Flags().Int64VarP(&opts.MinFileSizeInMb, "minFileSizeInMb", "", 0,
	"minimum size in mb to clean from target bucket, 0 means no lower limit")*/

	//findOptions.RootOptions = options.GetRootOptions()
}

// GetFindOptions returns the pointer of FindOptions
func GetFindOptions() *FindOptions {
	findOptions.RootOptions = options.GetRootOptions()

	return findOptions
}

func (opts *FindOptions) SetZeroValues() {

	opts.RootOptions = options.GetRootOptions()
}
