terraform {
  backend "s3" {
    bucket       = "sampay-tf-backend"
    key          = "common/state/terraform.tfstate"
    region       = "ap-northeast-1"
    use_lockfile = true
    encrypt      = true
  }
}

provider "aws" {
  region = var.aws_region
}

module "ses" {
  source  = "../../modules/ses"
  domain  = var.domain
  zone_id = var.zone_id
}
