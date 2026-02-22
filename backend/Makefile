.PHONY: \
	all \
	build \
	test \
	lint \
	lint-fix \
	ci \
	help

all: build ## Default target: build the application

build: build-api ## Build all binaries into ./bin

build-api: ## Build the API binary into ./bin
	go build -o ./bin/api ./cmd/api/main.go

test: ## Run all tests with race detector enabled
	go test -race ./...

lint: ## Run golangci-lint without fixing issues
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint is not installed"; \
		exit 1; \
	}
	golangci-lint run

lint-fix: ## Run golangci-lint with auto-fix enabled
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint is not installed"; \
		exit 1; \
	}
	golangci-lint run --fix

ci: lint test ## Run linting and tests for CI pipelines

generate: ## Run go generate for all packages
	find . \( -name '*_gen.go' -o -name '*.gen.go' \) -delete
	go generate ./...
	go tool injector generate --must ./...

help: ## Display a list of available Makefile targets with their descriptions
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
