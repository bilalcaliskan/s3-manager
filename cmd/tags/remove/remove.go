package remove

import (
	"errors"
	"fmt"
	"strings"

	v2s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	"github.com/bilalcaliskan/s3-manager/internal/utils"

	"github.com/bilalcaliskan/s3-manager/cmd/tags/options"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

const outputStr string = "%s=%s\n"

func init() {
	tagOpts = options.GetTagOptions()
}

var (
	svc           *v2s3.Client
	logger        zerolog.Logger
	confirmRunner prompt.PromptRunner
	tagOpts       *options.TagOptions
	RemoveCmd     = &cobra.Command{
		Use:           "remove",
		Short:         "removes the tagging configuration for the target bucket",
		SilenceUsage:  false,
		SilenceErrors: true,
		Example: `# remove comma separated tagging configuration from bucket
s3-manager tags remove foo1=bar1,foo2=bar2
		`,
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			var rootOpts *rootopts.RootOptions
			svc, rootOpts, logger, confirmRunner = utils.PrepareConstants(cmd)
			tagOpts.RootOptions = rootOpts

			if err := utils.CheckArgs(args, 1); err != nil {
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

			if _, err := aws.DeleteAllBucketTags(svc, tagOpts, confirmRunner, logger); err != nil {
				logger.Error().
					Str("error", err.Error()).
					Msg("an error occurred while deleting all the tags")
				return err
			}

			tagOpts.TagsToAdd = tagOpts.ActualTags
			if err := aws.SetBucketTags(svc, tagOpts, confirmRunner, logger); err != nil {
				logger.Error().
					Str("error", err.Error()).
					Msg("an error occurred while setting desired tags")
				return err
			}

			logger.Info().Msg("successfully removed target tags")
			logger.Info().Msg("here is the list of current tags")
			for i, v := range tagOpts.ActualTags {
				fmt.Printf(outputStr, i, v)
			}

			return nil
		},
	}
)
