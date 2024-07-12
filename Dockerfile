# syntax=docker/dockerfile:1

ARG GO_VERSION="1.22"
ARG FINAL_IMAGE="alpine:latest"
ARG BUILD_TAGS="netgo,ledger,muslc"

# --------------------------------------------------------
# Builder
# --------------------------------------------------------

FROM golang:${GO_VERSION}-alpine3.18 as builder

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

# TODO: Remove this part once no more private go modules
# ssh private key for downloading private go modules
# use with: docker build --secret id=goprivate,src=$HOME/.ssh/YOUR_PRIVATE_KEY .
RUN --mount=type=secret,id=goprivate \
  if [ -f /run/secrets/goprivate ]; then \
    mkdir -p ~/.ssh; \
    cat /run/secrets/goprivate > ~/.ssh/id_ed25519; \
    chmod 600 ~/.ssh/id_ed25519; \
    apk add openssh; \
    git config --global --add url."ssh://git@github.com/".insteadOf "https://github.com/"; \
    ssh-keyscan github.com >> /root/.ssh/known_hosts; \
  fi

# Download go dependencies
WORKDIR /mantrachain
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/root/go/pkg/mod \
  go mod download

# Cosmwasm - Download correct libwasmvm version
RUN ARCH=$(uname -m) && WASMVM_VERSION=$(go list -m github.com/CosmWasm/wasmvm | sed 's/.* //') && \
  wget https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/libwasmvm_muslc.$ARCH.a \
  -O /lib/libwasmvm_muslc.a && \
  # verify checksum
  wget https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/checksums.txt -O /tmp/checksums.txt && \
  sha256sum /lib/libwasmvm_muslc.a | grep $(cat /tmp/checksums.txt | grep libwasmvm_muslc.$ARCH | cut -d ' ' -f 1)

# Copy the remaining files
COPY . .

# Build mantrachaind binary
# build tag info: https://github.com/cosmos/wasmd/blob/master/README.md#supported-systems
RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/root/go/pkg/mod \
  LEDGER_ENABLED=true BUILD_TAGS=muslc LINK_STATICALLY=true make build

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