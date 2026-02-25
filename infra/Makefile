PROFILE ?= default
ENV ?= prod
TF = AWS_PROFILE=$(PROFILE) tofu

.PHONY: init plan apply destroy fmt validate workspace

init:
	$(TF) init

workspace:
	$(TF) workspace select $(ENV) || $(TF) workspace new $(ENV)

plan: workspace
	$(TF) plan -var-file=$(ENV).tfvars -var-file=.secret.$(ENV).tfvars

apply: workspace
	$(TF) apply -var-file=$(ENV).tfvars -var-file=.secret.$(ENV).tfvars

destroy: workspace
	$(TF) destroy -var-file=$(ENV).tfvars -var-file=.secret.$(ENV).tfvars

fmt:
	$(TF) fmt -recursive

validate:
	$(TF) validate
