variable "aws_region" {
  description = "AWS region to deploy the Lightsail instance"
  default     = "ap-northeast-1"
}

variable "domain" {
  description = "Domain name"
  type        = string
}

variable "geo_locations" {
  description = "Geo locations to allow access"
  type = list(string)
  default = ["JP"]
}

variable "ssh_private_key_path" {
  description = "Path to the private key"
  type        = string
}

variable "ssh_public_key_path" {
  description = "Path to the public key"
  type        = string
}
