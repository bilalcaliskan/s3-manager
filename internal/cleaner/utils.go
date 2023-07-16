package cleaner

import (
	"sort"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/clean/options"
)

func sortObjects(slice []*s3.Object, opts *options.CleanOptions) {
	switch opts.SortBy {
	case "lastModificationDate":
		sort.Slice(slice, func(i, j int) bool {
			switch opts.Order {
			case "ascending":
				return slice[i].LastModified.Before(*slice[j].LastModified)
			case "descending":
				return slice[i].LastModified.After(*slice[j].LastModified)
			default:
				return slice[i].LastModified.Before(*slice[j].LastModified)
			}
		})
	case "size":
		sort.Slice(slice, func(i, j int) bool {
			switch opts.Order {
			case "ascending":
				return *slice[i].Size < *slice[j].Size
			case "descending":
				return *slice[i].Size > *slice[j].Size
			default:
				return *slice[i].Size < *slice[j].Size
			}
		})
	}
}
