variable "domain" {
  description = "Domain name"
  type        = string
}

variable "aws_region" {
  description = "AWS region to deploy the Lightsail instance"
  type        = string
  default     = "ap-northeast-1"
}

variable "env" {
  description = "Environment name (e.g., stg, prod)"
  type        = string
}

variable "instance_type" {
  description = "Instance type"
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

variable "volume_size" {
  description = "Volume size in GB"
  type        = number
}

variable "volume_type" {
  description = "Volume type"
  type        = string
}

variable "vpc_security_group_ids" {
  description = "VPC security group IDs"
  type        = list(string)
}

variable "sqs_worker_dlq_queue_arn" {
  description = "SQS worker dead-letter queue ARN"
  type        = string
}

variable "sqs_worker_queue_arn" {
  description = "SQS worker queue ARN"
  type        = string
}

variable "subnet_id" {
  description = "Subnet ID"
  type        = string
}

variable "route53_record_ttl" {
  description = "TTL for the A record"
  type        = number
  default     = 300
}
