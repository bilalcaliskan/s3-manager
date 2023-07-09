package utils

import (
	"errors"
)

const (
	ErrTooManyArguments = "too many arguments"
	ErrNoArgument       = "no argument provided"
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

/*func PrepareConstants(cmd *cobra.Command, searchOptions *options.SearchOptions) (s3iface.S3API, *options.SearchOptions, zerolog.Logger) {
	svc := cmd.Context().Value(rootopts.S3SvcKey{}).(s3iface.S3API)
	rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
	searchOptions.RootOptions = rootOpts

	logger := logging.GetLogger(searchOptions.RootOptions)

	return svc, searchOptions, logger
}*/
