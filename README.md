# S3 Manager
[![CI](https://github.com/bilalcaliskan/s3-manager/workflows/CI/badge.svg?event=push)](https://github.com/bilalcaliskan/s3-manager/actions?query=workflow%3ACI)
[![Go Report Card](https://goreportcard.com/badge/github.com/bilalcaliskan/s3-manager)](https://goreportcard.com/report/github.com/bilalcaliskan/s3-manager)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_s3-manager&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=bilalcaliskan_s3-manager)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_s3-manager&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=bilalcaliskan_s3-manager)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_s3-manager&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=bilalcaliskan_s3-manager)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_s3-manager&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=bilalcaliskan_s3-manager)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_s3-manager&metric=coverage)](https://sonarcloud.io/summary/new_code?id=bilalcaliskan_s3-manager)
[![Release](https://img.shields.io/github/release/bilalcaliskan/s3-manager.svg)](https://github.com/bilalcaliskan/s3-manager/releases/latest)
[![Go version](https://img.shields.io/github/go-mod/go-version/bilalcaliskan/s3-manager)](https://github.com/bilalcaliskan/s3-manager)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

[s3-manager](https://github.com/bilalcaliskan/s3-manager) is a tool written with Golang that helps you take the following actions on an AWS S3 bucket:
- Sets up a one-shot file cleaning rule that matches specific pattern (great idea to check flags on [cmd/clean/options/options.go](cmd/clean/options/options.go)).
  You can absolutely use [s3-manager](https://github.com/bilalcaliskan/s3-manager) to take advantage of file cleaning for your bucket on your automated operations for **decreasing your public cloud costs**.
- Searches a string in files (supports regex)
- Finds files (supports regex)
- Searches strings in files (supports regex)
- enables/disables/shows tags
- enables/disables/shows versioning configuration
- enables/disables/shows transfer acceleration configuration
- adds/removes/shows bucket policy

About the first and the best benefit of [s3-manager](https://github.com/bilalcaliskan/s3-manager); you can read more from [here](https://aws.amazon.com/s3/pricing/).

## Access Credentials

You can provide access credentials of your AWS account with below environment variables or CLI flags. Keep in mind that command line flags
will override environment variables if you set both of them:
```
"--access-key" CLI flag or "AWS_ACCESS_KEY" environment variable
"--secret-key" CLI flag or "AWS_SECRET_KEY" environment variable
"--region" CLI flag or "AWS_REGION" environment variable
"--bucket-name" CLI flag or "AWS_BUCKET_NAME" environment variable
```

## Available Subcommands
Here is the list of available subcommands of [s3-manager](https://github.com/bilalcaliskan/s3-manager):

- [clean](cmd/clean)
- [search](cmd/search)
- [tags](cmd/tags)
- [versioning](cmd/versioning)
- [bucketpolicy](cmd/bucketpolicy)
- [transferacceleration](cmd/transferacceleration)

## Configuration
```
Usage:
  s3-manager [flags]
  s3-manager [command]

Available Commands:
  bucketpolicy         shows/sets the bucket policy configuration of the target bucket
  clean                finds and clears desired files by a pre-configured rule set
  completion           Generate the autocompletion script for the specified shell
  help                 Help about any command
  search               searches the files which has desired substrings in it
  tags                 shows/sets the tagging configuration of the target bucket
  transferacceleration shows/sets the transfer acceleration configuration of the target bucket
  versioning           shows/sets the versioning configuration of the target bucket

Flags:
      --access-key string         access key credential to access S3 bucket, this value also can be passed via "AWS_ACCESS_KEY" environment variable (default "")
      --banner-file-path string   relative path of the banner file (default "banner.txt")
      --bucket-name string        name of the target bucket on S3, this value also can be passed via "AWS_BUCKET_NAME" environment variable (default "")
  -h, --help                      help for s3-manager
      --region string             region of the target bucket on S3, this value also can be passed via "AWS_REGION" environment variable (default "")
      --secret-key string         secret key credential to access S3 bucket, this value also can be passed via "AWS_SECRET_KEY" environment variable (default "")
      --verbose                   verbose output of the logging library (default false)
  -v, --version                   version for s3-manager

Use "s3-manager [command] --help" for more information about a command.
```

## Installation
### Binary
Binary can be downloaded from [Releases](https://github.com/bilalcaliskan/s3-manager/releases) page.

### Homebrew
This project can also be installed with [Homebrew](https://brew.sh/):
```shell
$ brew tap bilalcaliskan/tap
$ brew install bilalcaliskan/tap/s3-manager
```

## Examples
Here is the couple of examples:
```shell
# bucket cleaning with specified ruleset
$ export AWS_ACCESS_KEY=${YOUR_ACCESS_KEY}
$ export AWS_SECRET_KEY=${YOUR_SECRET_KEY}
$ export AWS_REGION=${YOUR_REGION}
$ export AWS_BUCKET_NAME=${YOUR_BUCKET_NAME}
$ s3-manager clean --min-size-mb=1 --max-size-mb=1000 --keep-last-n-files=2 --sort-by=lastModificationDate

# set bucket versioning as enabled
$ s3-manager versioning set enabled --access-key ${YOUR_ACCESS_KEY} --secret-key ${YOUR_SECRET_KEY} --bucketName ${TARGET_BUCKET_NAME} --region ${TARGET_REGION}

# text search
$ s3-manager search text "catch me if you can" --access-key asdasfasfasfasfasfas --secret-key asdasfasfasfasfasfas --bucket-name demo-bucket --region us-east-2
```

Every subcommand has its own usage examples, please refer to `s3-manager [command] --help` for more information about a command.

## Development
This project requires below tools while developing:
- [Golang 1.20](https://golang.org/doc/go1.20)
- [pre-commit](https://pre-commit.com/)
- [golangci-lint](https://golangci-lint.run/usage/install/) - required by [pre-commit](https://pre-commit.com/)
- [gocyclo](https://github.com/fzipp/gocyclo) - required by [pre-commit](https://pre-commit.com/)

Simply run below command to prepare your development environment:
```shell
$ python3 -m venv venv
$ source venv/bin/activate
$ pip3 install pre-commit
$ pre-commit install -c build/ci/.pre-commit-config.yaml
```
