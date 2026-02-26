resource "aws_acm_certificate" "cdn" {
  provider = aws.us_east_1

  domain_name       = local.cdn_domain
  validation_method = "DNS"

  tags = {
    Name = "${local.name_prefix}-cdn-cert"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_acm_certificate_validation" "cdn" {
  provider = aws.us_east_1

  certificate_arn         = aws_acm_certificate.cdn.arn
  validation_record_fqdns = [for record in aws_route53_record.acm_validation : record.fqdn]
}
