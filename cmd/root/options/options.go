package options

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootOptions = &RootOptions{}

type (
	OptsKey          struct{}
	LoggerKey        struct{}
	S3SvcKey         struct{}
	ConfirmRunnerKey struct{}
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
	// BannerFilePath is the relative path to the banner file
	BannerFilePath string
	// AutoApprove is the boolean flag that lets you bypass approval before non read only operations
	AutoApprove bool
	// DryRun is the boolean flag that lets you see what will be changed before non read only operations
	DryRun bool
}

func (opts *RootOptions) InitFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&opts.BucketName, "bucket-name", "", "", "name of "+
		"the target bucket on S3, this value also can be passed via \"AWS_BUCKET_NAME\" environment variable (default \"\")")
	cmd.PersistentFlags().StringVarP(&opts.AccessKey, "access-key", "", "",
		"access key credential to access S3 bucket, this value also can be passed via \"AWS_ACCESS_KEY\" "+
			"environment variable (default \"\")")
	cmd.PersistentFlags().StringVarP(&opts.SecretKey, "secret-key", "", "",
		"secret key credential to access S3 bucket, this value also can be passed via \"AWS_SECRET_KEY\" "+
			"environment variable (default \"\")")
	cmd.PersistentFlags().StringVarP(&opts.Region, "region", "", "",
		"region of the target bucket on S3, this value also can be passed via \"AWS_REGION\" environment "+
			"variable (default \"\")")
	cmd.PersistentFlags().BoolVarP(&opts.VerboseLog, "verbose", "", false,
		"verbose output of the logging library (default false)")
	cmd.PersistentFlags().StringVarP(&opts.BannerFilePath, "banner-file-path", "", "banner.txt",
		"relative path of the banner file")
	cmd.PersistentFlags().BoolVarP(&opts.AutoApprove, "auto-approve", "", false, "boolean flag "+
		"that lets you bypass approval before non read only operations")
	cmd.PersistentFlags().BoolVarP(&opts.DryRun, "dry-run", "", false, "boolean flag that lets "+
		"you see what will be changed before non read only operations")
}

func (opts *RootOptions) SetAccessFlagsRequired(cmd *cobra.Command) {
	if opts.AccessKey == "" {
		_ = cmd.MarkPersistentFlagRequired("access-key")
	}

	if opts.SecretKey == "" {
		_ = cmd.MarkPersistentFlagRequired("secret-key")
	}

	if opts.BucketName == "" {
		_ = cmd.MarkPersistentFlagRequired("bucket-name")
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

func GetRootOptions() *RootOptions {
	return rootOptions
}

func GetMockedRootOptions() *RootOptions {
	return &RootOptions{
		AccessKey:  "thisisaccesskey",
		SecretKey:  "thisissecretkey",
		Region:     "thisisregion",
		BucketName: "thisisbucketname",
	}
}

func (opts *RootOptions) SetZeroValues() {
	opts.BucketName = ""
	opts.AccessKey = ""
	opts.SecretKey = ""
	opts.Region = ""
	opts.VerboseLog = false
	opts.BannerFilePath = "banner.txt"
	opts.DryRun = false
	opts.AutoApprove = false
}
