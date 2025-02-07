variable "aws_region" {
  description = "AWS default region"
  type        = string
  default     = "ap-northeast-1"
}

variable "cloudfront_domain" {
  description = "CloudFront domain"
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

variable "email_from" {
  description = "Email from"
  type        = string
}

variable "frontend_base_url" {
  description = "Frontend base URL"
  type        = string
}

variable "github_repo" {
  description = "GitHub repository name"
  type        = string
}

variable "google_client_id" {
  description = "Google client ID"
  type        = string
}

variable "google_client_secret" {
  description = "Google client secret"
  type        = string
}

variable "oauth_redirect_url" {
  description = "OAuth redirect URL"
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

variable "s3_public_bucket_name" {
  description = "Public S3 bucket name"
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
