# syntax=docker/dockerfile:1

ARG GO_VERSION="1.23"
ARG FINAL_IMAGE="alpine:latest"
ARG BUILD_TAGS="netgo,ledger,muslc"

# --------------------------------------------------------
# Builder
# --------------------------------------------------------

FROM golang:${GO_VERSION}-alpine3.20 as builder

ARG GIT_VERSION
ARG GIT_COMMIT
ARG BUILD_TAGS
ARG CMT_VERSION

RUN apk add --no-cache \
    ca-certificates \
    build-base \
    linux-headers \
    binutils-gold \
    git

# Download go dependencies
WORKDIR /mantrachain
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download

# Cosmwasm - Download correct libwasmvm version
RUN ARCH=$(uname -m) && WASMVM_VERSION=$(go list -m github.com/CosmWasm/wasmvm/v2 | sed 's/.* //') && \
    wget https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/libwasmvm_muslc.$ARCH.a \
    -O /lib/libwasmvm_muslc.$ARCH.a  && \
    # verify checksum
    wget https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/checksums.txt -O /tmp/checksums.txt && \
    sha256sum /lib/libwasmvm_muslc.$ARCH.a | grep $(cat /tmp/checksums.txt | grep libwasmvm_muslc.$ARCH | cut -d ' ' -f 1)

# Copy the remaining files
COPY . .

# Build mantrachaind binary
# build tag info: https://github.com/cosmos/wasmd/blob/master/README.md#supported-systems
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    LEDGER_ENABLED=true BUILD_TAGS='muslc osusergo' LINK_STATICALLY=true make build

# --------------------------------------------------------
# Runner
# --------------------------------------------------------

FROM ${FINAL_IMAGE}

COPY --from=builder /mantrachain/build/mantrachaind /bin/mantrachaind

ENV HOME /mantrachain
WORKDIR $HOME

EXPOSE 26656
EXPOSE 26657
EXPOSE 1317
# Note: uncomment the line below if you need pprof in local mantrachain
# We disable it by default in out main Dockerfile for security reasons
# EXPOSE 6060

ENTRYPOINT ["mantrachaind"]