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

###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=0.14.0
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./scripts/protocgen.sh

proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@$(protoImage) sh ./scripts/protoc-swagger-gen.sh

proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint --error-format=json

proto-check-breaking:
	@$(protoImage) buf breaking --against $(HTTPS_GIT)#branch=main

SWAGGER_DIR=./swagger-proto
THIRD_PARTY_DIR=$(SWAGGER_DIR)/third_party

proto-download-deps:
	mkdir -p "$(THIRD_PARTY_DIR)/wasm_tmp" && \
	cd "$(THIRD_PARTY_DIR)/wasm_tmp" && \
	git clone --depth 1 --branch $(shell grep -o 'wasmd v[0-9]\+\.[0-9]\+\.[0-9]\+' go.mod | awk '{print $$2}') https://github.com/CosmWasm/wasmd.git . && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/wasm_tmp"

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