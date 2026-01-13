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

protoVer=0.15.2
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
	@echo "Downloading dependencies..."
	@go mod download
	@rm -rf "$(THIRD_PARTY_DIR)"
	@mkdir -p "$(THIRD_PARTY_DIR)"

	@echo "Copying cosmos-sdk proto..."
	@DIR=$$(go list -m -f '{{.Dir}}' github.com/cosmos/cosmos-sdk); \
	if [ -d "$$DIR/proto" ]; then \
		cp -r "$$DIR/proto"/* "$(THIRD_PARTY_DIR)"; \
		chmod -R 755 "$(THIRD_PARTY_DIR)"; \
	fi

	@echo "Copying evm proto..."
	@DIR=$$(go list -m -f '{{.Dir}}' github.com/cosmos/evm); \
	if [ -d "$$DIR/proto/cosmos" ]; then \
		mkdir -p "$(THIRD_PARTY_DIR)/cosmos"; \
		cp -r "$$DIR/proto/cosmos"/* "$(THIRD_PARTY_DIR)/cosmos"; \
		chmod -R 755 "$(THIRD_PARTY_DIR)"; \
	fi

	@echo "Copying wasmd proto..."
	@DIR=$$(go list -m -f '{{.Dir}}' github.com/CosmWasm/wasmd); \
	if [ -d "$$DIR/proto" ]; then \
		cp -r "$$DIR/proto"/* "$(THIRD_PARTY_DIR)"; \
		chmod -R 755 "$(THIRD_PARTY_DIR)"; \
	fi

	@echo "Copying ibc-go proto..."
	@DIR=$$(go list -m -f '{{.Dir}}' github.com/cosmos/ibc-go/v10); \
	if [ -d "$$DIR/proto" ]; then \
		cp -r "$$DIR/proto"/* "$(THIRD_PARTY_DIR)"; \
		chmod -R 755 "$(THIRD_PARTY_DIR)"; \
	fi

	@echo "Copying connect proto..."
	@DIR=$$(go list -m -f '{{.Dir}}' github.com/skip-mev/connect/v2); \
	if [ -d "$$DIR/proto" ]; then \
		cp -r "$$DIR/proto"/* "$(THIRD_PARTY_DIR)"; \
		chmod -R 755 "$(THIRD_PARTY_DIR)"; \
	fi

	@echo "Copying ibc-apps rate-limiting proto..."
	@DIR=$$(go list -m -f '{{.Dir}}' github.com/cosmos/ibc-apps/modules/rate-limiting/v10); \
	if [ -d "$$DIR/proto/ratelimit" ]; then \
		mkdir -p "$(THIRD_PARTY_DIR)/ratelimit"; \
		cp -r "$$DIR/proto/ratelimit"/* "$(THIRD_PARTY_DIR)/ratelimit"; \
		chmod -R 755 "$(THIRD_PARTY_DIR)"; \
	fi

	@# Remove buf.yaml and buf.lock from third_party to avoid module conflicts.
	@# We do not want to use buf.yaml from either cosmos-sdk or connect because
	@# third_party is a merged directory of multiple modules.
	@find "$(THIRD_PARTY_DIR)" -name "buf.yaml" -delete
	@find "$(THIRD_PARTY_DIR)" -name "buf.lock" -delete

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