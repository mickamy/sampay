variable "aws_region" {
  description = "AWS region to deploy the Lightsail instance"
  type        = string
}

variable "domain" {
  description = "Domain name"
  type        = string
}

variable "zone_id" {
  description = "Route53 zone ID"
  type        = string
}
