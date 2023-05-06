package options

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootOptions = &RootOptions{}
)

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
	// BannerFilePath is the relative path to the banner file
	BannerFilePath string
	Timeout        int64
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
	cmd.PersistentFlags().StringVarP(&opts.BannerFilePath, "bannerFilePath", "", "banner.txt",
		"relative path of the banner file")
	cmd.PersistentFlags().Int64VarP(&opts.Timeout, "timeout", "t", 10,
		"")
}

func (opts *RootOptions) SetAccessFlagsRequired(cmd *cobra.Command) {
	if opts.AccessKey == "" {
		_ = cmd.MarkPersistentFlagRequired("accessKey")
	}

	if opts.SecretKey == "" {
		_ = cmd.MarkPersistentFlagRequired("secretKey")
	}

	if opts.BucketName == "" {
		_ = cmd.MarkPersistentFlagRequired("bucketName")
	}

	if opts.Region == "" {
		_ = cmd.MarkPersistentFlagRequired("region")
	}
}

func (opts *RootOptions) SetAccessCredentialsFromEnv() error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("aws")
	if err := viper.BindEnv("access_key", "secret_key", "bucket_name", "region"); err != nil {
		return err
	}

	if accessKey := viper.Get("access_key"); accessKey != nil {
		opts.AccessKey = fmt.Sprintf("%v", accessKey)
	}

	if secretKey := viper.Get("secret_key"); secretKey != nil {
		opts.SecretKey = fmt.Sprintf("%v", secretKey)
	}

	if bucketName := viper.Get("bucket_name"); bucketName != nil {
		opts.BucketName = fmt.Sprintf("%v", bucketName)
	}

	if region := viper.Get("region"); region != nil {
		opts.Region = fmt.Sprintf("%v", region)
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
	opts.BannerFilePath = "banner.txt"
}
