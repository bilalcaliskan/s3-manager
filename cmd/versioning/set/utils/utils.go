package utils

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/internal/logging"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/versioning/options"
)

const (
	ErrTooManyArguments = "too many arguments. please provide just 'enabled' or 'disabled'"
	ErrUnknownStatus    = "unknown status '%s' returned from AWS SDK"

	WarnDesiredState = "versioning is already at the desired state, skipping configuration"

	InfSuccess           = "successfully configured versioning as %v"
	InfCurrentState      = "current versioning configuration is %s"
	InfSettingVersioning = "setting versioning as %v"
)

func CheckArgs(args []string) error {
	if len(args) != 0 {
		return errors.New(ErrTooManyArguments)
	}

	return nil
}

func DecideActualState(versioning *s3.GetBucketVersioningOutput, opts *options.VersioningOptions) error {
	switch *versioning.Status {
	case "Enabled":
		opts.ActualState = "enabled"
	case "Suspended":
		opts.ActualState = "disabled"
	default:
		return fmt.Errorf(ErrUnknownStatus, *versioning.Status)
	}
  
	return nil
}

func PrepareConstants(cmd *cobra.Command, versioningOpts *options.VersioningOptions) (s3iface.S3API, *options.VersioningOptions, zerolog.Logger) {
	svc := cmd.Context().Value(rootopts.S3SvcKey{}).(s3iface.S3API)
	rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
	versioningOpts.RootOptions = rootOpts

	logger := logging.GetLogger(versioningOpts.RootOptions)

	return svc, versioningOpts, logger
}
