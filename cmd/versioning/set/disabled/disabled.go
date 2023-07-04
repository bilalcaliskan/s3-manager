package disabled

import (
	"strings"

	"github.com/bilalcaliskan/s3-manager/cmd/versioning/set/utils"
	"github.com/bilalcaliskan/s3-manager/internal/prompt"
	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/versioning/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	versioningOpts = options.GetVersioningOptions()
}

var (
	svc            s3iface.S3API
	logger         zerolog.Logger
	confirmRunner  prompt.PromptRunner = prompt.GetConfirmRunner()
	versioningOpts *options.VersioningOptions
	DisabledCmd    = &cobra.Command{
		Use:           "disabled",
		Short:         "disables the versioning configuration for the target bucket",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: `# set the versioning configuration for bucket as disabled
s3-manager versioning set disabled
		`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			svc, versioningOpts, logger = utils.PrepareConstants(cmd, options.GetVersioningOptions())

			if err := utils.CheckArgs(args); err != nil {
				logger.Error().
					Msg(err.Error())
				return err
			}

			versioningOpts.DesiredState = "disabled"

			if versioningOpts.DryRun {
				logger.Info().Msg("skipping operation since '--dry-run' flag is passed")
				return nil
			}

			if !versioningOpts.AutoApprove {
				var res string
				if res, err = confirmRunner.Run(); err != nil {
					return err
				}

				if strings.ToLower(res) == "n" {
					return errors.New("user terminated the process")
				}
			}

			return aws.SetBucketVersioning(svc, versioningOpts, logger)
		},
	}
)
