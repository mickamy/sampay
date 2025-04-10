.PHONY: \
	init-backend \
	init-common \
	plan-common \
	apply-common \
	destroy-common \
	init-stg \
	init-prod \
	plan-stg \
	plan-prod \
	apply-stg \
	apply-prod \
	destroy-stg \
	destroy-prod \
	lint-init \
	lint \
	upgrade \

init-backend:
	set -e
	aws s3api create-bucket \
		--bucket $(BUCKET_NAME) \
		--region $(AWS_REGION) \
		--create-bucket-configuration LocationConstraint=$(AWS_REGION)

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

init-prod:
	cd ./environments/prod && terraform init

plan-stg:
	cd ./environments/stg && terraform plan

plan-prod:
	cd ./environments/prod && terraform plan

apply-stg:
	cd ./environments/stg && terraform apply

apply-prod:
	cd ./environments/prod && terraform apply

destroy-stg:
	cd ./environments/stg && \
	bash -c 'while read -r resource; do \
		terraform destroy -target="$${resource}" -auto-approve; \
	done < <( \
		terraform state list | \
		grep -vE "^data\\." | \
		grep -vE "module.ssm.aws_ssm_parameter.random_values|module.ssm.random_password.secure_values|module.s3.aws_cloudfront_distribution.cdn" \
	)'

destroy-prod:
	cd ./environments/prod && \
	bash -c 'while read -r resource; do \
		terraform destroy -target="$${resource}" -auto-approve; \
	done < <( \
		terraform state list | \
		grep -vE "^data\\." | \
		grep -vE "module.ssm.aws_ssm_parameter.random_values|module.ssm.random_password.secure_values|module.s3.aws_cloudfront_distribution.cdn" \
	)'


reapply-stg: destroy-stg apply-stg
	@echo "Reapply completed"

reapply-ec2-stg:
	cd ./environments/stg && terraform destroy -target=module.ec2
	cd ./environments/stg && terraform apply -target=module.ec2

lint-init:
	tflint --init

lint:
	tflint --fix --recursive

upgrade:
	cd ./environments/common && terraform init -upgrade
	cd ./environments/stg && terraform init -upgrade
