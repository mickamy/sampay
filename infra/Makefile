.PHONY: \
	init-common \
	plan-common \
	apply-common \
	destroy-common \
	init-stg \
	plan-stg \
	apply-stg \
	destroy-stg \
	lint \

init-common:
	cd ./environments/common && terraform init

plan-common:
	cd ./environments/common && terraform plan

apply-common:
	cd ./environments/common && terraform apply

destroy-common:
	cd ./environments/common && terraform destroy

init-stg:
	cd ./environments/stg && terraform init

plan-stg:
	cd ./environments/stg && terraform plan

apply-stg:
	cd ./environments/stg && terraform apply

destroy-stg:
	cd ./environments/staging && terraform destroy

lint:
	tflint --fix --recursive
