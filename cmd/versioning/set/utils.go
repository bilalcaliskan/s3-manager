package set

import (
	"errors"
	"strings"

	"github.com/rs/zerolog"
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

	ver := strings.ToLower(args[0])
	if ver != "enabled" && ver != "disabled" {
		err = errors.New(ErrWrongArgumentProvided)
		logger.Error().
			Msg(err.Error())
		return err
	}

	return nil
}
