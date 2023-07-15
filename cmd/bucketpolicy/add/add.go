package add

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bilalcaliskan/s3-manager/internal/constants"

	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/bucketpolicy/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/prompt"
	"github.com/bilalcaliskan/s3-manager/internal/utils"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	bucketPolicyOpts = options.GetBucketPolicyOptions()
}

var (
	svc              s3iface.S3API
	logger           zerolog.Logger
	confirmRunner    prompt.PromptRunner
	bucketPolicyOpts *options.BucketPolicyOptions
	AddCmd           = &cobra.Command{
		Use:           "add",
		Short:         "adds a bucket policy configuration for the target bucket by specifying a valid policy file",
		SilenceUsage:  false,
		SilenceErrors: true,
		Example: `# add a bucket policy configuration onto target bucket
s3-manager bucketpolicy add my_custom_policy.json
		`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var rootOpts *rootopts.RootOptions
			svc, rootOpts, logger, confirmRunner = utils.PrepareConstants(cmd)
			bucketPolicyOpts.RootOptions = rootOpts

			if err := utils.CheckArgs(args, 1); err != nil {
				logger.Error().Msg(err.Error())
				return err
			}

			logger = logger.With().Str("policyFilePath", args[0]).Logger()

			logger.Info().Msg("trying to read target policy file")
			file, err := os.Open(args[0])
			if err != nil {
				logger.Error().Msg(err.Error())
				return err
			}

			defer func() {
				if err := file.Close(); err != nil {
					panic(err)
				}
			}()

			content, err := io.ReadAll(file)
			if err != nil {
				logger.Error().Msg(err.Error())
				return err
			}

			logger.Info().Msg("successfully read target policy file")
			bucketPolicyOpts.BucketPolicyContent = string(content)

			logger.Info().Msg("will attempt to add below bucket policy")
			fmt.Println(bucketPolicyOpts.BucketPolicyContent)
			if bucketPolicyOpts.DryRun {
				logger.Info().Msg("skipping operation since '--dry-run' flag is passed")
				return nil
			}

			if !bucketPolicyOpts.AutoApprove {
				var res string
				if res, err = confirmRunner.Run(); err != nil {
					if strings.ToLower(res) == "n" {
						return constants.ErrUserTerminated
					}

					return constants.ErrInvalidInput
				}
			}

			logger.Info().Msg("trying to add bucket policy")
			_, err = aws.SetBucketPolicy(svc, bucketPolicyOpts)
			if err != nil {
				logger.Error().
					Str("error", err.Error()).
					Msg("an error occurred while setting bucket policy")
				return err
			}

			logger.Info().Msg("successfully set bucket policy with target file content on target bucket")

			return nil
		},
	}
)
