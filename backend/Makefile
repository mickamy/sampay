.PHONY: \
	install \
	install-bin \
	install-mod \
	fmt \
	lint \

install: install-bin install-mod ## Install all dependencies
	@echo "All dependencies installed"

install-bin: ## Install go bin dependencies
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/kisielk/errcheck@latest
	go install github.com/air-verse/air@latest

install-mod: ## Install go modules
	go mod tidy

fmt: ## Run gofmt
	gofmt -s -w .

lint: ## Run linters (vet, staticcheck, and errcheck)
	go vet ./...
	staticcheck ./...
	errcheck ./...

help: ## Display a list of available Makefile targets with their descriptions
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
