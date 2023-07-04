package add

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	"github.com/bilalcaliskan/s3-manager/cmd/tags/utils"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/bilalcaliskan/s3-manager/cmd/tags/options"
	"github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	tagOpts = options.GetTagOptions()
}

var (
	svc           s3iface.S3API
	logger        zerolog.Logger
	confirmRunner prompt.PromptRunner = prompt.GetConfirmRunner()
	tagOpts       *options.TagOptions
	AddCmd        = &cobra.Command{
		Use:           "add",
		Short:         "adds the tagging configuration for the target bucket",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: `# add comma separated tagging configuration into bucket
s3-manager tags add foo1=bar1,foo2=bar2
		`,
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			svc, tagOpts, logger = utils.PrepareConstants(cmd, options.GetTagOptions())

			if err := utils.CheckArgs(cmd, args); err != nil {
				logger.Error().Msg(err.Error())
				return err
			}

			tagOpts.ActualTags = make(map[string]string)
			tagOpts.TagsToAdd = make(map[string]string)

			tags, err := aws.GetBucketTags(svc, tagOpts)
			if err != nil {
				logger.Error().
					Str("error", err.Error()).
					Msg("an error occurred while fetching current tags")
				return err
			}

			logger.Info().Msg("fetched current bucket tags successfully")

			for _, v := range tags.TagSet {
				tagOpts.TagsToAdd[*v.Key] = *v.Value
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

				tagOpts.TagsToAdd[tag[0]] = tag[1]
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			logger.Info().Msg("will try to set tags as below")
			for i, v := range tagOpts.TagsToAdd {
				fmt.Printf("%s=%s\n", i, v)
			}

			if tagOpts.DryRun {
				logger.Info().Msg("skipping operation since '--dry-run' flag is passed")
				return nil
			}

			if !tagOpts.AutoApprove {
				var res string
				if res, err = confirmRunner.Run(); err != nil {
					return err
				}

				if strings.ToLower(res) == "n" {
					return errors.New("user terminated the process")
				}
			}

			if _, err := aws.SetBucketTags(svc, tagOpts); err != nil {
				logger.Error().
					Str("error", err.Error()).
					Msg("an error occurred while setting tags")
				return err
			}

			logger.Info().Msg("set bucket tags successfully")

			return nil
		},
	}
)
