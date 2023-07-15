//go:build e2e

package add

import (
	"context"
	"errors"
	"testing"

	"github.com/bilalcaliskan/s3-manager/internal/prompt"

	internalaws "github.com/bilalcaliskan/s3-manager/internal/aws"
	"github.com/bilalcaliskan/s3-manager/internal/constants"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
)

type promptMock struct {
	msg string
	err error
}

func (p promptMock) Run() (string, error) {
	return p.msg, p.err
}

type mockS3Client struct {
	mock.Mock
	s3iface.S3API
}

func (m *mockS3Client) PutBucketPolicy(input *s3.PutBucketPolicyInput) (*s3.PutBucketPolicyOutput, error) {
	// Return the mocked output values using the `On` method of testify/mock
	args := m.Called(input)
	return args.Get(0).(*s3.PutBucketPolicyOutput), args.Error(1)
}

func TestExecuteAddCmd(t *testing.T) {
	ctx := context.Background()
	AddCmd.SetContext(ctx)

	rootOpts := options.GetMockedRootOptions()
	svc, err := internalaws.CreateAwsService(rootOpts)
	assert.NotNil(t, svc)
	assert.Nil(t, err)

	confirmRunner := prompt.GetConfirmRunner()
	assert.NotNil(t, confirmRunner)

	cases := []struct {
		caseName              string
		args                  []string
		shouldPass            bool
		putBucketPolicyErr    error
		putBucketPolicyOutput *s3.PutBucketPolicyOutput
		promptMock            *promptMock
		dryRun                bool
		autoApprove           bool
	}{
		{
			"Success",
			[]string{"../../../testdata/bucketpolicy.json"},
			true,
			nil,
			&s3.PutBucketPolicyOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
		{
			"Success with dry-run",
			[]string{"../../../testdata/bucketpolicy.json"},
			true,
			nil,
			&s3.PutBucketPolicyOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			true,
			false,
		},
		{
			"Failure",
			[]string{"../../../testdata/bucketpolicy.json"},
			false,
			errors.New("dummy error"),
			&s3.PutBucketPolicyOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by user terminated process",
			[]string{"../../../testdata/bucketpolicy.json"},
			false,
			nil,
			&s3.PutBucketPolicyOutput{},
			&promptMock{
				msg: "n",
				err: constants.ErrInjected,
			},
			false,
			false,
		},
		{
			"Failure caused by prompt error",
			[]string{"../../../testdata/bucketpolicy.json"},
			false,
			nil,
			&s3.PutBucketPolicyOutput{},
			&promptMock{
				msg: "nasdasd",
				err: constants.ErrInjected,
			},
			false,
			false,
		},
		{
			"Failure caused by target file not found",
			[]string{"../../../testdata/bucketpolicy.jsonnnn"},
			false,
			nil,
			&s3.PutBucketPolicyOutput{},
			&promptMock{
				msg: "y",
				err: nil,
			},
			false,
			false,
		},
		{
			"Failure caused by too many arguments error",
			[]string{"enabled", "foo"},
			false,
			nil,
			&s3.PutBucketPolicyOutput{},
			nil,
			false,
			false,
		},
		{
			"Failure caused by no arguments provided error",
			[]string{},
			false,
			nil,
			&s3.PutBucketPolicyOutput{},
			nil,
			false,
			false,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		rootOpts.DryRun = tc.dryRun
		rootOpts.AutoApprove = tc.autoApprove

		mockS3 := new(mockS3Client)
		mockS3.On("PutBucketPolicy", mock.AnythingOfType("*s3.PutBucketPolicyInput")).Return(tc.putBucketPolicyOutput, tc.putBucketPolicyErr)

		AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.S3SvcKey{}, mockS3))
		AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.OptsKey{}, rootOpts))
		AddCmd.SetContext(context.WithValue(AddCmd.Context(), options.ConfirmRunnerKey{}, tc.promptMock))

		AddCmd.SetArgs(tc.args)

		err = AddCmd.Execute()
		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}

		//mockS3.AssertCalled(t, "PutBucketPolicy", mock.AnythingOfType("*s3.PutBucketPolicyInput"))
	}

	bucketPolicyOpts.SetZeroValues()
}
