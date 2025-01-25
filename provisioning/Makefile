.PHONY: \
	stg \

stg: ## Provisioning staging server
	ANSIBLE_HOST=stg.sampay.link ANSIBLE_USER=ec2-user ANSIBLE_SSH_KEY=~/.ssh/sampay-stg ansible-playbook -i $$(PWD)/inventory/web.yaml playbook.yaml
