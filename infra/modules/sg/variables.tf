variable "env" {
  description = "Environment name (e.g., stg, prod)"
  type        = string
}

variable "github_repo" {
  description = "GitHub repository name"
  type        = string
}

variable "ssh_port" {
  description = "SSH port"
  type        = number
}

variable "vpc_id" {
  description = "The ID of the VPC"
  type        = string
}
