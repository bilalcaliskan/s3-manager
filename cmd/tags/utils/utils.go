package utils

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	options2 "github.com/bilalcaliskan/s3-manager/cmd/tags/options"
	"github.com/bilalcaliskan/s3-manager/internal/logging"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func CheckArgs(cmd *cobra.Command, args []string) error {
	if cmd.Name() == "show" && len(args) == 0 {
		return nil
	} else if cmd.Name() == "show" && len(args) > 0 {
		return errors.New("too many argument provided")
	}

	fmt.Printf("here")
	fmt.Println(len(args))

	if len(args) == 0 {
		return errors.New("no argument provided")
	} else if len(args) > 1 {
		return errors.New("too many argument provided")
	}

	return nil
}

func PrepareConstants(cmd *cobra.Command, tagOpts *options2.TagOptions) (s3iface.S3API, *options2.TagOptions, zerolog.Logger) {
	svc := cmd.Context().Value(rootopts.S3SvcKey{}).(s3iface.S3API)
	rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
	tagOpts.RootOptions = rootOpts

	logger := logging.GetLogger(tagOpts.RootOptions)

	return svc, tagOpts, logger
}
