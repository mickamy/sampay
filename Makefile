.PHONY: \
	build \
	up \
	up-d \
	down \
	down-v \
	install \
	help \

build: ## Build the docker images
	docker compose build

up: down ## Start the containers
	docker compose up

up-d: down ## Start the containers in detached mode
	docker compose up -d

down: ## Stop the containers
	docker compose down

down-v: ## Stop the containers and remove volumes
	docker compose down --volumes

install: ## Install all dependencies
	cd backend && make install
	cd frontend && make install

help: ## Display a list of available Makefile targets with their descriptions
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
