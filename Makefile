export CGO_ENABLED=0
VERSION := $(shell git describe --tags --always)
BUILD_INFO := -X "github.com/runbilliam/billiam.Version=$(VERSION)"
FLAGS := -trimpath -ldflags='$(BUILD_INFO) -w -s -extldflags "-static"'
DEV_FLAGS := -trimpath -ldflags='$(BUILD_INFO)'

assets:
	cd setup && npm run build

build: clean
	go generate
	go build -o ./bin/billiam $(FLAGS) cmd/billiam/*
.PHONY: build

dev: clean
	go generate -tags dev
	go build -o ./bin/billiam -tags dev $(DEV_FLAGS) cmd/billiam/*
.PHONY: build-dev

release: clean
	go generate
	GOOS=darwin  GOARCH=amd64 go build -o ./bin/billiam-$(VERSION)-macos-x64 $(FLAGS) cmd/billiam/*
	GOOS=linux   GOARCH=amd64 go build -o ./bin/billiam-$(VERSION)-linux-x64 $(FLAGS) cmd/billiam/*
	GOOS=windows GOARCH=amd64 go build -o ./bin/billiam-$(VERSION)-win64.exe $(FLAGS) cmd/billiam/*
.PHONY: build-all-platforms

clean:
	rm -rf ./bin
	find . -name '.DS_Store' -type f -delete
.PHONY: clean

.DEFAULT_GOAL := build
