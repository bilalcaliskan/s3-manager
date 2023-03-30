package utils

import "github.com/aws/aws-sdk-go/service/s3"

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GetKeysOnly(s []*s3.Object) []string {
	var res []string

	for _, v := range s {
		res = append(res, *v.Key)
	}

	return res
}
