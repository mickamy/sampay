locals {
  common_tags = {
    Project     = "sampay"
    Environment = var.env
    ManagedBy   = "Terraform"
  }
}

resource "null_resource" "check_keys" {
  provisioner "local-exec" {
    command = "test -f ${var.private_key_path} && test -f ${var.public_key_path} || echo 'Key files are missing'"
  }
}

resource "aws_ssm_parameter" "private_key" {
  depends_on = [null_resource.check_keys]
  name        = "/sampay/creds/${var.env}/private_key"
  description = "Private key for the ${var.env} environment"
  type        = "SecureString"
  value = file(var.private_key_path)
  key_id      = "alias/aws/ssm"
  tags        = local.common_tags
}

resource "aws_ssm_parameter" "public_key" {
  depends_on = [null_resource.check_keys]
  name        = "/sampay/creds/${var.env}/public_key"
  description = "Public key for the ${var.env} environment"
  type        = "String"
  value = file(var.public_key_path)
  tags        = local.common_tags
}

output "public_key" {
  value       = aws_ssm_parameter.public_key.value
  description = "The public key stored in SSM Parameter Store"
}
