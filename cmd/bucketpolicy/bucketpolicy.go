package bucketpolicy

import (
	"github.com/bilalcaliskan/s3-manager/cmd/bucketpolicy/add"
	"github.com/bilalcaliskan/s3-manager/cmd/bucketpolicy/remove"
	"github.com/bilalcaliskan/s3-manager/cmd/bucketpolicy/show"

	"github.com/spf13/cobra"
)

func init() {
	BucketPolicyCmd.AddCommand(show.ShowCmd)
	BucketPolicyCmd.AddCommand(add.AddCmd)
	BucketPolicyCmd.AddCommand(remove.RemoveCmd)
}

var (
	BucketPolicyCmd = &cobra.Command{
		Use:           "bucketpolicy",
		Short:         "shows/sets the bucket policy configuration of the target bucket",
		SilenceUsage:  false,
		SilenceErrors: false,
	}
)
