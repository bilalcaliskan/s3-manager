package utils

import (
	"encoding/json"
	"errors"
	internalawstypes "github.com/bilalcaliskan/s3-manager/internal/pkg/aws/types"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/logging"
	"github.com/bilalcaliskan/s3-manager/internal/pkg/prompt"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/bilalcaliskan/s3-manager/cmd/root/options"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GetKeysOnly(s []types.Object) []string {
	var res []string

	for _, v := range s {
		res = append(res, *v.Key)
	}

	return res
}

func RemoveMapElements(source, toRemove map[string]string) {
	for key := range toRemove {
		delete(source, key)
	}
}

func HasKeyValuePair(m map[string]string, key, value string) bool {
	v, ok := m[key]
	return ok && v == value
}

func BeautifyJSON(jsonString string) (string, error) {
	var jsonData interface{}

	err := json.Unmarshal([]byte(jsonString), &jsonData)
	if err != nil {
		return "", err
	}

	beautifiedBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return "", err
	}

	return string(beautifiedBytes), nil
}

func PrepareConstants(cmd *cobra.Command) (internalawstypes.S3ClientAPI, *options.RootOptions, zerolog.Logger, prompt.PromptRunner) {
	svc := cmd.Context().Value(options.S3ClientKey{}).(internalawstypes.S3ClientAPI)
	rootOpts := cmd.Context().Value(options.OptsKey{}).(*options.RootOptions)

	confirmRunner, ok := cmd.Context().Value(options.ConfirmRunnerKey{}).(prompt.PromptRunner)
	if !ok {
		confirmRunner = nil
	}

	return svc, rootOpts, logging.GetLogger(rootOpts), confirmRunner
}

func CheckArgs(args []string, allowed int) error {
	if len(args) > allowed {
		return errors.New("too many arguments provided")
	} else if len(args) < allowed {
		return errors.New("too few arguments provided")
	}

	return nil
}
