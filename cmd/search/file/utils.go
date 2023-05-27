package file

import (
	"errors"

	"github.com/rs/zerolog"
)

const (
	ErrTooManyArguments = "too many arguments. please provide target file name. regex is accepted"
	ErrNoArgument       = "no argument provided. 'file' subcommand takes your desired file name to search in target bucket"
)

func checkFlags(logger zerolog.Logger, args []string) (err error) {
	if len(args) == 0 {
		err = errors.New(ErrNoArgument)
		logger.Error().
			Msg(err.Error())
		return err
	}

	if len(args) > 1 {
		err = errors.New(ErrTooManyArguments)
		logger.Error().
			Msg(err.Error())
		return err
	}

	return nil
}
