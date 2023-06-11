package utils

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	options2 "github.com/bilalcaliskan/s3-manager/cmd/transferacceleration/options"
	"github.com/bilalcaliskan/s3-manager/internal/logging"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	ErrTooManyArguments = "too many arguments. please provide just 'enabled' or 'disabled'"
	ErrUnknownStatus    = "unknown status '%s' returned from AWS SDK"
)

func CheckArgs(args []string) error {
	if len(args) != 0 {
		return errors.New(ErrTooManyArguments)
	}

	return nil
}

func DecideActualState(res *s3.GetBucketAccelerateConfigurationOutput, opts *options2.TransferAccelerationOptions) error {
	if *res.Status == "Enabled" {
		opts.ActualState = "enabled"
	} else if *res.Status == "Suspended" {
		opts.ActualState = "disabled"
	} else {
		return fmt.Errorf(ErrUnknownStatus, opts.ActualState)
	}

	return nil
}

func PrepareConstants(cmd *cobra.Command, transferAccelerationOpts *options2.TransferAccelerationOptions) (s3iface.S3API, *options2.TransferAccelerationOptions, zerolog.Logger) {
	svc := cmd.Context().Value(rootopts.S3SvcKey{}).(s3iface.S3API)
	rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
	transferAccelerationOpts.RootOptions = rootOpts

	logger := logging.GetLogger(transferAccelerationOpts.RootOptions)

	return svc, transferAccelerationOpts, logger
}
