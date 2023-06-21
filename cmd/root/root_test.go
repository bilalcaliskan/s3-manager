//go:build e2e
// +build e2e

package root

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/stretchr/testify/assert"
)

// TODO: uncomment when interactivity enabled again
/*type promptMock struct {
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
}*/

func TestOuterExecute(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	// TODO: uncomment when interactivity enabled again
	/*err = rootCmd.PersistentFlags().Set("interactive", "false")
	assert.Nil(t, err)*/

	err = Execute()
	assert.Nil(t, err)

	opts.SetZeroValues()
}

func TestExecuteCreateSvcFailure(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "")
	assert.Nil(t, err)

	// TODO: uncomment when interactivity enabled again
	/*err = rootCmd.PersistentFlags().Set("interactive", "false")
	assert.Nil(t, err)*/

	err = rootCmd.Execute()
	assert.NotNil(t, err)

	opts.SetZeroValues()
}

func TestExecute(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("verbose", "true")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("bannerFilePath", "./../../build/ci/banner.txt")
	assert.Nil(t, err)

	// TODO: uncomment when interactivity enabled again
	/*err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)*/

	err = rootCmd.Execute()
	// TODO: uncomment when interactivity enabled again
	//assert.NotNil(t, err)
	// TODO: remove when interactivity enabled again
	assert.Nil(t, err)

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

	return cmd.PersistentFlags().Set("region", region)
}

// TODO: uncomment when interactivity enabled again
/*func TestExecuteInteractiveSelectRunnerSearchSuccess(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	options.SelectRunner = selectMock{msg: "search", err: nil}

	err = rootCmd.Execute()
	assert.NotNil(t, err)

	opts.SetZeroValues()
	options.SelectRunner = prompt.GetSelectRunner("Select operation", []string{"search", "clean"})
}*/

// TODO: uncomment when interactivity enabled again
/*func TestExecuteInteractiveSelectRunnerSearchErr(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	options.SelectRunner = selectMock{msg: "search", err: errors.New("dummy error")}

	err = rootCmd.Execute()
	assert.NotNil(t, err)

	opts.SetZeroValues()
	options.SelectRunner = prompt.GetSelectRunner("Select operation", []string{"search", "clean"})
}*/

// TODO: uncomment when interactivity enabled again
/*func TestExecuteInteractiveSelectRunnerCleanSuccess(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	options.SelectRunner = selectMock{msg: "clean", err: nil}

	err = rootCmd.Execute()
	assert.NotNil(t, err)

	opts.SetZeroValues()
	options.SelectRunner = prompt.GetSelectRunner("Select operation", []string{"search", "clean"})
}*/

// TODO: uncomment when interactivity enabled again
/*func TestExecuteInteractiveSelectRunnerErr(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	options.SelectRunner = selectMock{msg: "", err: errors.New("dummy error")}

	err = rootCmd.Execute()
	assert.NotNil(t, err)

	opts.SetZeroValues()
	options.SelectRunner = prompt.GetSelectRunner("Select operation", []string{"search", "clean"})
}*/

// TODO: uncomment when interactivity enabled again
/*func TestExecuteInteractiveAccessPromptErr(t *testing.T) {
	err := setAccessFlags(rootCmd, "", "thisissecretkey", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	accessKeyRunnerOrg := options.AccessKeyRunner
	options.AccessKeyRunner = promptMock{
		msg: "",
		err: errors.New("new dummy error"),
	}

	err = rootCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, opts.AccessKey, "")

	opts.SetZeroValues()
	options.AccessKeyRunner = accessKeyRunnerOrg
}*/

// TODO: uncomment when interactivity enabled again
/*func TestExecuteInteractiveAccessPromptSuccess(t *testing.T) {
	err := setAccessFlags(rootCmd, "", "thisissecretkey", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	accessKeyRunnerOrg := options.AccessKeyRunner
	options.AccessKeyRunner = promptMock{
		msg: "thisisaccesskey",
		err: nil,
	}

	err = rootCmd.Execute()
	assert.NotNil(t, err)
	fmt.Println(opts)
	assert.Equal(t, opts.AccessKey, "thisisaccesskey")

	opts.SetZeroValues()
	options.AccessKeyRunner = accessKeyRunnerOrg
}*/

// TODO: uncomment when interactivity enabled again
/*func TestExecuteInteractiveSecretPromptErr(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	secretKeyRunnerOrg := options.SecretKeyRunner
	options.SecretKeyRunner = promptMock{
		msg: "",
		err: errors.New("new dummy error"),
	}

	err = rootCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, opts.SecretKey, "")

	opts.SetZeroValues()
	options.SecretKeyRunner = secretKeyRunnerOrg
}*/

// TODO: uncomment when interactivity enabled again
/*func TestExecuteInteractiveSecretPromptSuccess(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	secretKeyRunnerOrg := options.SecretKeyRunner
	options.SecretKeyRunner = promptMock{
		msg: "thisissecretkey",
		err: nil,
	}

	err = rootCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, opts.SecretKey, "thisissecretkey")

	opts.SetZeroValues()
	options.SecretKeyRunner = secretKeyRunnerOrg
}*/

// TODO: uncomment when interactivity enabled again
/*func TestExecuteInteractiveBucketPromptErr(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	bucketRunnerOrg := options.BucketRunner
	options.BucketRunner = promptMock{
		msg: "",
		err: errors.New("new dummy error"),
	}

	err = rootCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, opts.BucketName, "")

	opts.SetZeroValues()
	options.BucketRunner = bucketRunnerOrg
}*/

// TODO: uncomment when interactivity enabled again
/*func TestExecuteInteractiveBucketPromptSuccess(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	bucketRunnerOrg := options.BucketRunner
	options.BucketRunner = promptMock{
		msg: "thisisbucketname",
		err: nil,
	}

	err = rootCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, opts.BucketName, "thisisbucketname")

	opts.SetZeroValues()
	options.BucketRunner = bucketRunnerOrg
}*/

// TODO: uncomment when interactivity enabled again
/*func TestExecuteInteractiveRegionPromptErr(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	regionRunnerOrg := options.RegionRunner
	options.RegionRunner = promptMock{
		msg: "",
		err: errors.New("new dummy error"),
	}

	err = rootCmd.Execute()
	assert.NotNil(t, err)
	assert.Equal(t, opts.Region, "")

	opts.SetZeroValues()
	options.RegionRunner = regionRunnerOrg
}*/

// TODO: uncomment when interactivity enabled again
/*func TestExecuteInteractiveRegionPromptSuccess(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	regionRunnerOrg := options.RegionRunner
	options.RegionRunner = promptMock{
		msg: "thisisregion",
		err: nil,
	}

	err = rootCmd.Execute()
	assert.NotNil(t, err)
	//assert.Equal(t, opts.Region, "thisisregion")

	opts.SetZeroValues()
	options.RegionRunner = regionRunnerOrg
}*/
