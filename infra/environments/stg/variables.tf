variable "aws_region" {
  description = "AWS region to deploy the Lightsail instance"
  type        = string
  default     = "ap-northeast-1"
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
  default     = "Asia/Tokyo"
}

variable "domain" {
  description = "Domain name"
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

variable "github_repo_with_owner" {
  description = "GitHub repository with owner"
  type        = string
}

variable "github_token" {
  description = "GitHub token"
  type        = string
}

variable "geo_locations" {
  description = "Geo locations to allow access"
  type = list(string)
  default = ["JP"]
}

variable "google_client_id" {
  description = "Google client ID"
  type        = string
}

variable "google_client_secret" {
  description = "Google client secret"
  type        = string
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t2.medium"
}

variable "oauth_redirect_url" {
  description = "OAuth redirect URL"
  type = string
}

variable "redis_host" {
  description = "Redis host"
  type        = string
}

variable "redis_port" {
  description = "Redis port"
  type        = number
}

variable "ssh_port" {
  description = "SSH port"
  type        = number
}

variable "ssh_private_key_path" {
  description = "Path to the private key"
  type        = string
}

variable "ssh_public_key_path" {
  description = "Path to the public key"
  type        = string
}

variable "trusted_ip" {
  description = "Trusted IP address"
  type        = string
}

variable "volume_size" {
  description = "Volume size in GB"
  type        = number
  default     = 8
}

variable "volume_type" {
  description = "Volume type"
  type        = string
  default     = "gp2"
}

variable "vpc_cidr" {
  description = "VPC CIDR"
  type        = string
  default     = "10.0.0.0/16"
}
