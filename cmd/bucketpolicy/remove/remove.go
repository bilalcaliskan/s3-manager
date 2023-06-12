package remove

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	options2 "github.com/bilalcaliskan/s3-manager/cmd/bucketpolicy/options"
	"github.com/bilalcaliskan/s3-manager/cmd/bucketpolicy/utils"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	bucketPolicyOpts = options2.GetBucketPolicyOptions()
}

var (
	svc              s3iface.S3API
	logger           zerolog.Logger
	bucketPolicyOpts *options2.BucketPolicyOptions
	RemoveCmd        = &cobra.Command{
		Use:           "remove",
		Short:         "removes the current bucket policy configuration of the target bucket",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			svc, bucketPolicyOpts, logger = utils.PrepareConstants(cmd, options2.GetBucketPolicyOptions())

			if err := utils.CheckArgs(args); err != nil {
				logger.Error().
					Msg(err.Error())
				return err
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
