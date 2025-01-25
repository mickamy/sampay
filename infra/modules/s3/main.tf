locals {
  common_tags = {
    Project     = "sampay"
    Environment = var.env
    ManagedBy   = "Terraform"
  }
}

resource "aws_s3_bucket" "public" {
  bucket = "sampay-${var.env}-public"
  tags = local.common_tags
}

output "public_bucket_arn" {
  value = aws_s3_bucket.public.arn
}

# Enable server-side encryption
resource "aws_s3_bucket_server_side_encryption_configuration" "public" {
  bucket = aws_s3_bucket.public.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# Enable public access block
resource "aws_s3_bucket_public_access_block" "public_access_block" {
  bucket = aws_s3_bucket.public.id

  block_public_acls       = true
  ignore_public_acls      = true
  block_public_policy     = true
  restrict_public_buckets = true
}

# Create an OAI for the CloudFront distribution
resource "aws_cloudfront_origin_access_identity" "oai" {
  comment = "OAI for ${var.env}-sampay-public"
}

# Attach a bucket policy to the public bucket
resource "aws_s3_bucket_policy" "public_policy" {
  bucket = aws_s3_bucket.public.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          AWS = aws_cloudfront_origin_access_identity.oai.iam_arn
        },
        Action = "s3:GetObject",
        Resource = [
          "${aws_s3_bucket.public.arn}/*"
        ]
      },
      {
        Effect    = "Deny",
        Principal = "*",
        Action    = "s3:GetObject",
        Resource = [
          "${aws_s3_bucket.public.arn}/*"
        ],
        Condition = {
          StringNotEquals = {
            "aws:PrincipalArn" = aws_cloudfront_origin_access_identity.oai.iam_arn
          }
        }
      }
    ]
  })

  depends_on = [
    aws_cloudfront_origin_access_identity.oai
  ]
}

data "aws_cloudfront_cache_policy" "cache_policy" {
  name = "Managed-CachingOptimized"
}

# Create a CloudFront distribution
resource "aws_cloudfront_distribution" "cdn" {
  origin {
    domain_name = aws_s3_bucket.public.bucket_regional_domain_name
    origin_id   = aws_s3_bucket.public.bucket

    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.oai.cloudfront_access_identity_path
    }
  }

  enabled = true

  default_cache_behavior {
    allowed_methods = ["GET", "HEAD"]
    cached_methods = ["GET", "HEAD"]
    target_origin_id       = aws_s3_bucket.public.bucket
    viewer_protocol_policy = "redirect-to-https"
    cache_policy_id        = data.aws_cloudfront_cache_policy.cache_policy.id
  }

  price_class = "PriceClass_200"

  restrictions {
    geo_restriction {
      restriction_type = "whitelist"
      locations        = var.geo_locations
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }

  tags = local.common_tags
}

output "cloudfront_domain" {
  value = aws_cloudfront_distribution.cdn.domain_name
}
