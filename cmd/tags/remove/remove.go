package remove

import (
	"errors"
	"fmt"
	"strings"

	utils2 "github.com/bilalcaliskan/s3-manager/cmd/tags/utils"

	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/utils"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/tags/options"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

const outputStr string = "%s=%s\n"

func init() {
	tagOpts = options.GetTagOptions()
}

var (
	svc       s3iface.S3API
	logger    zerolog.Logger
	tagOpts   *options.TagOptions
	RemoveCmd = &cobra.Command{
		Use:           "remove",
		Short:         "removes the tagging configuration for the target bucket",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			svc, tagOpts, logger = utils2.PrepareConstants(cmd, options.GetTagOptions())

			if err := utils2.CheckArgs(cmd, args); err != nil {
				logger.Error().Msg(err.Error())
				return err
			}

			tagOpts.ActualTags = make(map[string]string)
			tagOpts.TagsToRemove = make(map[string]string)

			tags, err := aws.GetBucketTags(svc, tagOpts)
			if err != nil {
				logger.Error().
					Str("error", err.Error()).
					Msg("an error occurred while fetching current tags")
				return err
			}

			logger.Info().Msg("fetched current bucket tags successfully")

			for _, v := range tags.TagSet {
				tagOpts.ActualTags[*v.Key] = *v.Value
			}

			tagArr := strings.Split(args[0], ",")
			for _, v := range tagArr {
				tag := strings.Split(v, "=")
				if len(tag) != 2 {
					err = errors.New("each key value pair for a tag should be separated with '='")
					logger.Error().
						Msg(err.Error())
					return err
				}

				if utils.HasKeyValuePair(tagOpts.ActualTags, tag[0], tag[1]) {
					tagOpts.TagsToRemove[tag[0]] = tag[1]
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(tagOpts.TagsToRemove) == 0 {
				logger.Warn().Msg("specified tags are not deletable, exiting")
				return nil
			}

			logger.Info().Msg("here is the list of current tags")
			for i, v := range tagOpts.ActualTags {
				fmt.Printf(outputStr, i, v)
			}

			logger.Info().Msg("will try to remove below tags")
			for i, v := range tagOpts.TagsToRemove {
				fmt.Printf(outputStr, i, v)
			}

			utils.RemoveMapElements(tagOpts.ActualTags, tagOpts.TagsToRemove)

			if _, err := aws.DeleteAllBucketTags(svc, tagOpts); err != nil {
				logger.Error().
					Str("error", err.Error()).
					Msg("an error occurred while deleting all the tags")
				return err
			}

			tagOpts.TagsToAdd = tagOpts.ActualTags
			if _, err := aws.SetBucketTags(svc, tagOpts); err != nil {
				logger.Error().
					Str("error", err.Error()).
					Msg("an error occurred while setting desired tags")
				return err
			}

			logger.Info().Msg("successfully removed target tags")
			logger.Info().Msg("here is the list of current tags again")
			for i, v := range tagOpts.ActualTags {
				fmt.Printf(outputStr, i, v)
			}

			return nil
		},
	}
)
