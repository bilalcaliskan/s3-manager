package options

import (
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
)

type TagOptsKey struct{}

var tagOpts = &TagOptions{}

// TagOptions contains frequent command line and application options.
type TagOptions struct {
	// ActualState is state
	ActualTags map[string]string
	// TagsToAdd is state
	TagsToAdd map[string]string
	// TagsToRemove is state
	TagsToRemove map[string]string
	*options.RootOptions
}

// GetTagOptions returns the pointer of TagOptions
func GetTagOptions() *TagOptions {
	return tagOpts
}

func (opts *TagOptions) SetZeroValues() {
	opts.ActualTags = make(map[string]string)
	opts.TagsToAdd = make(map[string]string)
	opts.TagsToRemove = make(map[string]string)
}
