.PHONY: \
	build \
	install \
	install-build \
	install-dev \
	install-mod \
	db-create \
	db-migrate \
	db-seed \
	db-drop \
	db-prepare \
	db-reset \
	gen \
	wire \
	fmt \
	lint \
	test \
	help \

build: ## Build the application
	go build -o ./build/api ./cmd/api/main.go

install: install-build install-dev install-mod ## Install all dependencies
	@echo "All dependencies installed"

install-build: ## Install build dependencies
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

install-dev: ## Install development dependencies
	go install github.com/air-verse/air@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/kisielk/errcheck@latest
	go install github.com/google/wire/cmd/wire@latest
	go install go.uber.org/mock/mockgen@latest

install-mod: ## Install go modules
	go mod tidy

db-create: ## Create the database
	PACKAGE_ROOT=$$(pwd) go run ./cmd/db/create/main.go

db-migrate: ## Migrate the database
	PACKAGE_ROOT=$$(pwd) go run ./cmd/db/migrate/main.go

db-seed: ## Seed the database
	PACKAGE_ROOT=$$(pwd) go run ./cmd/db/seed/main.go

db-drop: ## Drop the database
	PACKAGE_ROOT=$$(pwd) go run ./cmd/db/drop/main.go

db-prepare: db-create db-migrate db-seed ## Prepare the database
	@echo "Database prepared"

db-reset: db-drop db-prepare ## Reset the database
	@echo "Database reset"

gen: ## Generate all files
	@echo "Generating files"
	find ./ -type d -name "mock_*" -exec rm -r {} +
	go generate ./...

wire: ## Generate wire files
	wire ./...

fmt: ## Run gofmt
	gofmt -s -w .

lint: ## Run linters (vet, staticcheck, and errcheck)
	go vet ./...
	staticcheck ./...
	errcheck ./...

test: ## Run tests
	PACKAGE_ROOT=$$(pwd) go test -v ./...

help: ## Display a list of available Makefile targets with their descriptions
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
