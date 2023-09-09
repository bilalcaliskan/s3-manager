# S3 Manager :floppy_disk: :cloud: :man_office_worker:
![CI](https://github.com/bilalcaliskan/s3-manager/workflows/CI/badge.svg?event=push)
![Go Report Card](https://goreportcard.com/badge/github.com/bilalcaliskan/s3-manager)
![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_s3-manager&metric=alert_status)
![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_s3-manager&metric=sqale_rating)
![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_s3-manager&metric=reliability_rating)
![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_s3-manager&metric=security_rating)
![Coverage](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_s3-manager&metric=coverage)
![Release](https://img.shields.io/github/release/bilalcaliskan/s3-manager.svg)
![Go version](https://img.shields.io/github/go-mod/go-version/bilalcaliskan/s3-manager)
![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)

S3 Manager is robust, powerful and an all-in-one tool that provides a number of operations to enhance your experience
with AWS S3 buckets. It offers an array of features from setting up file cleaning rules, search string in files,
manipulating tags, and much more. It's designed to help you automate and optimize your operations, lowering your [public cloud costs](https://aws.amazon.com/s3/pricing/)
in the process.

## Available Subcommands
S3 Manager provides the following subcommands:

- [clean](cmd/clean)
- [search](cmd/search)
- [tags](cmd/tags)
- [versioning](cmd/versioning)
- [bucketpolicy](cmd/bucketpolicy)
- [transferacceleration](cmd/transferacceleration)

<!-- Add a command and its description -->
## Configuration
```shell
Usage:
  s3-manager [flags]
  s3-manager [command]

Available Commands:
  bucketpolicy         Shows/sets the bucket policy configuration of the target bucket
  clean                Finds and clears desired files by a pre-configured rule set
  completion           Generate the autocompletion script for the specified shell
  help                 Help about any command
  search               Searches the files which has desired substrings in it
  tags                 Shows/sets the tagging configuration of the target bucket
  transferacceleration Shows/sets the transfer acceleration configuration of the target bucket
  versioning           Shows/sets the versioning configuration of the target bucket

Flags:
  --access-key string         Access key credential to access S3 bucket, this value also can be passed via "AWS_ACCESS_KEY" environment variable (default "")
  --banner-file-path string   Relative path of the banner file (default "banner.txt")
  --bucket-name string        Name of the target bucket on S3, this value also can be passed via "AWS_BUCKET_NAME" environment variable (default "")
  -h, --help                  Help for s3-manager
  --region string             Region of the target bucket on S3, this value also can be passed via "AWS_REGION" environment variable (default "")
  --secret-key string         Secret key credential to access S3 bucket, this value also can be passed via "AWS_SECRET_KEY" environment variable (default "")
  --verbose                   Verbose output of the logging library (default false)
  -v, --version               Version for s3-manager

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


After you installed [pre-commit](https://pre-commit.com/) and the rest, simply run below command to prepare your
development environment:
```shell
$ make pre-commit-setup
```
