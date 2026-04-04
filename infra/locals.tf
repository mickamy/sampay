locals {
  is_prod     = var.environment == "prod"
  app_domain  = local.is_prod ? var.domain : "stg.${var.domain}"
  api_domain  = local.is_prod ? "api.${var.domain}" : "api.stg.${var.domain}"
  cdn_domain  = local.is_prod ? "cdn.${var.domain}" : "cdn.stg.${var.domain}"
  name_prefix = "sampay-${var.environment}"
}
