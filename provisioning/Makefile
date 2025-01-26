.PHONY: \
	stg \

stg: mkdir_secrets touch_secrets ## Provisioning staging server
	@if [ -z "$(POSTGRES_PASSWORD)" ]; then \
		POSTGRES_PASSWORD=$$(openssl rand -base64 16); \
		echo "Generated POSTGRES_PASSWORD: $$POSTGRES_PASSWORD"; \
		truncate -s 0 ./secrets/postgres_password.stg; \
		echo $$POSTGRES_PASSWORD > ./secrets/postgres_password.stg; \
		echo "Saved POSTGRES_PASSWORD to ./secrets/postgres_password.stg"; \
		ANSIBLE_HOST=stg.sampay.link \
		ANSIBLE_USER=ec2-user \
		ANSIBLE_SSH_KEY=~/.ssh/sampay-stg \
		POSTGRES_PASSWORD=$$POSTGRES_PASSWORD \
		ansible-playbook -i $$(PWD)/inventory/web.yaml playbook.yaml \
		--vault-password-file ./secrets/ansible_vault_pass; \
	else \
		echo "Using provided POSTGRES_PASSWORD: $(POSTGRES_PASSWORD)"; \
		truncate -s 0 ./secrets/postgres_password.stg; \
		echo $(POSTGRES_PASSWORD) > ./secrets/postgres_password.stg; \
		echo "Saved POSTGRES_PASSWORD to ./secrets/postgres_password.stg"; \
		ANSIBLE_HOST=stg.sampay.link \
		ANSIBLE_USER=ec2-user \
		ANSIBLE_SSH_KEY=~/.ssh/sampay-stg \
		POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		ansible-playbook -i $$(PWD)/inventory/web.yaml playbook.yaml \
		--vault-password-file ./secrets/ansible_vault_pass; \
	fi

mkdir_secrets:
	@echo "Creating secrets directory..."
	@mkdir -p ./secrets

touch_secrets:
	@echo "Creating secrets files..."
	@touch ./secrets/postgres_password
