package list

import (
	"fmt"
	"github.com/bilalcaliskan/s3-manager/cmd/list/options"
	rootopts "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/aws"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/pkg/aws/types"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/utils"
	"github.com/rs/zerolog"

	"github.com/spf13/cobra"
)

func init() {
	listOpts = options.GetListOptions()
	listOpts.InitFlags(ListCmd)
}

var (
	svc      internalawstypes.S3ClientAPI
	logger   zerolog.Logger
	listOpts *options.ListOptions
	ListCmd  = &cobra.Command{
		Use:           "list",
		Short:         "lists the objects in the target bucket",
		SilenceUsage:  false,
		SilenceErrors: true,
		Example: `# show the current tagging configuration for bucket
s3-manager list --storage-class STANDARD --min-file-size-in-kb 128
		`,
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			var rootOpts *rootopts.RootOptions
			svc, rootOpts, logger, _ = utils.PrepareConstants(cmd)
			listOpts.RootOptions = rootOpts

			if err := utils.CheckArgs(args, 0); err != nil {
				logger.Error().Msg(err.Error())
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			logger.Info().Msg("listing objects")
			objects, err := aws.ListAllObjects(svc, listOpts.BucketName)
			if err != nil {
				logger.Error().
					Str("bucketName", listOpts.BucketName).
					Str("error", err.Error()).
					Msg("an error occurred while listing objects")
				return err
			}

			logger.Info().Msg("successfully listed objects")
			if len(objects) == 0 {
				logger.Warn().
					Str("bucketName", listOpts.BucketName).
					Msg("no objects found in the specified bucket")
				return nil
			}

			for _, obj := range objects {
				if listOpts.StorageClass != "" && string(obj.StorageClass) != listOpts.StorageClass {
					continue
				}
				if obj.Size == nil || obj.Key == nil {
					continue
				}
				if listOpts.MinFileSizeInKB > 0 && *obj.Size < int64(listOpts.MinFileSizeInKB)*1024 {
					continue
				}

				fmt.Printf("Name: %s, Size: %d kb, Storage Class: %s, Last Modified: %v\n", *obj.Key, *obj.Size/1024, obj.StorageClass, obj.LastModified)
			}

			return nil
		},
	}
)
