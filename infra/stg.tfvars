environment   = "stg"
domain        = "sampay.link"
instance_type = "t4g.micro"
alert_email   = "system@sampay.link"

# Generate CI/CD key pair: ssh-keygen -t ed25519 -C "deploy@sampay" -f sampay_deploy_stg
# Set public key (sampay_deploy_stg.pub) here
# Store private key (sampay_deploy_stg) in GitHub Secrets as SSH_PRIVATE_KEY_STG
ssh_public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOQaey4RXTHsx9X4HWRH1pdCr+lUR4wjbTTJbN1f2AQJ deploy@sampay"
