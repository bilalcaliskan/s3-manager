package options

import (
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
)

type TagOptsKey struct{}

var (
	//substringRunner prompt.PromptRunner = prompt.GetPromptRunner("Provide substring to search", nil)
	//extensionRunner prompt.PromptRunner = prompt.GetPromptRunner("Provide target file extensions (comma seperated)", nil)
	tagOpts = &TagOptions{}
)

// TagOptions contains frequent command line and application options.
type TagOptions struct {
	// ActualState is state
	ActualTags map[string]string
	// TagsToAdd is state
	TagsToAdd map[string]string
	// TagsToRemove is state
	TagsToRemove map[string]string
	*options.RootOptions
}

// GetTagOptions returns the pointer of TagOptions
func GetTagOptions() *TagOptions {
	return tagOpts
}

func (opts *TagOptions) SetZeroValues() {
	opts.ActualTags = make(map[string]string)
	opts.TagsToRemove = make(map[string]string)
	opts.TagsToAdd = make(map[string]string)
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
