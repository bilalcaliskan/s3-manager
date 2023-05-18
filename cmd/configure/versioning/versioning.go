package versioning

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/configure/versioning/options"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/logging"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	versioningOpts = options.GetVersioningOptions()
	//versioningOpts.InitFlags(VersioningCmd)
}

var (
	svc            s3iface.S3API
	logger         zerolog.Logger
	versioningOpts *options.VersioningOptions
	VersioningCmd  = &cobra.Command{
		Use:           "versioning",
		Short:         "",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
			svc = cmd.Context().Value(rootopts.S3SvcKey{}).(s3iface.S3API)

			versioningOpts.RootOptions = rootOpts
			logger = logging.GetLogger(rootOpts)

			if len(args) == 0 {
				err = errors.New("you must pass versioning state as enabled or disabled")
				logger.Error().
					Msg(err.Error())
				return err
			}

			if len(args) > 1 {
				err = errors.New("too many arguments")
				logger.Error().
					Msg(err.Error())
				return err
			}

			ver := strings.ToLower(args[0])
			if ver != "enabled" && ver != "disabled" {
				err = errors.New("versioning state must be enabled or disabled")
				logger.Error().
					Msg(err.Error())
				return err
			}

			versioningOpts.State = strings.ToLower(args[0])

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			versioning, err := aws.GetBucketVersioning(svc, versioningOpts)
			if err != nil {
				return err
			}

			if *versioning.Status == "Enabled" && versioningOpts.State == "enabled" || *versioning.Status == "Suspended" && versioningOpts.State == "disabled" {
				logger.Warn().
					Str("state", *versioning.Status).
					Msg("versioning is already at the desired state, skipping configuration")
				return nil
			}

			logger.Info().Msgf("setting versioning as %v", versioningOpts.State)
			_, err = aws.SetBucketVersioning(svc, versioningOpts)
			if err != nil {
				return err
			}

			logger.Info().Msgf("successfully configured versioning as %v", versioningOpts.State)

			return nil
		},
	}
)
