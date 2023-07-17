package constants

import "errors"

var (
	ErrInjected       = errors.New("injected error")
	ErrUserTerminated = errors.New("user terminated the process")
	ErrInvalidInput   = errors.New("invalid input")
)
