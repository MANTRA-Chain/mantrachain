# syntax=docker/dockerfile:1

ARG GO_VERSION="1.24"
ARG IMG_TAG="latest"
ARG BUILD_TAGS="netgo,ledger,muslc"

# --------------------------------------------------------
# Builder
# --------------------------------------------------------

FROM golang:${GO_VERSION}-alpine3.20 AS mantra-builder
WORKDIR /src/app/
RUN apk add --no-cache \
    ca-certificates \
    build-base \
    linux-headers \
    binutils-gold \
    git

# Download go dependencies
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/nonroot/.cache/go-build \
    --mount=type=cache,target=/nonroot/go/pkg/mod \
    go mod download

# Cosmwasm - Download correct libwasmvm version
RUN ARCH=$(uname -m) && WASMVM_VERSION=$(go list -m github.com/CosmWasm/wasmvm/v3 | sed 's/.* //') && \
    wget https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/libwasmvm_muslc.$ARCH.a \
    -O /lib/libwasmvm_muslc.$ARCH.a  && \
    # verify checksum
    wget https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/checksums.txt -O /tmp/checksums.txt && \
    sha256sum /lib/libwasmvm_muslc.$ARCH.a | grep $(cat /tmp/checksums.txt | grep libwasmvm_muslc.$ARCH | cut -d ' ' -f 1)

# Copy the remaining files
COPY . .

# Build mantrachaind binary
# build tag info: https://github.com/cosmos/wasmd/blob/master/README.md#supported-systems
RUN --mount=type=cache,target=/nonroot/.cache/go-build \
    --mount=type=cache,target=/nonroot/go/pkg/mod \
    LEDGER_ENABLED=true BUILD_TAGS='muslc osusergo' LINK_STATICALLY=true make build

# --------------------------------------------------------
# Runner
# --------------------------------------------------------

FROM alpine:$IMG_TAG
RUN apk add --no-cache build-base jq
RUN addgroup -g 1025 nonroot
RUN adduser -D nonroot -u 1025 -G nonroot
ARG IMG_TAG
COPY --from=mantra-builder /src/app/build/mantrachaind /usr/local/bin/
EXPOSE 26656 26657 1317 9090
USER nonroot

ENTRYPOINT ["mantrachaind", "start"]