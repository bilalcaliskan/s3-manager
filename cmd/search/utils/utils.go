package utils

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/cmd/search/options"
	"github.com/bilalcaliskan/s3-manager/internal/logging"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

const (
	ErrTooManyArguments = "too many arguments. please provide just your desired text to search in --fileExtensions files"
	ErrNoArgument       = "no argument provided. 'text' subcommand takes your desired text to search in --fileExtensions files"
)

func CheckFlags(args []string) (err error) {
	if len(args) == 0 {
		return errors.New(ErrNoArgument)
	}

	if len(args) > 1 {
		return errors.New(ErrTooManyArguments)
	}

	return nil
}

func PrepareConstants(cmd *cobra.Command, searchOptions *options.SearchOptions) (s3iface.S3API, *options.SearchOptions, zerolog.Logger) {
	svc := cmd.Context().Value(rootopts.S3SvcKey{}).(s3iface.S3API)
	rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
	searchOptions.RootOptions = rootOpts

	logger := logging.GetLogger(searchOptions.RootOptions)

	return svc, searchOptions, logger
}
