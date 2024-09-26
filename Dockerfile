# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23
ARG FINAL_IMAGE=alpine:latest

# --------------------------------------------------------
# Builder
# --------------------------------------------------------

FROM golang:${GO_VERSION}-alpine3.20 AS builder

ARG BUILD_TAGS="netgo,ledger,muslc"
ARG GIT_VERSION
ARG GIT_COMMIT
ARG CMT_VERSION

# Install build dependencies
RUN apk add --no-cache \
    binutils-gold \
    build-base \
    ca-certificates \
    git \
    linux-headers

WORKDIR /mantrachain

# Download go dependencies
# Download and verify libwasmvm
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download \
    ARCH=$(uname -m) && \
    WASMVM_VERSION=$(go list -m github.com/CosmWasm/wasmvm/v2 | sed 's/.* //') && \
    wget -O /lib/libwasmvm_muslc."$ARCH".a \
    "https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/libwasmvm_muslc.$ARCH.a" && \
    wget -O /tmp/checksums.txt \
    "https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/checksums.txt" && \
    sha256sum /lib/libwasmvm_muslc."$ARCH".a | \
    grep $(grep libwasmvm_muslc."$ARCH" /tmp/checksums.txt | cut -d ' ' -f 1) && \
    rm /tmp/checksums.txt

# Copy the remaining files and build
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    LEDGER_ENABLED=true BUILD_TAGS="${BUILD_TAGS}" LINK_STATICALLY=true make build

# --------------------------------------------------------
# Runner
# --------------------------------------------------------

FROM ${FINAL_IMAGE}

COPY --from=builder /mantrachain/build/mantrachaind /usr/local/bin/

ENV HOME=/mantrachain
WORKDIR $HOME

EXPOSE 26656 26657 1317

ENTRYPOINT ["mantrachaind"]