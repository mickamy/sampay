locals {
  is_prod     = var.environment == "prod"
  app_domain  = local.is_prod ? var.domain : "stg.${var.domain}"
  api_domain  = local.is_prod ? "api.${var.domain}" : "api.stg.${var.domain}"
  cdn_domain  = local.is_prod ? "cdn.${var.domain}" : "cdn.stg.${var.domain}"
  name_prefix = "sampay-${var.environment}"

  ecr_backend_url  = aws_ecr_repository.backend.repository_url
  ecr_frontend_url = aws_ecr_repository.frontend.repository_url
}
