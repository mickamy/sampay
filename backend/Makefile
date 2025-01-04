.PHONY: \
	install \
	install-bin \
	install-mod \
	gen \
	gen-pb \
	fmt \
	lint \

install: install-bin install-mod ## Install all dependencies
	@echo "All dependencies installed"

install-bin: ## Install go bin dependencies
	@echo "Installing dev tools"
	go install github.com/air-verse/air@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/kisielk/errcheck@latest

	@echo "Installing protobuf and Connect RPC tools"
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest

install-mod: ## Install go modules
	go mod tidy

gen: gen-pb ## Generate all files

gen-pb: ## Generate protobuf/connect files
	rm -rf ./gen
	buf generate ../proto

fmt: ## Run gofmt
	gofmt -s -w .

lint: ## Run linters (vet, staticcheck, and errcheck)
	go vet ./...
	staticcheck ./...
	errcheck ./...

help: ## Display a list of available Makefile targets with their descriptions
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
