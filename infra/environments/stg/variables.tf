variable "aws_region" {
  description = "AWS region to deploy the Lightsail instance"
  type        = string
  default     = "ap-northeast-1"
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
  default     = "Asia/Tokyo"
}

variable "domain" {
  description = "Domain name"
  type        = string
}

variable "frontend_base_url" {
  description = "Frontend base URL"
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

variable "geo_locations" {
  description = "Geo locations to allow access"
  type = list(string)
  default = ["JP"]
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
