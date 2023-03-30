# golang-cli-template
[![CI](https://github.com/bilalcaliskan/golang-cli-template/workflows/CI/badge.svg?event=push)](https://github.com/bilalcaliskan/golang-cli-template/actions?query=workflow%3ACI)
[![Docker pulls](https://img.shields.io/docker/pulls/bilalcaliskan/golang-cli-template)](https://hub.docker.com/r/bilalcaliskan/golang-cli-template/)
[![Go Report Card](https://goreportcard.com/badge/github.com/bilalcaliskan/golang-cli-template)](https://goreportcard.com/report/github.com/bilalcaliskan/golang-cli-template)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_golang-cli-template&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=bilalcaliskan_golang-cli-template)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_golang-cli-template&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=bilalcaliskan_golang-cli-template)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_golang-cli-template&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=bilalcaliskan_golang-cli-template)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_golang-cli-template&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=bilalcaliskan_golang-cli-template)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_golang-cli-template&metric=coverage)](https://sonarcloud.io/summary/new_code?id=bilalcaliskan_golang-cli-template)
[![Release](https://img.shields.io/github/release/bilalcaliskan/golang-cli-template.svg)](https://github.com/bilalcaliskan/golang-cli-template/releases/latest)
[![Go version](https://img.shields.io/github/go-mod/go-version/bilalcaliskan/golang-cli-template)](https://github.com/bilalcaliskan/golang-cli-template)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## Required Steps
- Single command is mostly enough to prepare project, it will prompt you with some questions about your new project:
  ```shell
  $ make -s prepare-initial-project
  ```

## Additional nice-to-have steps
- If you want to build and publish Docker image:
  - Ensure `DOCKER_USERNAME` has been added as **repository secret on GitHub**
  - Ensure `DOCKER_PASSWORD` has been added as **repository secret on GitHub**
  - Uncomment **line 166** to **line 173** in [.github/workflows/push.yml](.github/workflows/push.yml)
  - Uncomment **line 32** to **line 50** in [build/package/.goreleaser.yaml](build/package/.goreleaser.yaml)
- If you want to enable https://sonarcloud.io/ integration:
  - Ensure your created repository from that template has been added to https://sonarcloud.io/
  - Ensure `SONAR_TOKEN` has been added as **repository secret** on GitHub
  - Ensure `SONAR_TOKEN` has been added as **dependabot secret** on GitHub
  - Uncomment **line 137** in [.github/workflows/push.yml](.github/workflows/push.yml)
  - Uncomment **line 86** to **line 111** in [.github/workflows/push.yml](.github/workflows/push.yml)
- If you want to create banner:
  - Generate a banner from [here](https://devops.datenkollektiv.de/banner.txt/index.html) and place it inside of [build/ci](build/ci) directory into a file **banner.txt**
  - Uncomment **line 18** and **line 35** to **line 38** in [cmd/root.go](cmd/root.go)
  - Run `go get -u github.com/dimiro1/banner`
- If you want to release as Homebrew Formula:
  - At first, you must have a **formula repository** like https://github.com/bilalcaliskan/homebrew-tap
  - Ensure `TAP_GITHUB_TOKEN` has been added as **repository secret** on GitHub
  - Uncomment **line 186** in [.github/workflows/push.yml](.github/workflows/push.yml)
  - Uncomment **line 70** to **line 80** in [build/package/.goreleaser.yaml](build/package/.goreleaser.yaml)

## Used Libraries
- [spf13/cobra](https://github.com/spf13/cobra)
- [rs/zerolog](https://github.com/rs/zerolog)

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
