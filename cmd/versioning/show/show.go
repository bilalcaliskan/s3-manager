package show

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/cmd/versioning/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/logging"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	versioningOpts = options.GetVersioningOptions()
}

var (
	svc            s3iface.S3API
	logger         zerolog.Logger
	versioningOpts *options.VersioningOptions
	ShowCmd        = &cobra.Command{
		Use:           "show",
		Short:         "shows the versioning configuration for the target bucket",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: `# show the current versioning configuration for bucket
s3-manager versioning show
		`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
			svc = cmd.Context().Value(rootopts.S3SvcKey{}).(s3iface.S3API)

			versioningOpts.RootOptions = rootOpts
			logger = logging.GetLogger(rootOpts)

			if len(args) > 0 {
				err = errors.New("too many arguments provided")
				logger.Error().
					Msg(err.Error())
				return err
			}

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
				err := fmt.Errorf("unknown versioning status %s returned from S3 SDK", *versioning.Status)
				logger.Error().Msg(err.Error())
				return err
			}

			logger.Info().Msgf("current versioning configuration is %s", versioningOpts.ActualState)

			return nil
		},
	}
)
