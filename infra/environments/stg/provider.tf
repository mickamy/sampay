terraform {
  required_version = ">= 1.10.5"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.84.0"
    }
    github = {
      source  = "integrations/github"
      version = "6.2.2"
    }
    null = {
      source  = "hashicorp/null"
      version = "3.2.3"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.6.3"
    }
    tls = {
      source  = "hashicorp/tls"
      version = ">= 4.0.0"
    }
  }
}
