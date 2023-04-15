package options

import (
	"fmt"
	"github.com/manifoldco/promptui"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootOptions = &RootOptions{}

type (
	OptsKey   struct{}
	LoggerKey struct{}
	S3SvcKey  struct{}
)

// RootOptions contains frequent command line and application options.
type RootOptions struct {
	// AccessKey is the access key credentials for accessing AWS over client
	AccessKey string
	// SecretKey is the secret key credentials for accessing AWS over client
	SecretKey string
	// BucketName is the name of target bucket
	BucketName string
	// Region is the region of the target bucket
	Region string
	// VerboseLog is the verbosity of the logging library
	VerboseLog bool
	// Interactive is the decision of that if you want to use interactive feature
	Interactive bool
}

func (opts *RootOptions) InitFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&opts.BucketName, "bucketName", "", "", "name of "+
		"the target bucket on S3, this value also can be passed via \"AWS_BUCKET_NAME\" environment variable (default \"\")")
	cmd.PersistentFlags().StringVarP(&opts.AccessKey, "accessKey", "", "",
		"access key credential to access S3 bucket, this value also can be passed via \"AWS_ACCESS_KEY\" "+
			"environment variable (default \"\")")
	cmd.PersistentFlags().StringVarP(&opts.SecretKey, "secretKey", "", "",
		"secret key credential to access S3 bucket, this value also can be passed via \"AWS_SECRET_KEY\" "+
			"environment variable (default \"\")")
	cmd.PersistentFlags().StringVarP(&opts.Region, "region", "", "",
		"region of the target bucket on S3, this value also can be passed via \"AWS_REGION\" environment "+
			"variable (default \"\")")
	cmd.PersistentFlags().BoolVarP(&opts.VerboseLog, "verbose", "", false,
		"verbose output of the logging library (default false)")
	cmd.PersistentFlags().BoolVarP(&opts.Interactive, "interactive", "i", false,
		"decision of that if you want to use interactive feature (default false)")
}

func (opts *RootOptions) SetAccessCredentialsFromEnv(cmd *cobra.Command) error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("aws")
	if err := viper.BindEnv("access_key", "secret_key", "bucket_name", "region"); err != nil {
		return err
	}

	if accessKey := viper.Get("access_key"); accessKey != nil {
		opts.AccessKey = fmt.Sprintf("%v", accessKey)
	}
	//else {
	//	_ = cmd.MarkPersistentFlagRequired("accessKey")
	//}

	if secretKey := viper.Get("secret_key"); secretKey != nil {
		opts.SecretKey = fmt.Sprintf("%v", secretKey)
	}
	//else {
	//	_ = cmd.MarkPersistentFlagRequired("secretKey")
	//}

	if bucketName := viper.Get("bucket_name"); bucketName != nil {
		opts.BucketName = fmt.Sprintf("%v", bucketName)
	}
	//else {
	//	_ = cmd.MarkPersistentFlagRequired("bucketName")
	//}

	if region := viper.Get("region"); region != nil {
		opts.Region = fmt.Sprintf("%v", region)
	}
	//else {
	//	_ = cmd.MarkPersistentFlagRequired("region")
	//}

	return nil
}

func (opts *RootOptions) PromptAccessCredentials() error {
	if opts.AccessKey == "" {
		accessKeyPrompt := promptui.Prompt{
			Label: "Provide AWS Access Key",
		}

		result, err := accessKeyPrompt.Run()

		if err != nil {
			return err
		}
		opts.AccessKey = result
	}

	if opts.SecretKey == "" {
		secretKeyPrompt := promptui.Prompt{
			Label: "Provide AWS Secret Key",
		}

		result, err := secretKeyPrompt.Run()

		if err != nil {
			return err
		}
		opts.SecretKey = result
	}

	if opts.Region == "" {
		regionPrompt := promptui.Prompt{
			Label: "Provide AWS Region",
		}

		result, err := regionPrompt.Run()

		if err != nil {
			return err
		}
		opts.Region = result
	}

	if opts.BucketName == "" {
		bucketPrompt := promptui.Prompt{
			Label: "Provide AWS Bucket Name",
		}

		result, err := bucketPrompt.Run()

		if err != nil {
			return err
		}
		opts.BucketName = result
	}

	return nil
}

// GetRootOptions returns the pointer of S3CleanerOptions
func GetRootOptions() *RootOptions {
	return rootOptions
}

func (opts *RootOptions) SetZeroValues() {
	opts.BucketName = ""
	opts.AccessKey = ""
	opts.SecretKey = ""
	opts.Region = ""
	opts.VerboseLog = false
}
