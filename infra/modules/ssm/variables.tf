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

variable "db_reader_password" {
  description = "Database reader password"
  type        = string
}

variable "db_reader_user" {
  description = "Database reader user"
  type        = string
}

variable "db_timezone" {
  description = "Database timezone"
  type        = string
}

variable "db_writer_password" {
  description = "Database writer password"
  type        = string
}

variable "db_writer_user" {
  description = "Database writer user"
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

variable "private_key_path" {
  description = "Path to the private key"
  type        = string
}

variable "public_key_path" {
  description = "Path to the public key"
  type        = string
}

variable "redis_url" {
  description = "Redis URL"
  type        = string
}

variable "sqs_worker_dlq_url" {
  description = "SQS worker DLQ URL"
  type        = string
}

variable "sqs_worker_url" {
  description = "SQS worker URL"
  type        = string
}
