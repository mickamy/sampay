variable "aws_region" {
  description = "AWS default region"
  type        = string
  default     = "ap-northeast-1"
}

variable "cloudfront_domain" {
  description = "CloudFront domain"
  type        = string
}

variable "db_admin_password" {
  description = "Database admin password"
  type        = string
}

variable "db_admin_user" {
  description = "Database admin user"
  type        = string
}

variable "db_host" {
  description = "Database host"
  type        = string
}

variable "db_name" {
  description = "Database name"
  type        = string
}

variable "db_port" {
  description = "Database port"
  type        = number
}

variable "db_timezone" {
  description = "Database timezone"
  type        = string
}

variable "env" {
  description = "Environment name (e.g., stg, prod)"
  type        = string
}

variable "frontend_base_url" {
  description = "Frontend base URL"
  type        = string
}

variable "ssh_private_key_path" {
  description = "Path to the private key"
  type        = string
}

variable "ssh_public_key_path" {
  description = "Path to the public key"
  type        = string
}

variable "redis_host" {
  description = "Redis host"
  type        = string
}

variable "redis_port" {
  description = "Redis port"
  type        = number
}

variable "sqs_worker_dlq_url" {
  description = "SQS worker DLQ URL"
  type        = string
}

variable "sqs_worker_url" {
  description = "SQS worker URL"
  type        = string
}
