package versioning

const (
	ErrTooManyArguments      = "too many arguments. please provide just 'enabled' or 'disabled'"
	ErrWrongArgumentProvided = "wrong argument provided. versioning state must be 'enabled' or 'disabled'"
	ErrNoArgument            = "no argument provided. versioning subcommand takes 'enabled' or 'disabled' argument, please provide one of them"

	WarnDesiredState = "versioning is already at the desired state, skipping configuration"

	InfSuccess = "successfully configured versioning as %v"
)
