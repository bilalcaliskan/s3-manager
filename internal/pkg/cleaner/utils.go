package cleaner

import (
	"sort"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/bilalcaliskan/s3-manager/cmd/clean/options"
)

// sortObjects sorts a slice of *s3.Object based on the specified sorting criteria in the CleanOptions.
//
// The sorting criteria can be "lastModificationDate" or "size".
// For "lastModificationDate", the sorting order can be "ascending" or "descending".
// For "size", the sorting order can be "ascending" or "descending".
//
// Parameters:
// - slice: The slice of *s3.Object to be sorted.
// - opts: The CleanOptions struct containing the sorting criteria and order.
//
// Returns:
// - None
func sortObjects(slice []types.Object, opts *options.CleanOptions) {
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
