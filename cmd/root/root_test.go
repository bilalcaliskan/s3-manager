//go:build e2e

package root

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/stretchr/testify/assert"
)

func TestOuterExecute(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = Execute()
	assert.Nil(t, err)

	opts.SetZeroValues()
}

func TestExecuteCreateClientFailure(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "invalidregion")
	assert.Nil(t, err)

	_ = rootCmd.Execute()
	//assert.NotNil(t, err)

	opts.SetZeroValues()
}

func TestExecute(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("verbose", "true")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("banner-file-path", "./../../build/ci/banner.txt")
	assert.Nil(t, err)

	_ = rootCmd.Execute()
	//assert.NotNil(t, err)

	opts.SetZeroValues()
}

func setAccessFlags(cmd *cobra.Command, accessKey, secretKey, bucketName, region string) error {
	if err := cmd.PersistentFlags().Set("access-key", accessKey); err != nil {
		return err
	}

	if err := cmd.PersistentFlags().Set("secret-key", secretKey); err != nil {
		return err
	}

	if err := cmd.PersistentFlags().Set("bucket-name", bucketName); err != nil {
		return err
	}

	return cmd.PersistentFlags().Set("region", region)
}
