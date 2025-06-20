package remove

import (
	"fmt"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/aws"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/pkg/aws/types"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/prompt"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/utils"

	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"

	"github.com/bilalcaliskan/s3-manager/cmd/bucketpolicy/options"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	bucketPolicyOpts = options.GetBucketPolicyOptions()
}

var (
	svc              internalawstypes.S3ClientAPI
	logger           zerolog.Logger
	confirmRunner    prompt.PromptRunner
	bucketPolicyOpts *options.BucketPolicyOptions
	RemoveCmd        = &cobra.Command{
		Use:           "remove",
		Short:         "removes the current bucket policy configuration of the target bucket",
		SilenceUsage:  false,
		SilenceErrors: true,
		Example: `# remove the current bucket policy configuration onto target bucket
s3-manager bucketpolicy remove
		`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var rootOpts *rootopts.RootOptions
			svc, rootOpts, logger, confirmRunner = utils.PrepareConstants(cmd)
			bucketPolicyOpts.RootOptions = rootOpts

			if err := utils.CheckArgs(args, 0); err != nil {
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

			logger.Info().Msg("trying to remove current bucket policy if exists")
			_, err = aws.DeleteBucketPolicy(svc, bucketPolicyOpts, confirmRunner, logger)
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
