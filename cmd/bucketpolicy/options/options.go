package options

import (
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
)

type BucketPolicyOptsKey struct{}

var bucketPolicyOpts = &BucketPolicyOptions{}

type BucketPolicyOptions struct {
	BucketPolicyContent string
	*options.RootOptions
}

// GetBucketPolicyOptions returns the pointer of FindOptions
func GetBucketPolicyOptions() *BucketPolicyOptions {
	return bucketPolicyOpts
}

func (opts *BucketPolicyOptions) SetZeroValues() {
	opts.BucketPolicyContent = ""
}
