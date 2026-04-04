variable "aws_profile" {
  type        = string
  description = "AWS CLI profile name (null for CI/CD)"
  default     = null
}

variable "environment" {
  type        = string
  description = "Environment name (stg or prod)"

  validation {
    condition     = contains(["stg", "prod"], var.environment)
    error_message = "environment must be 'stg' or 'prod'"
  }
}

variable "domain" {
  type        = string
  description = "Base domain name"
  default     = "sampay.link"
}

variable "instance_type" {
  type        = string
  description = "EC2 instance type"
  default     = "t4g.small"
}

variable "ssh_public_key" {
  type        = string
  description = "SSH public key for deploy user"
}

variable "ssh_allowed_cidrs" {
  type        = list(string)
  description = "CIDR blocks allowed for SSH access"
  sensitive   = true
}

variable "alert_email" {
  type        = string
  description = "Email address for alarm notifications"
  default     = ""
}

variable "line_channel_id" {
  type        = string
  description = "LINE channel ID"
  sensitive   = true
}

variable "line_channel_secret" {
  type        = string
  description = "LINE channel secret"
  sensitive   = true
}

variable "email_from" {
  type        = string
  description = "Sender email address"
  default     = ""
}

variable "github_repo" {
  type        = string
  description = "GitHub repository (owner/repo)"
  default     = "mickamy/sampay"
}
