package add

import (
	"errors"
	"io"
	"os"

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
	AddCmd           = &cobra.Command{
		Use:           "add",
		Short:         "",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			svc, bucketPolicyOpts, logger = utils.PrepareConstants(cmd, options2.GetBucketPolicyOptions())

			if len(args) < 1 {
				err = errors.New("no argument provided, you should provide bucket policy path on your filesystem")
				logger.Error().Msg(err.Error())
				return err
			} else if len(args) > 1 {
				err = errors.New("too many argument provided, just provide bucket policy path on your filesystem")
				logger.Error().Msg(err.Error())
				return err
			}

			/*if err := utils.CheckArgs(cmd, args); err != nil {
				logger.Error().Msg(err.Error())
				return err
			}*/

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

			// Read the file's content
			content, err := io.ReadAll(file)
			if err != nil {
				logger.Error().Msg(err.Error())
				return err
			}

			bucketPolicyOpts.BucketPolicyContent = string(content)
			logger.Info().Msg("successfully read target policy file")

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
