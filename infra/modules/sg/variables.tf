variable "env" {
  description = "Environment name (e.g., stg, prod)"
  type        = string
}

variable "ssh_port" {
  description = "SSH port"
  type        = number
  default     = 22
}

variable "vpc_id" {
  description = "The ID of the VPC"
  type        = string
}
