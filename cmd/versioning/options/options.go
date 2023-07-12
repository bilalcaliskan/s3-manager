package options

import (
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
)

type VersioningOptsKey struct{}

var versioningOpts = &VersioningOptions{}

// VersioningOptions contains frequent command line and application options.
type VersioningOptions struct {
	ActualState  string
	DesiredState string
	*options.RootOptions
}

func GetVersioningOptions() *VersioningOptions {
	return versioningOpts
}

func (opts *VersioningOptions) SetZeroValues() {
	opts.ActualState = "Enabled"
	opts.DesiredState = "enabled"
	opts.RootOptions.SetZeroValues()
}
