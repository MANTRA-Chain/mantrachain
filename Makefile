#!/usr/bin/make -f

# the subcommands are located in the specific makefiles
include scripts/makefiles/lint.mk
include scripts/makefiles/proto.mk

.DEFAULT_GOAL := help
help:
	@echo "Available top-level commands:"
	@echo ""
	@echo "Usage:"
	@echo "    make [command]"
	@echo ""
	@echo "  make build                 Build mantrachaind binary"
	@echo "  make lint                  Show available lint commands"
	@echo "  make test                  Show available test commands"
	@echo "  make proto                 Show available proto commands"
	@echo ""
	@echo "Run 'make [subcommand]' to see the available commands for each subcommand."

LEDGER_ENABLED ?= true
BINDIR ?= $(GOPATH)/bin
BUILDDIR ?= $(CURDIR)/build
DOCKER := $(shell which docker)

BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2> /dev/null)
BRANCH_PRETTY := $(subst /,-,$(BRANCH))
export CMT_VERSION := $(shell go list -m github.com/cometbft/cometbft 2> /dev/null | sed 's:.* ::')
COMMIT := $(shell git log -1 --format='%h' 2> /dev/null)
# don't override user values
ifeq (,$(VERSION))
  VERSION := $(shell git describe --exact-match --tags 2>/dev/null)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH_PRETTY)-$(COMMIT)
  endif
endif

# Go version to be used in docker images
GO_VERSION := $(shell cat go.mod | grep -E 'go [0-9].[0-9]+' | cut -d ' ' -f 2)
# currently installed Go version
GO_MODULE := $(shell cat go.mod | grep "module " | cut -d ' ' -f 2)

###############################################################################
###                            Build Flags/Tags                             ###
###############################################################################

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger pebbledb
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger pebbledb
      endif
    endif
  endif
endif

ifeq (cleveldb,$(findstring cleveldb,$(MANTRACHAIN_BUILD_OPTIONS)))
  build_tags += gcc
else ifeq (rocksdb,$(findstring rocksdb,$(MANTRACHAIN_BUILD_OPTIONS)))
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace := $(whitespace) $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=mantrachain \
	-X github.com/cosmos/cosmos-sdk/version.AppName=mantrachaind \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
	-X github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep) \
	-X github.com/cometbft/cometbft/version.TMCoreSemVer=$(CMT_VERSION)

ifeq (cleveldb,$(findstring cleveldb,$(MANTRACHAIN_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
else ifeq (rocksdb,$(findstring rocksdb,$(MANTRACHAIN_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb
endif
ifeq (,$(findstring nostrip,$(MANTRACHAIN_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ifeq ($(LINK_STATICALLY),true)
	ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static -lm"
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(MANTRACHAIN_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

###############################################################################
###                                  Build                                  ###
###############################################################################

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)/
build-arm:
	GOOS=darwin GOARCH=arm64 $(MAKE) build
build-linux:
	GOOS=linux GOARCH=$(if $(findstring aarch64,$(shell uname -m)) || $(findstring arm64,$(shell uname -m)),arm64,amd64) $(MAKE) build
build-image:
	docker build -f Dockerfile -t mantra-chain/mantrachain .

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) $(GO_MODULE)/cmd/mantrachaind
$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

###############################################################################
###                           Tests                            		        ###
###############################################################################

PACKAGES_UNIT=$(shell go list ./... | grep -v -e '/tests/e2e' | grep -v '/simulation')
PACKAGES_E2E=$(shell cd tests/e2e && go list ./... | grep '/e2e')
TEST_PACKAGES=./...
TEST_TARGETS := test-unit test-e2e test-cover test-connect

DIR=$(CURDIR)
test-unit: ARGS=-timeout=5m -tags='norace'
test-unit: TEST_PACKAGES=$(PACKAGES_UNIT)
test-e2e: ARGS=-timeout=35m -v
test-e2e: TEST_PACKAGES=$(PACKAGES_E2E)
test-e2e: build-image
test-cover: ARGS=-timeout=30m -coverprofile=coverage.txt -covermode=atomic -tags='norace'
test-cover: TEST_PACKAGES=$(PACKAGES_UNIT)
test-connect: ARGS=-v -race
test-connect: DIR=$(CURDIR)/tests/connect
test-connect: build-image
$(TEST_TARGETS): run-tests

run-tests:
ifneq (,$(shell which tparse 2>/dev/null))
	@echo "--> Running tests"
	@cd $(DIR) && go test -mod=readonly -json $(ARGS) $(TEST_PACKAGES) | tparse
else
	@echo "--> Running tests"
	cd $(DIR) && go test -mod=readonly $(ARGS) $(TEST_PACKAGES)
endif

###############################################################################
###                                Release                                  ###
###############################################################################
ifeq ($(strip $(GORELEASER_CROSS_DISABLE)),true)
GORELEASER_IMAGE := goreleaser/goreleaser:v2.8.2
else
GORELEASER_CROSS := ghcr.io/goreleaser/goreleaser-cross
GO_VERSION_FALLBACK := 1.24.1
GORELEASER_IMAGE := $(shell docker manifest inspect $(GORELEASER_CROSS):v$(GO_VERSION) > /dev/null 2>&1 && echo $(GORELEASER_CROSS):v$(GO_VERSION) || echo $(GORELEASER_CROSS):v$(GO_VERSION_FALLBACK))
endif
GORELEASER_PLATFORM ?= linux/amd64
COSMWASM_VERSION := $(shell go list -m github.com/CosmWasm/wasmvm/v3 | sed 's/.* //')
REPO_OWNER ?= MANTRA-Chain
REPO_NAME ?= mantrachain

# Check if GITHUB_TOKEN is defined
ifndef GITHUB_TOKEN
MISSING_TOKEN := GITHUB_TOKEN
endif

ifeq ($(strip $(MISSING_TOKEN)),)
release:
	docker run \
		--rm \
		-e GITHUB_TOKEN=$(GITHUB_TOKEN) \
		-e COSMWASM_VERSION=$(COSMWASM_VERSION) \
		-e CMT_VERSION=$(CMT_VERSION) \
		-e REPO_OWNER=$(REPO_OWNER) \
		-e REPO_NAME=$(REPO_NAME) \
		-v `pwd`:/go/src/mantrachaind \
		-w /go/src/mantrachaind \
		--platform=$(GORELEASER_PLATFORM) \
		$(GORELEASER_IMAGE) \
		release $(if $(GORELEASER_SKIP),--skip=$(GORELEASER_SKIP)) $(if $(GORELEASER_CONFIG),--config=$(GORELEASER_CONFIG)) \
		--clean
else
release:
	@echo "Error: $(MISSING_TOKEN) is not defined. Please define it before running 'make release'."
endif

# uses goreleaser to create static binaries for linux and darwin on local machine
# platform is set because not setting it results in broken builds for linux-amd64
goreleaser-build-local:
	docker run \
		--rm \
		-e COSMWASM_VERSION=$(COSMWASM_VERSION) \
		-e CMT_VERSION=$(CMT_VERSION) \
		-e REPO_OWNER=$(REPO_OWNER) \
		-e REPO_NAME=$(REPO_NAME) \
		-v `pwd`:/go/src/mantrachaind \
		-w /go/src/mantrachaind \
		--platform=$(GORELEASER_PLATFORM) \
		$(GORELEASER_IMAGE) \
		build $(if $(GORELEASER_IDS),$(shell echo $(GORELEASER_IDS) | tr ',' ' ' | sed 's/[^ ]*/--id=&/g')) \
		--skip=validate $(if $(filter true,$(GORELEASER_SNAPSHOT)),--snapshot) \
		--clean \
		--timeout 90m \
		--verbose

.PHONY: build build-linux lint release

###############################################################################
###                           Mocks                            		    ###
###############################################################################

mocks:
	go generate ./...


###############################################################################
###                           Single Node Test                              ###
###############################################################################

build-and-run-single-node: build
	@echo "Building and running a single node for testing..."
	@mkdir -p .mantrasinglenodetest
	@if [ ! -f .mantrasinglenodetest/config/config.toml ]; then \
		./build/mantrachaind init single-node-test --chain-id test-chain --home .mantrasinglenodetest --default-denom amantra; \
		./build/mantrachaind keys add validator --keyring-backend test --home .mantrasinglenodetest; \
		./build/mantrachaind genesis add-genesis-account $$(./build/mantrachaind keys show validator -a --keyring-backend test --home .mantrasinglenodetest) 100000000000000000000000000amantra --home .mantrasinglenodetest; \
		./build/mantrachaind genesis gentx validator 100000000000000000000amantra --chain-id test-chain --keyring-backend test --home .mantrasinglenodetest; \
		./build/mantrachaind genesis collect-gentxs --home .mantrasinglenodetest; \
		sed -i'' -e 's/"fee_denom": "stake"/"fee_denom": "amantra"/' .mantrasinglenodetest/config/genesis.json; \
	fi
	./build/mantrachaind start --home .mantrasinglenodetest --minimum-gas-prices 0amantra

.PHONY: build-and-run-single-node
