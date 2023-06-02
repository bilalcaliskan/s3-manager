package show

import (
	"errors"
	"fmt"

	"github.com/bilalcaliskan/s3-manager/cmd/tags/options"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/logging"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	tagOpts = options.GetTagOptions()
}

var (
	svc     s3iface.S3API
	logger  zerolog.Logger
	tagOpts *options.TagOptions
	ShowCmd = &cobra.Command{
		Use:           "show",
		Short:         "",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
			svc = cmd.Context().Value(rootopts.S3SvcKey{}).(s3iface.S3API)

			tagOpts.RootOptions = rootOpts
			logger = logging.GetLogger(rootOpts)

			if len(args) > 0 {
				err = errors.New("too many arguments provided")
				logger.Error().
					Msg(err.Error())
				return err
			}

			tagOpts.ActualTags = make(map[string]string)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			tags, err := aws.GetBucketTags(svc, tagOpts)
			if err != nil {
				logger.Error().
					Str("error", err.Error()).
					Msg("an error occurred while fetching current tags")
				return err
			}

			logger.Info().Msg("fetched bucket tags successfully")

			for _, v := range tags.TagSet {
				tagOpts.ActualTags[*v.Key] = *v.Value
			}

			for key, val := range tagOpts.ActualTags {
				fmt.Printf("%s=%s\n", key, val)
			}

			return nil
		},
	}
)
