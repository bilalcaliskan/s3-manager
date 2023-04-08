package cleaner

import (
	"errors"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	start "github.com/bilalcaliskan/s3-manager/cmd/clean/options"
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog"
)

func getProperObjects(cleanOpts *start.CleanOptions, allFiles *s3.ListObjectsOutput, logger zerolog.Logger) (res []*s3.Object) {
	extensions := strings.Split(cleanOpts.FileExtensions, ",")

	for _, v := range allFiles.Contents {
		if strings.HasSuffix(*v.Key, "/") {
			logger.Debug().Str("key", *v.Key).Msg("object has directory suffix, skipping that one")
			continue
		}

		if len(extensions) > 0 && !arrayContains(extensions, *v.Key) {
			continue
		}

		res = makeDecisionBySize(cleanOpts, res, v)
	}

	return res
}

func makeDecisionBySize(startOpts *start.CleanOptions, res []*s3.Object, object *s3.Object) []*s3.Object {
	if (startOpts.MinFileSizeInMb == 0 && startOpts.MaxFileSizeInMb != 0) && *object.Size < startOpts.MaxFileSizeInMb*1000000 {
		res = append(res, object)
	} else if (startOpts.MinFileSizeInMb != 0 && startOpts.MaxFileSizeInMb == 0) && *object.Size >= startOpts.MinFileSizeInMb*1000000 {
		res = append(res, object)
	} else if startOpts.MinFileSizeInMb == 0 && startOpts.MaxFileSizeInMb == 0 {
		res = append(res, object)
	} else if startOpts.MinFileSizeInMb != 0 && startOpts.MaxFileSizeInMb != 0 && (*object.Size >= startOpts.MinFileSizeInMb*1000000 && *object.Size < startOpts.MaxFileSizeInMb*1000000) {
		res = append(res, object)
	}

	return res
}

func sortObjects(slice []*s3.Object, startOpts *start.CleanOptions) {
	switch startOpts.SortBy {
	case "lastModificationDate":
		sort.Slice(slice, func(i, j int) bool {
			return slice[i].LastModified.Before(*slice[j].LastModified)
		})
	case "size":
		sort.Slice(slice, func(i, j int) bool {
			return *slice[i].Size < *slice[j].Size
		})
	}
}

func checkLength(targetObjects []*s3.Object) error {
	if len(targetObjects) == 0 {
		return errors.New("no deletable file found on the target bucket")
	}

	return nil
}

func promptDeletion(startOpts *start.CleanOptions, logger zerolog.Logger, keys []string) error {
	if !startOpts.AutoApprove {
		logger.Info().Any("files", keys).Msg("these files will be removed if you approve:")

		prompt := promptui.Prompt{
			Label:     "Delete Files? (y/N)",
			IsConfirm: true,
			Validate: func(s string) error {
				if len(s) == 1 {
					return nil
				}

				return errors.New("invalid input")
			},
		}

		if res, err := prompt.Run(); err != nil {
			if strings.ToLower(res) == "n" {
				return errors.New("user terminated the process")
			}

			return errors.New("invalid input")
		}
	}

	return nil
}

func arrayContains(sl []string, name string) bool {
	for _, value := range sl {
		if strings.Contains(name, value) {
			return true
		}
	}

	return false
}
