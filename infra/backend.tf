terraform {
  backend "s3" {
    bucket         = "sampay-tfstate"
    key            = "infra/terraform.tfstate"
    region         = "ap-northeast-1"
    dynamodb_table = "sampay-tfstate-lock"
    encrypt        = true
  }
}
