GOLANGCI_LINT_VERSION = latest
REVIVE_VERSION = latest
GOIMPORTS_VERSION = latest
INEFFASSIGN_VERSION = latest


LOCAL_BIN := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))/.bin

.PHONY: all
all: clean tools lint fmt test build

.PHONY: clean
clean:
	rm -rf $(LOCAL_BIN)

.PHONY: tools
tools:  golangci-lint-install revive-install go-imports-install ineffassign-install
	go mod tidy

.PHONY: golangci-lint-install
golangci-lint-install:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: revive-install
revive-install:
	GOBIN=$(LOCAL_BIN) go install github.com/mgechev/revive@$(REVIVE_VERSION)

.PHONY: ineffassign-install
ineffassign-install:
	GOBIN=$(LOCAL_BIN) go install github.com/gordonklaus/ineffassign@$(INEFFASSIGN_VERSION)

.PHONY: lint
lint: tools run-lint

.PHONY: run-lint
run-lint: lint-golangci-lint lint-revive

.PHONY: lint-golangci-lint
lint-golangci-lint:
	#$(info running golangci-lint...)
	echo "running golangci-lint..."
	$(LOCAL_BIN)/golangci-lint -v run ./... || (echo golangci-lint returned an error, exiting!; sh -c 'exit 1';)
	echo "golangci-lint exited successfully!"

.PHONY: lint-revive
lint-revive:
	echo "running revive..."
	$(LOCAL_BIN)/revive -formatter=stylish -config=build/ci/.revive.toml -exclude ./vendor/... ./... || (echo revive returned an error, exiting!; sh -c 'exit 1';)
	echo "revive exited successfully!"

.PHONY: upgrade-direct-deps
upgrade-direct-deps: tidy
	for item in `grep -v 'indirect' go.mod | grep '/' | cut -d ' ' -f 1`; do \
		echo "trying to upgrade direct dependency $$item" ; \
		go get -u $$item ; \
  	done
	go mod tidy
	go mod vendor

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: run-goimports
run-goimports: go-imports-install
	for item in `find . -type f -name '*.go' -not -path './vendor/*'`; do \
		$(LOCAL_BIN)/goimports -l -w $$item ; \
	done

.PHONY: go-imports-install
go-imports-install:
	GOBIN=$(LOCAL_BIN) go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)

.PHONY: fmt
fmt: tools run-fmt run-ineffassign run-vet

.PHONY: run-fmt
run-fmt:
	echo "running fmt..."
	go fmt ./... || (echo fmt returned an error, exiting!; sh -c 'exit 1';)
	echo "fmt exited successfully!"

.PHONY: run-ineffassign
run-ineffassign:
	echo "running ineffassign..."
	$(LOCAL_BIN)/ineffassign ./... || (echo ineffassign returned an error, exiting!; sh -c 'exit 1';)
	echo "ineffassign exited successfully!"

.PHONY: run-vet
run-vet:
	echo "running vet..."
	go vet ./... || (echo vet returned an error, exiting!; sh -c 'exit 1';)
	echo "vet exited successfully!"

.PHONY: test
test: tidy
	echo "starting the test for whole module..."
	CGO_ENABLED=1 go test -failfast -vet=off -race ./... || (echo an error while testing, exiting!; sh -c 'exit 1';)

.PHONY: test-with-coverage
test-with-coverage: tidy
	CGO_ENABLED=1 go test ./... -race -coverprofile=coverage.txt -covermode=atomic
	go tool cover -html=coverage.txt -o cover.html

.PHONY: update
update: tidy
	go get -u ./...

.PHONY: build
build: tidy
	echo "building binary..."
	go build -o bin/main main.go || (echo an error while building binary, exiting!; sh -c 'exit 1';)
	echo "binary built successfully!"

.PHONY: run
run: tidy
	go run main.go

.PHONY: cross-compile
cross-compile:
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 main.go
	GOOS=darwin GOARCH=386 go build -o bin/main-darwin-386 main.go
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 main.go
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 main.go
	GOOS=freebsd GOARCH=amd64 go build -o bin/main-freebsd-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/main-darwin-amd64 main.go
	GOOS=linux GOARCH=amd64 go build -o bin/main-linux-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -o bin/main-windows-amd64 main.go


.PHONY: prepare-initial-project
GITHUB_USERNAME ?= $(shell read -p "Your Github username(ex: bilalcaliskan): " github_username; echo $$github_username)
PROJECT_NAME ?= $(shell read -p "'Kebab-cased' Project Name(ex: s3-manager): " project_name; echo $$project_name)
prepare-initial-project:
	grep -rl bilalcaliskan . --exclude={README.md,Makefile} --exclude-dir=.git --exclude-dir=.idea | xargs sed -i 's/bilalcaliskan/$(GITHUB_USERNAME)/g'
	grep -rl s3-manager . --exclude-dir=.git --exclude-dir=.idea | xargs sed -i 's/s3-manager/$(PROJECT_NAME)/g'
	echo "Please refer to *Additional nice-to-have steps* in README.md for additional features"
	echo "Cheers!"
