SHELL := /bin/bash
BASEDIR = $(shell pwd)
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
VersionDir = github.com/telepair/telepair/pkg/version
TZ := Asia/Shanghai	
VERSION := v0.0.1
ENV := dev
BINARY_NAME := telepair
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

gitHash = $(shell git rev-parse HEAD) 
gitBranch = $(shell git rev-parse --abbrev-ref HEAD)
gitTag = $(shell \
            if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ]; \
            then \
                git describe --tags --abbrev=0; \
            else \
                git log --pretty=format:'%h' -n 1; \
            fi)
gitCommit = $(shell git log --pretty=format:'%H' -n 1)
gitTreeState = $(shell if git status | grep -q 'clean'; then echo clean; else echo dirty; fi)
buildDate = $(shell TZ=$(TZ) date +%FT%T%z)

LDFLAGS += -X '${VersionDir}.gitTag=${gitTag}'
LDFLAGS += -X '${VersionDir}.buildDate=${buildDate}'
LDFLAGS += -X '${VersionDir}.gitCommit=${gitCommit}'
LDFLAGS += -X '${VersionDir}.gitTreeState=${gitTreeState}'
LDFLAGS += -X '${VersionDir}.gitBranch=${gitBranch}'
LDFLAGS += -X '${VersionDir}.version=${VERSION}'

ifeq ($(ENV), dev)
    BUILD_FLAGS = -race
endif

ifeq ($(ENV), pro)
    LDFLAGS = -w 
endif

all: lint test run

run: build
	@echo "  >  Running binary..."
	$(GOBIN)/$(BINARY_NAME)

build:
	mkdir -p $(GOBIN)
	go build  -v -ldflags "$(LDFLAGS)" $(BUILD_FLAGS) -gcflags=all="-N -l"  -o $(GOBIN)/$(BINARY_NAME) .

lint:
	go fmt ./...
	go vet ./...
	goimports -l -w .
	golangci-lint run ./...

test: lint
	go test -v -run '^Test' ./...

example:
	go test -v -run '^Example' ./...

bench:
	go test -v -bench=. ./...

cover:
	go test ./... -v -short -coverprofile .coverage.txt
	go tool cover -func .coverage.txt

clean: 
	rm -f $(GOBIN)/$(BINARY_NAME) || true
	find . -name "[._]*.s[a-w][a-z]" | xargs rm -f || true
	rm -rf ./log || true
	go clean
	rm -f go.sum; go mod tidy

help:
	@echo "make all      - run lint, test, build"
	@echo "make run      - run the binary file"
	@echo "make build    - compile the source code"
	@echo "make lint     - run go tool 'fmt', 'vet', 'goimports', 'golangci-lint' "
	@echo "make test     - run go test"
	@echo "make example  - run go test with example"
	@echo "make bench    - run go test with benchmark"
	@echo "make cover    - run go test with coverage"
	@echo "make clean    - remove binary file and vim swp files"
	@echo "make help     - show the help info"

.PHONY: all run build lint test example bench cover clean help
