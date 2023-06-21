package utils

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/bucketpolicy/options"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/internal/logging"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

const errTooManyArgsProvided = "too many arguments provided"

func CheckArgs(args []string) error {
	if len(args) > 0 {
		return errors.New(errTooManyArgsProvided)
	}

	return nil
}

func PrepareConstants(cmd *cobra.Command, tagOpts *options.BucketPolicyOptions) (s3iface.S3API, *options.BucketPolicyOptions, zerolog.Logger) {
	svc := cmd.Context().Value(rootopts.S3SvcKey{}).(s3iface.S3API)
	rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
	tagOpts.RootOptions = rootOpts

	logger := logging.GetLogger(tagOpts.RootOptions)

	return svc, tagOpts, logger
}
