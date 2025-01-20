variable "api_domain" {
  description = "API domain name"
  type        = string
}

variable "aws_region" {
  description = "AWS region to deploy the Lightsail instance"
  default     = "ap-northeast-1"
}

variable "blueprint_id" {
  description = "Blueprint ID for Lightsail instance"
  default     = "amazon_linux_2"
}

variable "bundle_id" {
  description = "Instance plan (CPU, Memory, Disk size)"
  default     = "nano_2_0"
}

variable "email_domain" {
  description = "Email domain"
  type        = string
}

variable "env" {
  description = "Environment name (e.g., stg, prod)"
  type        = string
}

variable "public_key" {
  description = "Public key for SSH access"
  type        = string
}

variable "s3_public_bucket_arn" {
  description = "S3 public bucket ARN"
  type        = string
}

variable "sqs_worker_dlq_queue_arn" {
  description = "SQS worker dead-letter queue ARN"
  type        = string
}

variable "sqs_worker_queue_arn" {
  description = "SQS worker queue ARN"
  type        = string
}

variable "route53_record_ttl" {
  description = "TTL for the A record"
  type        = number
  default     = 300
}
