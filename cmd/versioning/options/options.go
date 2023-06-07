package options

import (
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
)

type VersioningOptsKey struct{}

var (
	//substringRunner prompt.PromptRunner = prompt.GetPromptRunner("Provide text to search", nil)
	//extensionRunner prompt.PromptRunner = prompt.GetPromptRunner("Provide target file extensions (comma seperated)", nil)
	versioningOpts = &VersioningOptions{}
)

// VersioningOptions contains frequent command line and application options.
type VersioningOptions struct {
	// ActualState is state
	ActualState string
	// DesiredState is state
	DesiredState string
	*options.RootOptions
}

// GetVersioningOptions returns the pointer of FindOptions
func GetVersioningOptions() *VersioningOptions {
	return versioningOpts
}

func (opts *VersioningOptions) SetZeroValues() {
	opts.ActualState = "Enabled"
	opts.DesiredState = "enabled"
	opts.RootOptions.SetZeroValues()
}

/*func (opts *ConfigureOptions) PromptInteractiveValues() error {
	res, err := substringRunner.Run()
	if err != nil {
		return err
	}
	opts.Foo = res

	res, err = extensionRunner.Run()
	if err != nil {
		return err
	}
	opts.FileExtensions = res

	return nil
}
*/
