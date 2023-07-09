package show

import (
	"errors"
	"fmt"

	"github.com/bilalcaliskan/s3-manager/internal/utils"

	"github.com/bilalcaliskan/s3-manager/cmd/tags/options"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
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
		Short:         "shows the tagging configuration for the target bucket",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: `# show the current tagging configuration for bucket
s3-manager tags show
		`,
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			var rootOpts *rootopts.RootOptions
			svc, rootOpts, logger = utils.PrepareConstants(cmd)
			tagOpts.RootOptions = rootOpts

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
				fmt.Printf("%s=%s\n", *v.Key, *v.Value)
			}

			return nil
		},
	}
)
