environment   = "prod"
domain        = "sampay.link"
instance_type = "t4g.small"
alert_email   = "system@sampay.link"

# Generate CI/CD key pair: ssh-keygen -t ed25519 -C "deploy@sampay" -f sampay_deploy_prod
# Set public key (sampay_deploy_prod.pub) here
# Store private key (sampay_deploy_prod) in GitHub Secrets as SSH_PRIVATE_KEY_PROD
ssh_public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBjk2wi+IQSip+KvVZMZq6OdrNCS8IL8hmnJTpjmkOkQ deploy@sampay"
