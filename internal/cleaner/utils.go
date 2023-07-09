package cleaner

import (
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/clean/options"
	"github.com/rs/zerolog"
)

func getProperObjects(cleanOpts *options.CleanOptions, allFiles *s3.ListObjectsOutput, logger zerolog.Logger) (res []*s3.Object) {
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

func makeDecisionBySize(opts *options.CleanOptions, res []*s3.Object, object *s3.Object) []*s3.Object {
	if (opts.MinFileSizeInMb == 0 && opts.MaxFileSizeInMb != 0) && *object.Size < opts.MaxFileSizeInMb*1000000 {
		res = append(res, object)
	} else if (opts.MinFileSizeInMb != 0 && opts.MaxFileSizeInMb == 0) && *object.Size >= opts.MinFileSizeInMb*1000000 {
		res = append(res, object)
	} else if opts.MinFileSizeInMb == 0 && opts.MaxFileSizeInMb == 0 {
		res = append(res, object)
	} else if opts.MinFileSizeInMb != 0 && opts.MaxFileSizeInMb != 0 && (*object.Size >= opts.MinFileSizeInMb*1000000 && *object.Size < opts.MaxFileSizeInMb*1000000) {
		res = append(res, object)
	}

	return res
}

func sortObjects(slice []*s3.Object, opts *options.CleanOptions) {
	switch opts.SortBy {
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

func arrayContains(sl []string, name string) bool {
	for _, value := range sl {
		if strings.Contains(name, value) {
			return true
		}
	}

	return false
}
