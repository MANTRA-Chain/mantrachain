###############################################################################
###                                Linting                                  ###
###############################################################################
lint-help:
	@echo "lint subcommands"
	@echo ""
	@echo "Usage:"
	@echo "  make lint-[command]"
	@echo ""
	@echo "Available Commands:"
	@echo "  all                   Run all linters"
	@echo "  fix-typo              Run codespell to fix typos"
	@echo "  format                Run linters with auto-fix"
	@echo "  markdown              Run markdown linter with auto-fix"
	@echo "  mdlint                Run markdown linter"
	@echo "  setup-pre-commit      Set pre-commit git hook"
	@echo "  typo                  Run codespell to check typos"
lint: lint-help

golangci_version=v1.59.0

#? lint-install: Install golangci-lint
lint-install:
	@echo "--> Installing golangci-lint $(golangci_version)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)

lint-all:
	@echo "--> Running linter"
	$(MAKE) lint-install
	@golangci-lint run --timeout=10m
	@docker run -v $(PWD):/workdir ghcr.io/igorshubovych/markdownlint-cli:latest "**/*.md"

formatter-install:
	@echo "--> Installing gofumpt"
	@go install mvdan.cc/gofumpt@latest

lint-format:
	$(MAKE) lint-install formatter-install
	@golangci-lint run ./... --fix
	@gofumpt -l -w x/ app/ cmd/ internal/ tests/ 
	@docker run -v $(PWD):/workdir ghcr.io/igorshubovych/markdownlint-cli:latest "**/*.md" --fix

lint-mdlint:
	@echo "--> Running markdown linter"
	@docker run -v $(PWD):/workdir ghcr.io/igorshubovych/markdownlint-cli:latest "**/*.md"

lint-markdown:
	@docker run -v $(PWD):/workdir ghcr.io/igorshubovych/markdownlint-cli:latest "**/*.md" --fix

lint-typo:
	@codespell

lint-fix-typo:
	@codespell -w
