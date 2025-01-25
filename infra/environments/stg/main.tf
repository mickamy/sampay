terraform {
  backend "s3" {
    bucket       = "sampay-tf-backend"
    key          = "stg/state/terraform.tfstate"
    region       = "ap-northeast-1"
    use_lockfile = true
    encrypt      = true
  }
}

locals {
  env = "stg"
}

provider "aws" {
  region = var.aws_region
}

module "ec2" {
  source = "../../modules/ec2"

  env = local.env

  aws_region               = var.aws_region
  domain                   = var.domain
  instance_type            = var.instance_type
  public_key               = module.ssm.public_key
  s3_public_bucket_arn     = module.s3.public_bucket_arn
  sqs_worker_dlq_queue_arn = module.sqs.worker_dlq_arn
  sqs_worker_queue_arn     = module.sqs.worker_arn
  subnet_id                = module.vpc.public_subnet_id
  volume_size              = var.volume_size
  volume_type              = var.volume_type
  vpc_security_group_ids = [
    module.sg.sg_eic_id,
    module.sg.sg_ssh_id,
    module.sg.sg_web_id,
  ]

  depends_on = [
    module.s3,
    module.sg,
    module.sqs,
    module.ssm,
    module.vpc,
  ]
}

module "s3" {
  source = "../../modules/s3"
  env    = local.env

  geo_locations = var.geo_locations
}

module "sg" {
  source = "../../modules/sg"
  env    = local.env

  vpc_id = module.vpc.vpc_id
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
  db_timezone        = var.db_timezone
  frontend_base_url  = var.frontend_base_url
  private_key_path   = var.private_key_path
  public_key_path    = var.public_key_path
  redis_host         = var.redis_host
  redis_port         = var.redis_port
  sqs_worker_dlq_url = module.sqs.worker_url
  sqs_worker_url     = module.sqs.worker_dlq_url
}

module "vpc" {
  source = "../../modules/vpc"
  env    = local.env

  vpc_cidr = var.vpc_cidr
}
