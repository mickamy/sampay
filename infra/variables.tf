variable "aws_profile" {
  type        = string
  description = "AWS CLI profile name (null for CI/CD)"
  default     = null
}

variable "domain" {
  type        = string
  description = "Base domain name"
  default     = "sampay.link"
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
