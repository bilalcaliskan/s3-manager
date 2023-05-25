package substring

import (
	"errors"
	"github.com/rs/zerolog"
)

const (
	ErrTooManyArguments = "too many arguments. please provide just your desired substring to search in --fileExtensions files"
	ErrNoArgument       = "no argument provided. 'substring' subcommand takes your desired substring to search in --fileExtensions files"
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
