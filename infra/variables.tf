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

variable "google_client_id" {
  type        = string
  description = "Google OAuth client ID"
  sensitive   = true
}

variable "google_client_secret" {
  type        = string
  description = "Google OAuth client secret"
  sensitive   = true
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
