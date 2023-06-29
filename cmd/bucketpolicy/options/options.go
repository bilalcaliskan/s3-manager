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

// GetBucketPolicyOptions returns the pointer of FindOptions
func GetBucketPolicyOptions() *BucketPolicyOptions {
	return bucketPolicyOpts
}

/*func (opts *BucketPolicyOptions) InitFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&opts.AutoApprove, "auto-approve", "", false, "Skip interactive approval (default false)")
	cmd.Flags().BoolVarP(&opts.DryRun, "dry-run", "", false, "specifies that if you "+
		"just want to see on what content to take action (default false)")
}*/

func (opts *BucketPolicyOptions) SetZeroValues() {
	opts.BucketPolicyContent = ""
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
