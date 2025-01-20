variable "env" {
  description = "Environment name (e.g., stg, prod)"
  type        = string
}

variable "private_key_path" {
  description = "Path to the private key file"
  type        = string
}

variable "public_key_path" {
  description = "Path to the public key file"
  type        = string
}
