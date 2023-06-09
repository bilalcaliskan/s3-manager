package options

import (
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
)

type BucketPolicyOptsKey struct{}

var (
	//substringRunner prompt.PromptRunner = prompt.GetPromptRunner("Provide text to search", nil)
	//extensionRunner prompt.PromptRunner = prompt.GetPromptRunner("Provide target file extensions (comma seperated)", nil)
	bucketPolicyOpts = &BucketPolicyOptions{}
)

// BucketPolicyOptions contains frequent command line and application options.
type BucketPolicyOptions struct {
	BucketPolicyContent string

	*options.RootOptions
}

/*func (opts *BucketPolicyOptions) InitFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&opts.FilePath, "file-path", "", "", "file path to set "+
		"policy into")

	_ = cmd.MarkPersistentFlagRequired("file-path")
}*/

// GetBucketPolicyOptions returns the pointer of FindOptions
func GetBucketPolicyOptions() *BucketPolicyOptions {
	return bucketPolicyOpts
}

func (opts *BucketPolicyOptions) SetZeroValues() {
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
