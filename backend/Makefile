.PHONY: \
	build \
	install \
	install-build \
	install-dev \
	install-mod \
	update-mod \
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
	ci \
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
	go install github.com/a-h/templ/cmd/templ@latest
	go install golang.org/x/tools/cmd/stringer@latest

install-mod: ## Install go modules
	go mod tidy

update-mod: ## Update go modules
	go get -u buf.build/gen/go/mickamy/sampay/bufbuild/connect-go@latest
	go get -u buf.build/gen/go/mickamy/sampay/protocolbuffers/go@latest
	go get -u ./...
	go mod tidy
	go generate ./... # Run generate after updating modules for some packages are dropped on go mod tidy

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
	templ generate
	gofmt -s -w .

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

ci: install-mod gen fmt lint test ## Run all CI checks
	@echo "CI pipeline passed"

help: ## Display a list of available Makefile targets with their descriptions
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
