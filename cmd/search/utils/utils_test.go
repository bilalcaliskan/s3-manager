package utils

import (
	"context"
	"testing"

	options3 "github.com/bilalcaliskan/s3-manager/cmd/search/options"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	options2 "github.com/bilalcaliskan/s3-manager/cmd/root/options"
	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func createSvc(rootOpts *options2.RootOptions) (*s3.S3, error) {
	return internalaws.CreateAwsService(rootOpts)
}

func TestCheckFlags(t *testing.T) {
	cases := []struct {
		caseName   string
		flags      []string
		shouldPass bool
	}{
		{"Success",
			[]string{"file"},
			true,
		},
		{"Failure caused by no arguments provided",
			[]string{},
			false,
		},
		{"Failure caused by too many arguments provided",
			[]string{"foo", "bar"},
			false,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case '%s'", tc.caseName)

		err := CheckFlags(tc.flags)

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		t.Logf("ending case '%s'", tc.caseName)
	}
}

func TestPrepareConstants(t *testing.T) {
	var (
		svc        s3iface.S3API
		searchOpts *options3.SearchOptions
		logger     zerolog.Logger
	)

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	rootOpts := options2.GetRootOptions()
	rootOpts.AccessKey = "thisisaccesskey"
	rootOpts.SecretKey = "thisissecretkey"
	rootOpts.Region = "thisisregion"
	rootOpts.BucketName = "thisisbucketname"

	svc, err := createSvc(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	cmd.SetContext(context.WithValue(context.Background(), options2.OptsKey{}, rootOpts))
	cmd.SetContext(context.WithValue(cmd.Context(), options2.S3SvcKey{}, svc))

	svc, searchOpts, logger = PrepareConstants(cmd, options3.GetSearchOptions())
	assert.NotNil(t, svc)
	assert.NotNil(t, searchOpts)
	assert.NotNil(t, logger)
}
