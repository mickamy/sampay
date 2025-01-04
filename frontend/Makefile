.PHONY: \
	build \
	install \
	help \

build: ## Build the application
	npm run build

install: ## Install all dependencies
	npm install

gen: ## Generate protobuf/connect files
	rm -rf ./app/gen
	buf generate ../proto

help: ## Display a list of available Makefile targets with their descriptions
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
