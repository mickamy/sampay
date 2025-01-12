.PHONY: \
	build \
	install \
	tsc \
	check \
	ci \
	help \

build: ## Build the application
	npm run build

install: ## Install all dependencies
	npm install

tsc: ## Run the TypeScript compiler
	npm run typecheck

check: ## Run the TypeScript compiler and linter
	npm run check

ci: tsc check ## Run the CI pipeline
	@echo "CI pipeline passed"

help: ## Display a list of available Makefile targets with their descriptions
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
