package add

import (
	v2s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/internal/utils"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	"github.com/bilalcaliskan/s3-manager/cmd/tags/options"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	tagOpts = options.GetTagOptions()
}

var (
	svc           *v2s3.Client
	logger        zerolog.Logger
	confirmRunner prompt.PromptRunner
	tagOpts       *options.TagOptions
	AddCmd        = &cobra.Command{
		Use:           "add",
		Short:         "adds the tagging configuration for the target bucket",
		SilenceUsage:  false,
		SilenceErrors: true,
		Example: `# add comma separated tagging configuration into bucket
s3-manager tags add foo1=bar1,foo2=bar2
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
			tagOpts.TagsToAdd = make(map[string]string)

			//tags, err := aws.GetBucketTags(svc, tagOpts)
			//if err != nil {
			//	logger.Error().
			//		Str("error", err.Error()).
			//		Msg("an error occurred while fetching current tags")
			//	return err
			//}
			//
			//logger.Info().Msg("fetched current bucket tags successfully")
			//
			//for _, v := range tags.TagSet {
			//	tagOpts.TagsToAdd[*v.Key] = *v.Value
			//}
			//
			//tagArr := strings.Split(args[0], ",")
			//for _, v := range tagArr {
			//	tag := strings.Split(v, "=")
			//	if len(tag) != 2 {
			//		err = errors.New("each key value pair for a tag should be separated with '='")
			//		logger.Error().
			//			Msg(err.Error())
			//		return err
			//	}
			//
			//	tagOpts.TagsToAdd[tag[0]] = tag[1]
			//}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			//logger.Info().Msg("will try to set tags as below")
			//for i, v := range tagOpts.TagsToAdd {
			//	fmt.Printf("%s=%s\n", i, v)
			//}
			//
			//if err := aws.SetBucketTags(svc, tagOpts, confirmRunner, logger); err != nil {
			//	logger.Error().
			//		Str("error", err.Error()).
			//		Msg("an error occurred while setting tags")
			//	return err
			//}
			//
			//logger.Info().Msg("set bucket tags successfully")

			return nil
		},
	}
)
