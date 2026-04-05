locals {
  app_domain  = var.domain
  api_domain  = "api.${var.domain}"
  cdn_domain  = "cdn.${var.domain}"
  name_prefix = "sampay"

  ecr_backend_url  = aws_ecr_repository.backend.repository_url
  ecr_frontend_url = aws_ecr_repository.frontend.repository_url
}
