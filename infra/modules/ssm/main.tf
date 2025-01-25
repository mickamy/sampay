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

locals {
  creds = {
    "private_key" = {
      name        = "private_key"
      description = "SSH private key"
      file_path   = var.private_key_path
    }
    "public_key" = {
      name        = "public_key"
      description = "SSH public key"
      file_path   = var.public_key_path
    }
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

output "public_key" {
  value       = aws_ssm_parameter.creds["public_key"].value
  description = "The public key stored in SSM Parameter Store"
}

locals {
  app = {
    "aws_region" = {
      name        = "aws_region"
      description = "AWS region"
      value       = var.aws_region
    }
    "cloudfront_domain" = {
      name        = "cloudfront_domain"
      description = "CloudFront domain"
      value       = var.cloudfront_domain
    }
    "db_admin_password" = {
      name        = "db_admin_password"
      description = "Database admin password"
      value       = var.db_admin_password
    }
    "db_admin_user" = {
      name        = "db_admin_user"
      description = "Database admin user"
      value       = var.db_admin_user
    }
    "db_host" = {
      name        = "db_host"
      description = "Database host"
      value       = var.db_host
    }
    "db_name" = {
      name        = "db_name"
      description = "Database name"
      value       = var.db_name
    }
    "db_port" = {
      name        = "db_port"
      description = "Database port"
      value       = var.db_port
    }
    "db_reader_password" = {
      name        = "db_reader_password"
      description = "Database reader password"
      value       = var.db_reader_password
    }
    "db_reader_user" = {
      name        = "db_reader_user"
      description = "Database reader user"
      value       = var.db_reader_user
    }
    "db_timezone" = {
      name        = "db_timezone"
      description = "Database timezone"
      value       = var.db_timezone
    }
    "db_writer_password" = {
      name        = "db_writer_password"
      description = "Database writer password"
      value       = var.db_writer_password
    }
    "db_writer_user" = {
      name        = "db_writer_user"
      description = "Database writer user"
      value       = var.db_writer_user
    }
    "frontend_base_url" = {
      name        = "frontend_base_url"
      description = "Frontend base URL"
      value       = var.frontend_base_url
    }
    "redis_url" = {
      name        = "redis_url"
      description = "Redis URL"
      value       = var.redis_url
    }
    "sqs_worker_dlq_url" = {
      name        = "sqs_worker_dlq_url"
      description = "SQS Worker DLQ URL"
      value       = var.sqs_worker_dlq_url
    }
    "sqs_worker_url" = {
      name        = "sqs_worker_url"
      description = "SQS Worker URL"
      value       = var.sqs_worker_url
    }
  }
}

resource "aws_ssm_parameter" "app" {
  for_each    = local.app
  name        = "/sampay/app/${var.env}/${each.value.name}"
  description = each.value.description
  type        = "SecureString"
  value       = each.value.value
  tags        = local.common_tags
}
