PROFILE ?= default
ENV ?= prod
TF = AWS_PROFILE=$(PROFILE) tofu

.PHONY: init plan apply destroy fmt validate workspace

init: ## Initialize the Terraform working directory
	$(TF) init

workspace: ## Ensure the correct Terraform workspace is selected or created
	$(TF) workspace select $(ENV) || $(TF) workspace new $(ENV)

plan: workspace ## Generate and show an execution plan for the Terraform configuration
	$(TF) plan -var-file=$(ENV).tfvars -var-file=.secret.$(ENV).tfvars

apply: workspace ## Apply the Terraform configuration to create or update infrastructure
	$(TF) apply -var-file=$(ENV).tfvars -var-file=.secret.$(ENV).tfvars

destroy: workspace ## Destroy the Terraform-managed infrastructure
	$(TF) destroy -var-file=$(ENV).tfvars -var-file=.secret.$(ENV).tfvars

fmt: ## Format the Terraform configuration files
	$(TF) fmt -recursive

validate: ## Validate the Terraform configuration
	$(TF) validate

help: ## Display a list of available Makefile targets with their descriptions
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'