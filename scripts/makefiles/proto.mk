###############################################################################
###                                  Proto                                  ###
###############################################################################

proto-help:
	@echo "proto subcommands"
	@echo ""
	@echo "Usage:"
	@echo "  make proto-[command]"
	@echo ""
	@echo "Available Commands:"
	@echo "  all        Run proto-format and proto-gen"
	@echo "  format     Format Protobuf files"
	@echo "  gen        Generate Protobuf files"

proto: proto-help
proto-all: proto-format proto-gen

PROTO_BUILDER_IMAGE=ghcr.io/cosmos/proto-builder:0.14.0
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(PROTO_BUILDER_IMAGE)

proto-all: proto-format proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@$(DOCKER) run --rm -u 0 -v $(CURDIR):/workspace --workdir /workspace $(PROTO_BUILDER_IMAGE) sh ./scripts/protocgen.sh

proto-format:
	@echo "Formatting Protobuf files"
	@$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/docker-build-proto \
		find ./proto -name "*.proto" -exec clang-format -i {} \;


SWAGGER_DIR=./swagger-proto
THIRD_PARTY_DIR=$(SWAGGER_DIR)/third_party

proto-download-deps:
	mkdir -p "$(THIRD_PARTY_DIR)/cosmos_tmp" && \
	cd "$(THIRD_PARTY_DIR)/cosmos_tmp" && \
	git clone --depth 1 --branch release/v0.50.x https://github.com/cosmos/cosmos-sdk.git . && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/cosmos_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/ibc_tmp" && \
	cd "$(THIRD_PARTY_DIR)/ibc_tmp" && \
	git clone --depth 1 https://github.com/cosmos/ibc-go.git . && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/ibc_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/osmosis_tmp" && \
	cd "$(THIRD_PARTY_DIR)/osmosis_tmp" && \
	git clone --depth 1 --branch v26.x https://github.com/osmosis-labs/osmosis.git . && \
	mkdir -p "../osmosis" && \
	mv ./proto/osmosis/tokenfactory ../osmosis
	rm -rf "$(THIRD_PARTY_DIR)/osmosis_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/feemarket_tmp" && \
	cd "$(THIRD_PARTY_DIR)/feemarket_tmp" && \
	git clone --depth 1 --branch $(shell grep -o 'feemarket v[0-9]\+\.[0-9]\+\.[0-9]\+' go.mod | awk '{print $$2}') https://github.com/skip-mev/feemarket.git . && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/feemarket_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/connect_tmp" && \
	cd "$(THIRD_PARTY_DIR)/connect_tmp" && \
	git clone --depth 1 https://github.com/skip-mev/connect.git . && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/connect_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/ibc_apps_tmp" && \
	cd "$(THIRD_PARTY_DIR)/ibc_apps_tmp" && \
	git clone --depth 1 https://github.com/cosmos/ibc-apps.git . && \
	mkdir -p "../ratelimit" && \
	mv ./modules/rate-limiting/proto/ratelimit ..
	rm -rf "$(THIRD_PARTY_DIR)/ibc_apps_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/cosmos_proto" && \
	curl -SSL https://raw.githubusercontent.com/cosmos/cosmos-proto/main/proto/cosmos_proto/cosmos.proto > "$(THIRD_PARTY_DIR)/cosmos_proto/cosmos.proto"

	mkdir -p "$(THIRD_PARTY_DIR)/gogoproto" && \
	curl -SSL https://raw.githubusercontent.com/cosmos/gogoproto/main/gogoproto/gogo.proto > "$(THIRD_PARTY_DIR)/gogoproto/gogo.proto"

	mkdir -p "$(THIRD_PARTY_DIR)/google/api" && \
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto > "$(THIRD_PARTY_DIR)/google/api/annotations.proto"
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto > "$(THIRD_PARTY_DIR)/google/api/http.proto"

	mkdir -p "$(THIRD_PARTY_DIR)/cosmos/ics23/v1" && \
	curl -sSL https://raw.githubusercontent.com/cosmos/ics23/master/proto/cosmos/ics23/v1/proofs.proto > "$(THIRD_PARTY_DIR)/cosmos/ics23/v1/proofs.proto"

docs:
	@echo
	@echo "=========== Generate Message ============"
	@echo
	@make proto-download-deps
	./scripts/generate-docs.sh

	statik -src=client/docs/static -dest=client/docs -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
        echo "\033[91mSwagger docs are out of sync!!!\033[0m";\
        exit 1;\
    else \
        echo "\033[92mSwagger docs are in sync\033[0m";\
    fi
	@echo
	@echo "=========== Generate Complete ============"
	@echo
.PHONY: docs