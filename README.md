# S3 Manager
[![CI](https://github.com/bilalcaliskan/s3-manager/workflows/CI/badge.svg?event=push)](https://github.com/bilalcaliskan/s3-manager/actions?query=workflow%3ACI)
[![Docker pulls](https://img.shields.io/docker/pulls/bilalcaliskan/s3-manager)](https://hub.docker.com/r/bilalcaliskan/s3-manager/)
[![Go Report Card](https://goreportcard.com/badge/github.com/bilalcaliskan/s3-manager)](https://goreportcard.com/report/github.com/bilalcaliskan/s3-manager)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_s3-manager&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=bilalcaliskan_s3-manager)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_s3-manager&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=bilalcaliskan_s3-manager)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_s3-manager&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=bilalcaliskan_s3-manager)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_s3-manager&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=bilalcaliskan_s3-manager)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_s3-manager&metric=coverage)](https://sonarcloud.io/summary/new_code?id=bilalcaliskan_s3-manager)
[![Release](https://img.shields.io/github/release/bilalcaliskan/s3-manager.svg)](https://github.com/bilalcaliskan/s3-manager/releases/latest)
[![Go version](https://img.shields.io/github/go-mod/go-version/bilalcaliskan/s3-manager)](https://github.com/bilalcaliskan/s3-manager)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

**TBD**

You can provide access credentials of your AWS account with below environment variables or CLI flags. Keep in mind that command line flags
will override environment variables if you set both of them:
```
"--accessKey" CLI flag or "AWS_ACCESS_KEY" environment variable
"--secretKey" CLI flag or "AWS_SECRET_KEY" environment variable
"--region" CLI flag or "AWS_REGION" environment variable
"--bucketName" CLI flag or "AWS_BUCKET_NAME" environment variable
```

## Configuration
**TBD**

## Installation
### Binary
Binary can be downloaded from [Releases](https://github.com/bilalcaliskan/s3-manager/releases) page.

After then, you can simply run binary by providing required command line arguments:
```shell
$ ./s3-manager search --accessKey asdasfasfasfasfasfas --secretKey asdasfasfasfasfasfas --bucketName demo-bucket --region us-east-2 --substring "catch me if you can"
```

### Homebrew
This project can also be installed with [Homebrew](https://brew.sh/):
```shell
$ brew tap bilalcaliskan/tap
$ brew install bilalcaliskan/tap/s3-manager
```

Then similar to binary method, you can run it by calling below command:
```shell
$ s3-manager search --accessKey asdasfasfasfasfasfas --secretKey asdasfasfasfasfasfas --bucketName demo-bucket --region us-east-2 --substring "catch me if you can"
```

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
