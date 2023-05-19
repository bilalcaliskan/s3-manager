package options

import (
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
)

var (
	//substringRunner prompt.PromptRunner = prompt.GetPromptRunner("Provide substring to search", nil)
	//extensionRunner prompt.PromptRunner = prompt.GetPromptRunner("Provide target file extensions (comma seperated)", nil)
	tagOpts = &TagOptions{}
)

// TagOptions contains frequent command line and application options.
type TagOptions struct {
	// Tags is Tags
	ActualTags  map[string]string
	DesiredTags map[string]string
	*options.RootOptions
}

// GetTagOptions returns the pointer of TagOptions
func GetTagOptions() *TagOptions {
	return tagOpts
}

func (opts *TagOptions) SetZeroValues() {
	opts.ActualTags = map[string]string{}
	opts.DesiredTags = map[string]string{}
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
