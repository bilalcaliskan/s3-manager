package root

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	"github.com/spf13/cobra"

	"github.com/stretchr/testify/assert"
)

type promptMock struct {
	msg string
	err error
}

func (p promptMock) Run() (string, error) {
	// return expected result
	return p.msg, p.err
}

type selectMock struct {
	msg string
	err error
}

func (p selectMock) Run() (int, string, error) {
	// return expected result
	return 1, p.msg, p.err
}

func TestOuterExecute(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "false")
	assert.Nil(t, err)

	err = Execute()
	assert.Nil(t, err)

	opts.SetZeroValues()
}

func TestExecute(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("verbose", "true")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("bannerFilePath", "./../../build/ci/banner.txt")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	err = rootCmd.Execute()
	assert.NotNil(t, err)

	opts.SetZeroValues()
}

func setAccessFlags(cmd *cobra.Command, accessKey, secretKey, bucketName, region string) error {
	if err := cmd.PersistentFlags().Set("accessKey", accessKey); err != nil {
		return err
	}

	if err := cmd.PersistentFlags().Set("secretKey", secretKey); err != nil {
		return err
	}

	if err := cmd.PersistentFlags().Set("bucketName", bucketName); err != nil {
		return err
	}

	if err := cmd.PersistentFlags().Set("region", region); err != nil {
		return err
	}

	return nil
}

func TestExecuteInteractiveSelectRunnerSearchSuccess(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	selectRunner = selectMock{msg: "search", err: nil}

	err = rootCmd.Execute()
	assert.NotNil(t, err)

	opts.SetZeroValues()
	selectRunner = prompt.GetSelectRunner("Select operation", []string{"search", "clean"})
}

func TestExecuteInteractiveSelectRunnerSearchErr(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	selectRunner = selectMock{msg: "search", err: errors.New("dummy error")}

	err = rootCmd.Execute()
	assert.NotNil(t, err)

	opts.SetZeroValues()
	selectRunner = prompt.GetSelectRunner("Select operation", []string{"search", "clean"})
}

func TestExecuteInteractiveSelectRunnerCleanSuccess(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	selectRunner = selectMock{msg: "clean", err: nil}

	err = rootCmd.Execute()
	assert.NotNil(t, err)

	opts.SetZeroValues()
	selectRunner = prompt.GetSelectRunner("Select operation", []string{"search", "clean"})
}

func TestExecuteInteractiveSelectRunnerErr(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	selectRunner = selectMock{msg: "", err: errors.New("dummy error")}

	err = rootCmd.Execute()
	assert.NotNil(t, err)

	opts.SetZeroValues()
	selectRunner = prompt.GetSelectRunner("Select operation", []string{"search", "clean"})
}

func TestExecuteInteractiveAccessPromptErr(t *testing.T) {
	err := setAccessFlags(rootCmd, "", "thisissecretkey", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	accessKeyRunnerOrg := accessKeyRunner
	accessKeyRunner = promptMock{
		msg: "",
		err: errors.New("new dummy error"),
	}

	err = rootCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, opts.AccessKey, "")

	opts.SetZeroValues()
	accessKeyRunner = accessKeyRunnerOrg
}

func TestExecuteInteractiveAccessPromptSuccess(t *testing.T) {
	err := setAccessFlags(rootCmd, "", "thisissecretkey", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	accessKeyRunnerOrg := accessKeyRunner
	accessKeyRunner = promptMock{
		msg: "thisisaccesskey",
		err: nil,
	}

	err = rootCmd.Execute()
	assert.NotNil(t, err)
	fmt.Println(opts)
	assert.Equal(t, opts.AccessKey, "thisisaccesskey")

	opts.SetZeroValues()
	accessKeyRunner = accessKeyRunnerOrg
}

func TestExecuteInteractiveSecretPromptErr(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	secretKeyRunnerOrg := secretKeyRunner
	secretKeyRunner = promptMock{
		msg: "",
		err: errors.New("new dummy error"),
	}

	err = rootCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, opts.SecretKey, "")

	opts.SetZeroValues()
	secretKeyRunner = secretKeyRunnerOrg
}

func TestExecuteInteractiveSecretPromptSuccess(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	secretKeyRunnerOrg := secretKeyRunner
	secretKeyRunner = promptMock{
		msg: "thisissecretkey",
		err: nil,
	}

	err = rootCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, opts.SecretKey, "thisissecretkey")

	opts.SetZeroValues()
	secretKeyRunner = secretKeyRunnerOrg
}

func TestExecuteInteractiveBucketPromptErr(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	bucketRunnerOrg := bucketRunner
	bucketRunner = promptMock{
		msg: "",
		err: errors.New("new dummy error"),
	}

	err = rootCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, opts.BucketName, "")

	opts.SetZeroValues()
	bucketRunner = bucketRunnerOrg
}

func TestExecuteInteractiveBucketPromptSuccess(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	bucketRunnerOrg := bucketRunner
	bucketRunner = promptMock{
		msg: "thisisbucketname",
		err: nil,
	}

	err = rootCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, opts.BucketName, "thisisbucketname")

	opts.SetZeroValues()
	bucketRunner = bucketRunnerOrg
}

func TestExecuteInteractiveRegionPromptErr(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	regionRunnerOrg := regionRunner
	regionRunner = promptMock{
		msg: "",
		err: errors.New("new dummy error"),
	}

	err = rootCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, opts.Region, "")

	opts.SetZeroValues()
	regionRunner = regionRunnerOrg
}

func TestExecuteInteractiveRegionPromptSuccess(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	regionRunnerOrg := regionRunner
	regionRunner = promptMock{
		msg: "thisisregion",
		err: nil,
	}

	err = rootCmd.Execute()
	assert.NotNil(t, err)
	//assert.Equal(t, opts.Region, "thisisregion")

	opts.SetZeroValues()
	regionRunner = regionRunnerOrg
}
