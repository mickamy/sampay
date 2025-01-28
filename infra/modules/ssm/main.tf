locals {
  common_tags = {
    Project     = "sampay"
    Environment = var.env
    ManagedBy   = "Terraform"
  }

  creds = {
    "private_key" = {
      name        = "private_key"
      description = "SSH private key"
      file_path   = var.ssh_private_key_path
    }
    "public_key" = {
      name        = "public_key"
      description = "SSH public key"
      file_path   = var.ssh_public_key_path
    }
  }

  random_values = {
    "DB_ADMIN_PASSWORD"  = "Database admin password"
    "DB_NAME"            = "Database name"
    "DB_WRITER_USER"     = "Database writer user"
    "DB_WRITER_PASSWORD" = "Database writer password"
    "DB_READER_USER"     = "Database reader user"
    "DB_READER_PASSWORD" = "Database reader password"
    "JWT_SIGNING_SECRET" = "JWT signing secret"
    "REDIS_PASSWORD"     = "Redis password"
  }

  non_random_values = {
    "aws_region" = {
      name        = "AWS_REGION"
      description = "AWS region"
      value       = var.aws_region
    }
    "cloudfront_domain" = {
      name        = "CLOUDFRONT_DOMAIN"
      description = "CloudFront domain"
      value       = var.cloudfront_domain
    }
    "db_admin_user" = {
      name        = "DB_ADMIN_USER"
      description = "Database admin user"
      value       = var.db_admin_user
    }
    "db_host" = {
      name        = "DB_HOST"
      description = "Database host"
      value       = var.db_host
    }
    "db_port" = {
      name        = "DB_PORT"
      description = "Database port"
      value       = var.db_port
    }
    "db_timezone" = {
      name        = "DB_TIMEZONE"
      description = "Database timezone"
      value       = var.db_timezone
    }
    "frontend_base_url" = {
      name        = "FRONTEND_BASE_URL"
      description = "Frontend base URL"
      value       = var.frontend_base_url
    }
    "redis_host" = {
      name        = "REDIS_HOST"
      description = "Redis host"
      value       = var.redis_host
    }
    "redis_port" = {
      name        = "REDIS_PORT"
      description = "Redis port"
      value       = var.redis_port
    }
    "sqs_worker_dlq_url" = {
      name        = "SQS_WORKER_DLQ_URL"
      description = "SQS Worker DLQ URL"
      value       = var.sqs_worker_dlq_url
    }
    "sqs_worker_url" = {
      name        = "SQS_WORKER_URL"
      description = "SQS Worker URL"
      value       = var.sqs_worker_url
    }
  }
}

resource "null_resource" "check_keys" {
  provisioner "local-exec" {
    command = "test -f ${var.ssh_private_key_path} && test -f ${var.ssh_public_key_path} || echo 'Key files are missing'"
  }
}

resource "aws_ssm_parameter" "creds" {
  for_each    = local.creds
  name        = "/sampay/creds/${var.env}/${each.value.name}"
  description = each.value.description
  type        = "SecureString"
  value = file(each.value.file_path)
  tags        = local.common_tags

  depends_on = [null_resource.check_keys]
}

output "private_key" {
  value       = aws_ssm_parameter.creds["private_key"].value
  description = "The private key stored in SSM Parameter Store"
}

output "public_key" {
  value       = aws_ssm_parameter.creds["public_key"].value
  description = "The public key stored in SSM Parameter Store"
}

resource "random_password" "secure_values" {
  for_each         = local.random_values
  length           = 16
  special          = true
  override_special = "_%@"
  upper            = true
  lower            = true
  numeric          = true
}

resource "aws_ssm_parameter" "random_values" {
  for_each    = random_password.secure_values
  name        = "/sampay/app/${var.env}/${each.key}"
  description = local.random_values[each.key]
  type        = "SecureString"
  value       = each.value.result

  tags = local.common_tags
}

resource "aws_ssm_parameter" "non_random_values" {
  for_each    = local.non_random_values
  name        = "/sampay/app/${var.env}/${each.value.name}"
  description = each.value.description
  type        = "SecureString"
  value       = each.value.value
  tags        = local.common_tags
}

resource "github_actions_secret" "postgres_admin_password" {
  repository      = var.github_repo
  secret_name     = "POSTGRES_ADMIN_PASSWORD_${upper(var.env)}"
  plaintext_value = aws_ssm_parameter.random_values["DB_ADMIN_PASSWORD"].value
}

resource "github_actions_secret" "redis_password" {
  repository      = var.github_repo
  secret_name     = "REDIS_PASSWORD_${upper(var.env)}"
  plaintext_value = aws_ssm_parameter.random_values["REDIS_PASSWORD"].value
}
