package remove

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	options "github.com/bilalcaliskan/s3-manager/cmd/bucketpolicy/options"
	"github.com/bilalcaliskan/s3-manager/cmd/bucketpolicy/utils"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/prompt"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	bucketPolicyOpts = options.GetBucketPolicyOptions()
	//bucketPolicyOpts.InitFlags(RemoveCmd)
}

var (
	svc              s3iface.S3API
	logger           zerolog.Logger
	confirmRunner    prompt.PromptRunner = prompt.GetConfirmRunner()
	bucketPolicyOpts *options.BucketPolicyOptions
	RemoveCmd        = &cobra.Command{
		Use:           "remove",
		Short:         "removes the current bucket policy configuration of the target bucket",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: `# remove the current bucket policy configuration onto target bucket
s3-manager bucketpolicy remove
		`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			svc, bucketPolicyOpts, logger = utils.PrepareConstants(cmd, options.GetBucketPolicyOptions())

			if err := utils.CheckArgs(args); err != nil {
				logger.Error().
					Msg(err.Error())
				return err
			}

			res, err := aws.GetBucketPolicyString(svc, bucketPolicyOpts)
			if err != nil {
				logger.Error().
					Str("error", err.Error()).
					Msg("an error occurred while getting bucket policy")
				return err
			}

			logger.Info().Msg("fetched bucket policy successfully")
			bucketPolicyOpts.BucketPolicyContent = res

			logger.Info().Msg("will attempt to delete below bucket policy")
			fmt.Println(bucketPolicyOpts.BucketPolicyContent)
			if bucketPolicyOpts.DryRun {
				logger.Info().Msg("skipping operation since '--dry-run' flag is passed")
				return nil
			}

			if !bucketPolicyOpts.AutoApprove {
				var res string
				if res, err = confirmRunner.Run(); err != nil {
					return err
				}

				if strings.ToLower(res) == "n" {
					return errors.New("user terminated the process")
				}
			}

			logger.Info().Msg("trying to remove current bucket policy if exists")
			_, err = aws.DeleteBucketPolicy(svc, bucketPolicyOpts)
			if err != nil {
				logger.Error().
					Str("error", err.Error()).
					Msg("an error occurred while deleting current bucket policy")
				return err
			}

			logger.Info().Msg("successfully deleted bucket policy on target bucket")

			return nil
		},
	}
)
