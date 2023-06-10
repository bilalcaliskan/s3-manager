package show

import (
	"encoding/json"
	"fmt"

	"github.com/bilalcaliskan/s3-manager/cmd/bucketpolicy/utils"

	options2 "github.com/bilalcaliskan/s3-manager/cmd/bucketpolicy/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
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
	ShowCmd          = &cobra.Command{
		Use:           "show",
		Short:         "",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			svc, bucketPolicyOpts, logger = utils.PrepareConstants(cmd, options2.GetBucketPolicyOptions())

			if err := utils.CheckArgs(args); err != nil {
				logger.Error().
					Msg(err.Error())
				return err
			}

			res, err := aws.GetBucketPolicy(svc, bucketPolicyOpts)
			if err != nil {
				logger.Error().
					Str("error", err.Error()).
					Msg("an error occurred while getting bucket policy")
				return err
			}

			logger.Info().Msg("fetched bucket policy successfully")

			beautifiedJSON, err := beautifyJSON(*res.Policy)
			if err != nil {
				logger.Error().Msg(err.Error())
				return err
			}

			fmt.Println(beautifiedJSON)

			return nil
		},
	}
)

func beautifyJSON(jsonString string) (string, error) {
	var jsonData interface{}

	err := json.Unmarshal([]byte(jsonString), &jsonData)
	if err != nil {
		return "", err
	}

	beautifiedBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return "", err
	}

	return string(beautifiedBytes), nil
}
