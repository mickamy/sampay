PROFILE ?= default
TF = AWS_PROFILE=$(PROFILE) tofu
TF_VARS = -var-file=.secret.tfvars -var="aws_profile=$(PROFILE)"

.PHONY: init plan apply destroy fmt validate workspace

init: ## Initialize the Terraform working directory
	$(TF) init -backend-config="profile=$(PROFILE)"

plan: ## Generate and show an execution plan for the Terraform configuration
	$(TF) plan $(TF_VARS)

apply: ## Apply the Terraform configuration to create or update infrastructure
	$(TF) apply $(TF_VARS)

destroy: ## Destroy the Terraform-managed infrastructure
	$(TF) destroy $(TF_VARS)

fmt: ## Format the Terraform configuration files
	$(TF) fmt -recursive

validate: ## Validate the Terraform configuration
	$(TF) validate

help: ## Display a list of available Makefile targets with their descriptions
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
