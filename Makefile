.PHONY: \
	build \
	up \
	up-d \
	down \
	down-v \
	install \
	ci \
	push-proto \
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

ci:
	cd backend && make ci
	cd frontend && make ci

push-proto:
	buf lint ./proto
	buf push ./proto
	cd backend && make update-mod
	cd frontend && npm update @buf/mickamy_sampay.bufbuild_es

help: ## Display a list of available Makefile targets with their descriptions
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
