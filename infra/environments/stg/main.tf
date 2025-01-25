locals {
  env = "stg"
}

provider "aws" {
  region = var.aws_region
}

module "lightsail" {
  source = "../../modules/lightsail"

  env = local.env

  aws_region               = var.aws_region
  domain                   = var.domain
  email_domain             = var.domain
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

  cloudfront_domain  = module.s3.cloudfront_domain
  db_admin_password  = var.db_admin_password
  db_admin_user      = var.db_admin_user
  db_host            = var.db_host
  db_name            = var.db_name
  db_port            = var.db_port
  db_reader_password = var.db_reader_password
  db_reader_user     = var.db_reader_user
  db_timezone        = var.db_timezone
  db_writer_password = var.db_writer_password
  db_writer_user     = var.db_writer_user
  frontend_base_url  = var.frontend_base_url
  private_key_path   = var.private_key_path
  public_key_path    = var.public_key_path
  redis_url          = var.redis_url
  sqs_worker_dlq_url = module.sqs.worker_url
  sqs_worker_url     = module.sqs.worker_dlq_url
}
