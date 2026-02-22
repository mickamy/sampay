.PHONY: \
	all \
	build \
	test \
	lint \
	lint-fix \
	ci \
	generate \
	db-create \
	db-migrate \
	db-drop \
	db-seed \
	db-setup \
	db-reset \
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

db-create: ## Create the database
	MODULE_ROOT=$$(pwd) go run ./cmd/db/create/main.go

db-migrate: ## Run database migrations
	MODULE_ROOT=$$(pwd) go run ./cmd/db/migrate/main.go

db-drop: ## Drop the database
	MODULE_ROOT=$$(pwd) go run ./cmd/db/drop/main.go

db-seed: ## Seed the database
	MODULE_ROOT=$$(pwd) go run ./cmd/db/seed/main.go

db-setup: db-create db-migrate db-seed ## Setup the database
	@echo "Database setup completed"

db-reset: db-drop db-setup ## Reset the database
	@echo "Database reset completed"

help: ## Display a list of available Makefile targets with their descriptions
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
