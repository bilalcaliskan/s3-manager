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
	// Foo is foo
	Foo string
	// Bar is bar
	Bar string
	*options.RootOptions
}

func (opts *ConfigureOptions) InitFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&opts.Foo, "foo", "", "bar",
		"foo is foo")
	cmd.Flags().StringVarP(&opts.Bar, "bar", "", "foo",
		"bar is bar")
}

// GetConfigureOptions returns the pointer of FindOptions
func GetConfigureOptions() *ConfigureOptions {
	return configureOptions
}

func (opts *ConfigureOptions) SetZeroValues() {
	opts.Foo = "bar"
	opts.Bar = "foo"
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
