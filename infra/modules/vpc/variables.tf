variable "env" {
  description = "Environment name (e.g., stg, prod)"
  type        = string
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
}
