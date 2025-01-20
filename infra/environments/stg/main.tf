locals {
  env = "stg"
}

provider "aws" {
  region = var.aws_region
}

module "lightsail" {
  source = "../../modules/lightsail"

  env = local.env

  api_domain               = var.api_domain
  email_domain             = var.email_domain
  public_key               = module.ssm.public_key
  s3_public_bucket_arn     = module.s3.public_bucket_arn
  sqs_worker_dlq_queue_arn = module.sqs.worker_dlq_arn
  sqs_worker_queue_arn     = module.sqs.worker_arn

  depends_on = [
    module.s3,
    module.sqs,
    module.ssm,
  ]
}

module "s3" {
  source = "../../modules/s3"
  env    = local.env

  geo_locations = var.geo_locations
}

module "sqs" {
  source = "../../modules/sqs"
  env    = local.env
}

module "ssm" {
  source = "../../modules/ssm"
  env    = local.env

  private_key_path = var.ssh_private_key_path
  public_key_path  = var.ssh_public_key_path
}
