package prompt

import (
	"fmt"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog"
)

type SelectRunner interface {
	Run() (int, string, error)
}

func GetSelectRunner(label string, items []string) *promptui.Select {
	return &promptui.Select{
		Label: label,
		Items: items,
	}
}

type PromptRunner interface {
	Run() (string, error)
}

func GetPromptRunner(label string, valFunc func(s string) error) *promptui.Prompt {
	return &promptui.Prompt{
		Label:    label,
		Validate: valFunc,
	}
}

func PromptAccessCreds(opts *options.RootOptions, logger zerolog.Logger) error {
	infoLog := "skipping %s prompt since it is provided either by environment variable or flag"
	var promptRunner *promptui.Prompt

	if opts.AccessKey == "" {
		promptRunner = GetPromptRunner("Provide AWS Access Key", nil)
		result, err := promptRunner.Run()

		if err != nil {
			return err
		}
		opts.AccessKey = result
	} else {
		logger.Info().Msg(fmt.Sprintf(infoLog, "accessKey"))
	}

	if opts.SecretKey == "" {
		promptRunner = GetPromptRunner("Provide AWS Secret Key", nil)
		result, err := promptRunner.Run()

		if err != nil {
			return err
		}
		opts.SecretKey = result
	} else {
		logger.Info().Msg(fmt.Sprintf(infoLog, "secretKey"))
	}

	if opts.Region == "" {
		promptRunner = GetPromptRunner("Provide AWS Region", nil)
		result, err := promptRunner.Run()

		if err != nil {
			return err
		}
		opts.Region = result
	} else {
		logger.Info().Msg(fmt.Sprintf(infoLog, "region"))
	}

	if opts.BucketName == "" {
		promptRunner = GetPromptRunner("Provide AWS Bucket Name", nil)
		result, err := promptRunner.Run()

		if err != nil {
			return err
		}
		opts.BucketName = result
	} else {
		logger.Info().Msg(fmt.Sprintf(infoLog, "bucketName"))
	}

	return nil
}
