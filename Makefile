VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u '+%Y-%m-%d %H:%M:%S')
GO_VERSION := $(shell go version | cut -d' ' -f3)

LDFLAGS := -X 'github.com/conormkelly/yts-cli/cmd.version=$(VERSION)' \
           -X 'github.com/conormkelly/yts-cli/cmd.commit=$(COMMIT)' \
           -X 'github.com/conormkelly/yts-cli/cmd.buildDate=$(BUILD_DATE)' \
           -X 'github.com/conormkelly/yts-cli/cmd.goVersion=$(GO_VERSION)'

.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o bin/yts

.PHONY: install
install:
	go install -ldflags "$(LDFLAGS)"

.PHONY: release
release:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/yts-linux-amd64
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/yts-linux-arm64
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/yts-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/yts-darwin-arm64
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/yts-windows-amd64.exe
