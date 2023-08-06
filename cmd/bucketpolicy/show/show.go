package show

import (
	"fmt"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/aws/types"

	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"

	"github.com/bilalcaliskan/s3-manager/internal/utils"

	options2 "github.com/bilalcaliskan/s3-manager/cmd/bucketpolicy/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	bucketPolicyOpts = options2.GetBucketPolicyOptions()
}

var (
	svc              internalawstypes.S3ClientAPI
	logger           zerolog.Logger
	bucketPolicyOpts *options2.BucketPolicyOptions
	ShowCmd          = &cobra.Command{
		Use:           "show",
		Short:         "shows the bucket policy configuration of the target bucket",
		SilenceUsage:  false,
		SilenceErrors: true,
		Example: `# show the current bucket policy configuration for target bucket
s3-manager bucketpolicy show
		`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var rootOpts *rootopts.RootOptions
			svc, rootOpts, logger, _ = utils.PrepareConstants(cmd)
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
			fmt.Println(res)

			return nil
		},
	}
)
