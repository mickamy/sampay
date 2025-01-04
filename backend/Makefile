.PHONY: \
	build \
	install \
	install-build \
	install-dev \
	install-mod \
	gen \
	fmt \
	lint \
	help \

build: ## Build the application
	go build -o ./build/api ./cmd/api/main.go

install: install-build install-dev install-mod ## Install all dependencies
	@echo "All dependencies installed"

install-build: ## Install build dependencies
	go install github.com/bufbuild/buf/cmd/buf@latest

install-dev: ## Install development dependencies
	go install github.com/air-verse/air@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/kisielk/errcheck@latest

install-mod: ## Install go modules
	go mod tidy

gen: ## Generate all files
	@echo "Generating files"
	go generate ./...

fmt: ## Run gofmt
	gofmt -s -w .

lint: ## Run linters (vet, staticcheck, and errcheck)
	go vet ./...
	staticcheck ./...
	errcheck ./...

help: ## Display a list of available Makefile targets with their descriptions
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
