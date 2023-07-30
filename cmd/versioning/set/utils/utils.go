package utils

import (
	"fmt"

	v2s3 "github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/bilalcaliskan/s3-manager/cmd/versioning/options"
)

const (
	ErrUnknownStatus = "unknown status '%s' returned from AWS SDK"

	WarnDesiredState = "versioning is already at the desired state, skipping configuration"

	InfSuccess           = "successfully configured versioning as %v"
	InfCurrentState      = "current versioning configuration is %s"
	InfSettingVersioning = "setting versioning as %v"
)

func DecideActualState(versioning *v2s3.GetBucketVersioningOutput, opts *options.VersioningOptions) error {
	switch versioning.Status {
	case "Enabled":
		opts.ActualState = "enabled"
	case "Suspended":
		opts.ActualState = "disabled"
	default:
		return fmt.Errorf(ErrUnknownStatus, versioning.Status)
	}

	return nil
}
