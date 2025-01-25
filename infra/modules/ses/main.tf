locals {
  common_tags = {
    Project     = "sampay"
    Environment = "common"
    ManagedBy   = "Terraform"
  }
}

########################################################################################################################
# SES
########################################################################################################################
resource "aws_ses_domain_identity" "domain" {
  domain = var.domain
}

resource "aws_route53_record" "ses_domain_verification" {
  zone_id = var.zone_id
  name    = "_amazonses.${aws_ses_domain_identity.domain.domain}"
  type    = "TXT"
  ttl     = 300
  records = [aws_ses_domain_identity.domain.verification_token]
}

resource "null_resource" "wait_for_ses_verification" {
  depends_on = [aws_ses_domain_identity.domain]

  provisioner "local-exec" {
    command = "aws ses verify-domain-identity --domain ${var.domain}"
  }
}

########################################################################################################################
# DKIM
########################################################################################################################
resource "aws_ses_domain_dkim" "dkim" {
  domain = aws_ses_domain_identity.domain.domain
}

resource "aws_route53_record" "dkim" {
  count   = 3
  zone_id = var.zone_id
  name    = "${aws_ses_domain_dkim.dkim.dkim_tokens[count.index]}.${aws_ses_domain_identity.domain.domain}"
  type    = "CNAME"
  ttl     = 300
  records = ["${aws_ses_domain_dkim.dkim.dkim_tokens[count.index]}.dkim.amazonses.com"]

  depends_on = [aws_ses_domain_dkim.dkim]
}

########################################################################################################################
# DMARC
########################################################################################################################
resource "aws_route53_record" "dmarc" {
  zone_id = var.zone_id
  name    = "_dmarc.${aws_ses_domain_identity.domain.domain}"
  type    = "TXT"
  ttl     = 300
  records = [
    "v=DMARC1; p=reject; rua=mailto:dmarc-reports@${aws_ses_domain_identity.domain.domain}; ruf=mailto:dmarc-reports@${aws_ses_domain_identity.domain.domain}; sp=none; adkim=s; aspf=s"
  ]
}

resource "aws_route53_record" "dmarc_reports_mx" {
  zone_id = var.zone_id
  name    = "dmarc-reports.${aws_ses_domain_identity.domain.domain}"
  type    = "MX"
  ttl     = 300
  records = ["10 inbound-smtp.us-east-1.amazonaws.com"]
}

resource "aws_s3_bucket" "dmarc_reports" {
  bucket = "dmarc-reports-${var.domain}"
  tags   = local.common_tags
}

resource "aws_s3_bucket_server_side_encryption_configuration" "dmarc_reports" {
  bucket = aws_s3_bucket.dmarc_reports.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "dmarc_reports_lifecycle" {
  bucket = aws_s3_bucket.dmarc_reports.id

  rule {
    id     = "DeleteOldReports"
    status = "Enabled"

    expiration {
      days = 90
    }
  }
}

resource "aws_s3_bucket_public_access_block" "dmarc_reports" {
  bucket = aws_s3_bucket.dmarc_reports.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_iam_role" "ses_dmarc_role" {
  name = "ses-dmarc-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ses.amazonaws.com"
        }
      }
    ]
  })

  tags = local.common_tags
}

data "aws_caller_identity" "default" {}

resource "aws_s3_bucket_policy" "dmarc_reports" {
  bucket = aws_s3_bucket.dmarc_reports.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Service = "ses.amazonaws.com"
        },
        Action   = "s3:PutObject",
        Resource = "${aws_s3_bucket.dmarc_reports.arn}/*",
        Condition = {
          StringEquals = {
            "AWS:SourceAccount" : data.aws_caller_identity.default.account_id
          }
        }
      }
    ]
  })
}


resource "aws_iam_role_policy" "ses_dmarc_policy" {
  name = "ses-dmarc-policy"
  role = aws_iam_role.ses_dmarc_role.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "s3:PutObject",
          "s3:GetBucketLocation"
        ],
        Resource = [
          aws_s3_bucket.dmarc_reports.arn,
          "${aws_s3_bucket.dmarc_reports.arn}/*"
        ]
      }
    ]
  })
}

resource "aws_ses_receipt_rule_set" "dmarc_ruleset" {
  rule_set_name = "dmarc-ruleset"
}

resource "aws_ses_receipt_rule" "store_dmarc" {
  name          = "store-dmarc-reports"
  rule_set_name = aws_ses_receipt_rule_set.dmarc_ruleset.rule_set_name
  recipients = ["dmarc-reports@${var.domain}"]
  enabled       = true
  scan_enabled  = true

  s3_action {
    bucket_name       = aws_s3_bucket.dmarc_reports.id
    object_key_prefix = "dmarc-reports/"
    position          = 1
  }
}

resource "aws_ses_active_receipt_rule_set" "dmarc_active" {
  rule_set_name = aws_ses_receipt_rule_set.dmarc_ruleset.rule_set_name
}

########################################################################################################################
# Email Identity
########################################################################################################################
resource "aws_ses_email_identity" "no_reply" {
  email = "no-reply@${aws_ses_domain_identity.domain.domain}"
}
