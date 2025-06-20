//go:build unit

package utils

import (
	"context"
	"errors"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/logging"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/prompt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/spf13/cobra"

	"github.com/stretchr/testify/assert"
)

var mockObjects = []types.Object{{
	ChecksumAlgorithm: nil,
	ETag:              aws.String("233b4ce689c7086b7958eb31d8f8b811"),
	Key:               aws.String("bar-service/233b4ce689c7086b7958eb31d8f8b811.template"),
	LastModified:      aws.Time(time.Now()),
	Owner: &types.Owner{
		DisplayName: aws.String("developer1"),
		ID:          aws.String("3becc289963dfc26fe632e4d2fc78d2c7875fc4f030813629e28db2c1fbba4b7"),
	},
	Size:         aws.Int64(2129),
	StorageClass: types.ObjectStorageClassStandard,
}, {
	ChecksumAlgorithm: nil,
	ETag:              aws.String("233b4ce689c7086b7958eb31d8f8b811"),
	Key:               aws.String("foo-service/233b4ce689c7086b7958eb31d8f8b811.template"),
	LastModified:      aws.Time(time.Now()),
	Owner: &types.Owner{
		DisplayName: aws.String("developer1"),
		ID:          aws.String("3becc289963dfc26fe632e4d2fc78d2c7875fc4f030813629e28db2c1fbba4b7"),
	},
	Size:         aws.Int64(2129),
	StorageClass: types.ObjectStorageClassStandard,
}}

func TestContains(t *testing.T) {
	res := Contains([]string{"size", "lastModificationDate"}, "size")
	assert.True(t, res)
}

func TestNotContains(t *testing.T) {
	res := Contains([]string{"size", "lastModificationDate"}, "sizee")
	assert.False(t, res)
}

func TestGetKeysOnly(t *testing.T) {
	keys := GetKeysOnly(mockObjects)
	assert.NotEmpty(t, keys)
}

func TestHasKeyValuePair(t *testing.T) {
	map1 := make(map[string]string)

	ok := HasKeyValuePair(map1, "foo", "bar")
	assert.False(t, ok)
}

func TestRemoveMapElements(t *testing.T) {
	map1 := make(map[string]string)
	map1["foo"] = "bar"
	map2 := make(map[string]string)
	map2["foo"] = "bar"

	RemoveMapElements(map1, map2)
}

func TestPrepareConstants(t *testing.T) {
	cmd := &cobra.Command{}
	ctx := context.Background()
	cmd.SetContext(ctx)

	rootOpts := options.GetMockedRootOptions()
	appCreds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(rootOpts.AccessKey, rootOpts.SecretKey, ""))
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(rootOpts.Region),
		config.WithCredentialsProvider(appCreds),
	)
	assert.Nil(t, err)

	client := s3.NewFromConfig(cfg)
	assert.NotNil(t, client)

	cases := []struct {
		caseName string
		prompt.PromptRunner
	}{
		{
			"Success with all passed",
			prompt.GetConfirmRunner(),
		},
		{
			"Success with not-passed confirm runner",
			nil,
		},
	}

	for _, tc := range cases {
		cmd.SetContext(context.WithValue(cmd.Context(), options.LoggerKey{}, logging.GetLogger(rootOpts)))
		cmd.SetContext(context.WithValue(cmd.Context(), options.OptsKey{}, rootOpts))
		cmd.SetContext(context.WithValue(cmd.Context(), options.S3ClientKey{}, client))
		cmd.SetContext(context.WithValue(cmd.Context(), options.ConfirmRunnerKey{}, tc.PromptRunner))

		returnSvc, returnOpts, returnLogger, returnPrompt := PrepareConstants(cmd)
		if tc.PromptRunner != nil {
			assert.NotNil(t, returnSvc, returnOpts, returnLogger, returnPrompt)
		} else {
			assert.NotNil(t, returnSvc, returnOpts, returnLogger)
			assert.Nil(t, returnPrompt)
		}
	}
}

func TestBeautifyJSON(t *testing.T) {
	cases := []struct {
		caseName   string
		input      string
		shouldPass bool
	}{
		{
			"Success",
			`
{
  "Statement": [
    {
      "Action": "s3:*",
      "Condition": {
        "Bool": {
          "aws:SecureTransport": "false"
        }
      },
      "Effect": "Deny",
      "Principal": "*",
      "Resource": [
        "arn:aws:s3:::thevpnbeast-releases-1",
        "arn:aws:s3:::thevpnbeast-releases-1/*"
      ],
      "Sid": "RestrictToTLSRequestsOnly"
    }
  ],
  "Version": "2012-10-17"
}
`,
			true,
		},
		{
			"Failure caused by invalid json",
			`
{
  "Statement": [
    {
      "Action": "s3:*",
      "Condition": {
        "Bool": {
          "aws:SecureTransport": "false"
        }
      },
      "Effect": "Deny",
      "Principal": "*",
      "Resource": [
        "arn:aws:s3:::thevpnbeast-releases-1",
        "arn:aws:s3:::thevpnbeast-releases-1/*"
      ],
      "Sid": "RestrictToTLSRequestsOnly"
    }
  ]
  "Version": "2012-10-17"

`,
			false,
		},
	}

	for _, tc := range cases {
		res, err := BeautifyJSON(tc.input)

		if tc.shouldPass {
			assert.Nil(t, err)
			assert.NotEqual(t, "", res)
		} else {
			assert.NotNil(t, err)
			assert.Equal(t, "", res)
		}
	}
}

func TestCheckArgs(t *testing.T) {
	args := []string{"foo", "bar"}
	cases := []struct {
		caseName    string
		allowed     int
		expectedErr error
	}{
		{
			caseName:    "Success",
			allowed:     2,
			expectedErr: nil,
		},
		{
			caseName:    "Failure caused by too much arguments",
			allowed:     1,
			expectedErr: errors.New("too many arguments provided"),
		},
		{
			caseName:    "Failure caused by too few arguments",
			allowed:     3,
			expectedErr: errors.New("too few arguments provided"),
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		err := CheckArgs(args, tc.allowed)
		assert.Equal(t, tc.expectedErr, err)
	}
}
