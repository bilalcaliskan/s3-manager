package options

import (
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
)

var (
	//substringRunner prompt.PromptRunner = prompt.GetPromptRunner("Provide substring to search", nil)
	//extensionRunner prompt.PromptRunner = prompt.GetPromptRunner("Provide target file extensions (comma seperated)", nil)
	versioningOpts = &VersioningOptions{}
)

// VersioningOptions contains frequent command line and application options.
type VersioningOptions struct {
	// State is state
	State string
	*options.RootOptions
}

// GetVersioningOptions returns the pointer of FindOptions
func GetVersioningOptions() *VersioningOptions {
	return versioningOpts
}

func (opts *VersioningOptions) SetZeroValues() {
	opts.State = "enabled"
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
