provider "aws" {
  region = "ap-northeast-1"
}

module "ses" {
  source  = "../../modules/ses"
  domain  = var.domain
  zone_id = var.zone_id
}
