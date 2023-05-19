package tags

import (
	"errors"
	"strings"

	"github.com/bilalcaliskan/s3-manager/internal/aws"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/configure/tags/options"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/internal/logging"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	tagOpts = options.GetTagOptions()
	//versioningOpts.InitFlags(VersioningCmd)
}

var (
	svc     s3iface.S3API
	logger  zerolog.Logger
	tagOpts *options.TagOptions
	TagsCmd = &cobra.Command{
		Use:           "tags",
		Short:         "",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			rootOpts := cmd.Context().Value(rootopts.OptsKey{}).(*rootopts.RootOptions)
			svc = cmd.Context().Value(rootopts.S3SvcKey{}).(s3iface.S3API)

			tagOpts.RootOptions = rootOpts
			tagOpts.ActualTags = make(map[string]string)
			tagOpts.DesiredTags = make(map[string]string)
			logger = logging.GetLogger(rootOpts)

			if len(args) == 0 {
				err = errors.New(ErrNoArgument)
				logger.Error().
					Msg(err.Error())
				return err
			}

			if len(args) > 1 {
				err = errors.New(ErrTooManyArguments)
				logger.Error().
					Msg(err.Error())
				return err
			}

			// TODO: split tags and append to array
			tagArr := strings.Split(args[0], ",")
			for _, v := range tagArr {
				key := strings.Split(v, "=")[0]
				val := strings.Split(v, "=")[1]
				tagOpts.DesiredTags[key] = val
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			tags, err := aws.GetBucketTags(svc, tagOpts)
			if err != nil {
				return err
			}

			for _, v := range tags.TagSet {
				tagOpts.ActualTags[*v.Key] = *v.Value
			}

			logger.Info().Any("actualTags", tagOpts.ActualTags).Any("desiredTags", tagOpts.DesiredTags).Msg("")

			return nil
		},
	}
)
