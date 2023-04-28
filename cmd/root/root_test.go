package root

import (
	"errors"
	"testing"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	"github.com/spf13/cobra"

	"github.com/stretchr/testify/assert"
)

/*type promptMock struct {
	msg string
	err error
}

func (p promptMock) Run() (string, error) {
	// return expected result
	return p.msg, p.err
}
*/

type selectMock struct {
	msg string
	err error
}

func (p selectMock) Run() (int, string, error) {
	// return expected result
	return 1, p.msg, p.err
}

func TestOuterExecute(t *testing.T) {
	err := setAccessFlags(rootCmd, "", "", "", "")
	assert.Nil(t, err)
}

func TestExecute(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("verbose", "true")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("bannerFilePath", "./../../build/ci/banner.txt")
	assert.Nil(t, err)

	err = rootCmd.Execute()
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

// 45.5%

/*func TestExecuteWithPromptsSuccessSelectFailPrompt(t *testing.T) {
	// get original value for valid domains
	predefinedValidDomainsOrg := mail.PredefinedValidDomains

	// override valid domains
	mail.PredefinedValidDomains = []string{"ssss.com"}

	selectRunner = selectMock{msg: "Yes please!", err: nil}
	promptRunner = promptMock{msg: "nonexistedemailaddress@example.com", err: errors.New("dummy error")}
	err := rootCmd.Execute()
	assert.NotNil(t, err)

	// revert valid domains
	mail.PredefinedValidDomains = predefinedValidDomainsOrg
	selectRunner = prompt.GetSelectRunner()
	promptRunner = prompt.GetPromptRunner()
}*/

/*func TestExecuteInteractive(t *testing.T) {
	err := setAccessFlags(rootCmd, "thisisaccesskey", "thisissecretkey", "thisisbucketname", "thisisregion")
	assert.Nil(t, err)

	err = rootCmd.PersistentFlags().Set("interactive", "true")
	assert.Nil(t, err)

	selectRunner = selectMock{msg: "search", err: nil}

	err = rootCmd.Execute()
	assert.NotNil(t, err)

	opts.SetZeroValues()
	selectRunner = prompt.GetSelectRunner("Select operation", []string{"search", "clean"})
}*/

func TestExecuteInteractiveSelectRunnerErr(t *testing.T) {
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
