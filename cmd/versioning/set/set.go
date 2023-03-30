package set

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/cmd/versioning/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/logging"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

const (
	ErrTooManyArguments      = "too many arguments. please provide just 'enabled' or 'disabled'"
	ErrWrongArgumentProvided = "wrong argument provided. versioning state must be 'enabled' or 'disabled'"
	ErrNoArgument            = "no argument provided. versioning subcommand takes 'enabled' or 'disabled' argument, please provide one of them"
	ErrUnknownStatus         = "unknown status '%s' returned from AWS SDK"

	WarnDesiredState = "versioning is already at the desired state, skipping configuration"

	InfSuccess           = "successfully configured versioning as %v"
	InfCurrentState      = "current versioning configuration is %s"
	InfSettingVersioning = "setting versioning as %v"
)

func init() {
	versioningOpts = options.GetVersioningOptions()
}

var (
	svc            s3iface.S3API
	logger         zerolog.Logger
	versioningOpts *options.VersioningOptions
	SetCmd         = &cobra.Command{
		Use:           "set",
		Short:         "",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
			svc = cmd.Context().Value(rootopts.S3SvcKey{}).(s3iface.S3API)

			versioningOpts.RootOptions = rootOpts
			logger = logging.GetLogger(rootOpts)

			if err = checkFlags(logger, args); err != nil {
				return err
			}

			versioningOpts.DesiredState = strings.ToLower(args[0])
			versioning, err := aws.GetBucketVersioning(svc, versioningOpts.RootOptions)
			if err != nil {
				return err
			}

			switch *versioning.Status {
			case "Enabled":
				versioningOpts.ActualState = "enabled"
			case "Suspended":
				versioningOpts.ActualState = "disabled"
			default:
				err := fmt.Errorf(ErrUnknownStatus, *versioning.Status)
				logger.Error().Msg(err.Error())
				return err
			}

			logger.Info().Msgf(InfCurrentState, versioningOpts.ActualState)
			if versioningOpts.ActualState == "enabled" && versioningOpts.DesiredState == "enabled" || versioningOpts.ActualState == "disabled" && versioningOpts.DesiredState == "disabled" {
				logger.Warn().
					Str("state", versioningOpts.ActualState).
					Msg(WarnDesiredState)
				return nil
			}

			logger.Info().Msgf(InfSettingVersioning, versioningOpts.DesiredState)
			_, err = aws.SetBucketVersioning(svc, versioningOpts)
			if err != nil {
				return err
			}

			logger.Info().Msgf(InfSuccess, versioningOpts.DesiredState)

			return nil
		},
	}
)
