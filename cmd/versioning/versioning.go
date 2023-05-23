package versioning

import (
	"github.com/bilalcaliskan/s3-manager/cmd/versioning/set"
	"github.com/bilalcaliskan/s3-manager/cmd/versioning/show"
	"github.com/spf13/cobra"
)

func init() {
	VersioningCmd.AddCommand(show.ShowCmd)
	VersioningCmd.AddCommand(set.SetCmd)
}

var (
	VersioningCmd = &cobra.Command{
		Use:           "versioning",
		Short:         "",
		SilenceUsage:  true,
		SilenceErrors: true,
		/*PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
			svc = cmd.Context().Value(rootopts.S3SvcKey{}).(s3iface.S3API)
			fmt.Println("hello")

			versioningOpts.RootOptions = rootOpts
			logger = logging.GetLogger(rootOpts)

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
				logger.Error().Msgf(ErrUnknownStatus, *versioning.Status)
			}

			cmd.SetContext(context.WithValue(cmd.Context(), options.VersioningOptsKey{}, versioningOpts))

			return nil
		},*/
		/*RunE: func(cmd *cobra.Command, args []string) (err error) {
			if versioningOpts.ActualState == "enabled" && versioningOpts.DesiredState == "enabled" || versioningOpts.ActualState == "disabled" && versioningOpts.DesiredState == "disabled" {
				logger.Warn().
					Str("state", versioningOpts.ActualState).
					Msg(WarnDesiredState)
				return nil
			}

			logger.Info().Msgf("setting versioning as %v", versioningOpts.DesiredState)
			_, err = aws.SetBucketVersioning(svc, versioningOpts)
			if err != nil {
				return err
			}

			logger.Info().Msgf(InfSuccess, versioningOpts.DesiredState)

			return nil
		},*/
	}
)
