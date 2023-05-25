package options

import (
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/spf13/cobra"
)

var (
	// TODO: uncomment when interactivity enabled again
	/*substringRunner prompt.PromptRunner = prompt.GetPromptRunner("Provide substring to search", nil)
	extensionRunner prompt.PromptRunner = prompt.GetPromptRunner("Provide target file extensions (comma seperated)", nil)*/
	searchOptions = &SearchOptions{}
)

// SearchOptions contains frequent command line and application options.
type SearchOptions struct {
	// Substring is the target string to search in a bucket
	Substring string
	// FileExtensions is a comma separated list of file extensions to search on S3 bucket (txt, json etc)
	FileExtensions string
	*options.RootOptions
}

func (opts *SearchOptions) InitFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&opts.FileExtensions, "fileExtensions", "", "txt",
		"comma separated list of file extensions to search on S3 bucket")
}

// GetFindOptions returns the pointer of FindOptions
func GetSearchOptions() *SearchOptions {
	return searchOptions
}

func (opts *SearchOptions) SetZeroValues() {
	opts.FileExtensions = "txt"
	opts.RootOptions.SetZeroValues()
}

// TODO: uncomment when interactivity enabled again
/*func (opts *SearchOptions) PromptInteractiveValues() error {
	res, err := substringRunner.Run()
	if err != nil {
		return err
	}
	opts.Substring = res

	res, err = extensionRunner.Run()
	if err != nil {
		return err
	}
	opts.FileExtensions = res

	return nil
}*/
