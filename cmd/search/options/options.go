package options

import (
	"errors"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var searchOptions = &SearchOptions{}

// SearchOptions contains frequent command line and application options.
type SearchOptions struct {
	// Substring is the target string to search in a bucket
	Substring string
	// FileExtensions is a comma separated list of file extensions to search on S3 bucket (txt, json etc)
	FileExtensions string
	*options.RootOptions
}

func (opts *SearchOptions) InitFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&opts.Substring, "substring", "", "",
		"substring to search on txt files on target bucket (default \"\")")
	cmd.Flags().StringVarP(&opts.FileExtensions, "fileExtensions", "", "txt",
		"comma separated list of file extensions to search on S3 bucket")
}

// GetFindOptions returns the pointer of FindOptions
func GetSearchOptions() *SearchOptions {
	return searchOptions
}

func (opts *SearchOptions) SetZeroValues() {
	opts.Substring = ""
	opts.FileExtensions = "txt"

	opts.RootOptions.SetZeroValues()
}

func (opts *SearchOptions) PromptInteractiveValues() error {
	prompt := promptui.Prompt{
		Label: "Substring to search",
		Validate: func(s string) error {
			if len(s) > 50 {
				return errors.New("to long substring to search")
			}

			return nil
		},
	}

	res, err := prompt.Run()
	if err != nil {
		return err
	}
	opts.Substring = res

	prompt = promptui.Prompt{
		Label: "Target file extensions (comma seperated)",
	}
	res, err = prompt.Run()
	if err != nil {
		return err
	}
	opts.FileExtensions = res

	return nil
}
