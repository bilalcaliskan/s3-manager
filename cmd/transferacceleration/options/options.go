package options

import (
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
)

type TransferAccelerationOptsKey struct{}

var (
	transferAccelerationOpts = &TransferAccelerationOptions{}
)

type TransferAccelerationOptions struct {
	ActualState  string
	DesiredState string
	*options.RootOptions
}

func GetTransferAccelerationOptions() *TransferAccelerationOptions {
	return transferAccelerationOpts
}
