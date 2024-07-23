#!/usr/bin/make -f

BINDIR ?= $(GOPATH)/bin
BUILDDIR ?= $(CURDIR)/build
PACKAGES_NOSIMULATION=$(shell go list ./... | grep -v '/simulation')

###############################################################################
###                                  Build                                  ###
###############################################################################

BUILD_TARGETS := build install

build: BUILD_FLAGS=-ldflags="-X 'github.com/cosmos/cosmos-sdk/version.Version=0.5.0'"
build: BUILD_ARGS=-o $(BUILDDIR)/
build-arm:
	GOOS=darwin GOARCH=arm64 $(MAKE) build
build-linux:
	GOOS=linux GOARCH=$(if $(findstring aarch64,$(shell uname -m)) || $(findstring arm64,$(shell uname -m)),arm64,amd64) $(MAKE) build

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

.PHONY: build build-linux

###############################################################################
###                           Tests                            							###
###############################################################################

test: test-unit

test-unit: 
	@VERSION=$(VERSION) go test ./x/... -mod=readonly -vet=all -tags='norace' $(PACKAGES_NOSIMULATION)

test-cover:
	@VERSION=$(VERSION) go test ./x/... -mod=readonly -timeout 30m -coverprofile=coverage.txt -covermode=atomic -tags='norace' $(PACKAGES_NOSIMULATION)

mocks:
	@go install github.com/golang/mock/mockgen@v1.6.0
	sh ./scripts/mockgen.sh