package options

import (
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
)

type TransferAccelerationOptsKey struct{}

var (
	//substringRunner prompt.PromptRunner = prompt.GetPromptRunner("Provide text to search", nil)
	//extensionRunner prompt.PromptRunner = prompt.GetPromptRunner("Provide target file extensions (comma seperated)", nil)
	transferAccelerationOpts = &TransferAccelerationOptions{}
)

// TransferAccelerationOptions contains frequent command line and application options.
type TransferAccelerationOptions struct {
	// ActualState is state
	ActualState string
	// DesiredState is state
	DesiredState string
	*options.RootOptions
}

// GetTransferAccelerationOptions returns the pointer of FindOptions
func GetTransferAccelerationOptions() *TransferAccelerationOptions {
	return transferAccelerationOpts
}

func (opts *TransferAccelerationOptions) SetZeroValues() {
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
