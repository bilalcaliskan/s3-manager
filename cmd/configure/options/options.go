package options

import (
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/spf13/cobra"
)

var (
	//substringRunner prompt.PromptRunner = prompt.GetPromptRunner("Provide substring to search", nil)
	//extensionRunner prompt.PromptRunner = prompt.GetPromptRunner("Provide target file extensions (comma seperated)", nil)
	configureOptions = &ConfigureOptions{}
)

// ConfigureOptions contains frequent command line and application options.
type ConfigureOptions struct {
	// Versioning is versioning
	Versioning bool
	*options.RootOptions
}

func (opts *ConfigureOptions) InitFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&opts.Versioning, "versioning", "", false,
		"Versioning is versioning")
}

// GetConfigureOptions returns the pointer of FindOptions
func GetConfigureOptions() *ConfigureOptions {
	return configureOptions
}

func (opts *ConfigureOptions) SetZeroValues() {
	opts.Versioning = false
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