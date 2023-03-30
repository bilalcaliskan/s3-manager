/*
Copyright © 2022 bilalcaliskan bilalcaliskan@protonmail.com
*/
package main

import (
	"os"

	"github.com/bilalcaliskan/s3-manager/cmd/root"
)

func main() {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
