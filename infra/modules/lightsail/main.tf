locals {
  instance_name = "sampay-${var.env}"

  common_tags = {
    Project     = "sampay"
    Environment = var.env
    ManagedBy   = "Terraform"
  }
}

data "aws_caller_identity" "default" {}

resource "aws_lightsail_key_pair" "sampay" {
  name       = "${local.instance_name}-key"
  public_key = var.public_key

  tags = local.common_tags
}

resource "aws_lightsail_instance" "web" {
  availability_zone = "${var.aws_region}a"
  blueprint_id      = var.blueprint_id
  bundle_id         = var.bundle_id
  key_pair_name     = aws_lightsail_key_pair.sampay.name
  name              = local.instance_name

  user_data = templatefile("${path.module}/user_data.sh.tpl", {
    aws_region : var.aws_region,
  })

  depends_on = [
    aws_lightsail_key_pair.sampay,
  ]

  tags = local.common_tags
}

resource "aws_lightsail_static_ip" "static_ip" {
  name = "${local.instance_name}-static-ip"
}

resource "aws_lightsail_static_ip_attachment" "attach_static_ip" {
  static_ip_name = aws_lightsail_static_ip.static_ip.name
  instance_name  = aws_lightsail_instance.web.name

  depends_on = [
    aws_lightsail_instance.web,
    aws_lightsail_static_ip.static_ip,
  ]
}

data "aws_route53_zone" "main" {
  name = "sampay.link"
}

resource "aws_route53_record" "public_record" {
  zone_id = data.aws_route53_zone.main.zone_id
  name    = "api.${var.domain}"
  type    = "A"
  ttl     = var.route53_record_ttl
  records = [aws_lightsail_static_ip.static_ip.ip_address]

  lifecycle {
    create_before_destroy = true
  }
}

########################################################################################################################
# IAM
########################################################################################################################
resource "aws_iam_role" "lightsail_role" {
  name = "lightsail-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Service = "lightsail.amazonaws.com"
        },
        Action = "sts:AssumeRole"
      }
    ]
  })

  tags = local.common_tags
}

resource "aws_iam_instance_profile" "lightsail_instance_profile" {
  name = "lightsail-instance-profile"
  role = aws_iam_role.lightsail_role.name

  tags = local.common_tags
}

resource "aws_iam_policy" "s3_access_policy" {
  name        = "sampay-s3-access-${var.env}"
  description = "Allows read and write actions for S3"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:ListBucket",
          "s3:GetBucketLocation",
        ],
        Resource = [
          var.s3_public_bucket_arn,
          "${var.s3_public_bucket_arn}/*",
        ],
      },
    ],
  })

  tags = local.common_tags
}

resource "aws_iam_role_policy_attachment" "attach_s3_access_policy" {
  policy_arn = aws_iam_policy.s3_access_policy.arn
  role       = aws_iam_role.lightsail_role.name

  depends_on = [
    aws_iam_role.lightsail_role,
    aws_iam_policy.s3_access_policy,
  ]
}

resource "aws_iam_policy" "ses_access_policy" {
  name        = "sampay-ses-access-${var.env}"
  description = "Allows send email actions for SES"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "ses:SendEmail",
        ],
        Resource = [
          "arn:aws:ses:${var.aws_region}:${data.aws_caller_identity.default.account_id}:identity/${var.email_domain}",
        ],
      },
    ],
  })

  tags = local.common_tags
}

resource "aws_iam_role_policy_attachment" "attach_ses_access_policy" {
  policy_arn = aws_iam_policy.ses_access_policy.arn
  role       = aws_iam_role.lightsail_role.name

  depends_on = [
    aws_iam_role.lightsail_role,
    aws_iam_policy.ses_access_policy,
  ]
}

resource "aws_iam_policy" "sqs_access_policy" {
  name        = "sampay-sqs-access-${var.env}"
  description = "Allows read and write actions for SQS"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "sqs:SendMessage",
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes",
          "sqs:ListQueues",
          "sqs:CreateQueue",
          "sqs:DeleteQueue",
        ],
        Resource = [
          var.sqs_worker_queue_arn,
          var.sqs_worker_dlq_queue_arn,
        ],
      },
    ],
  })

  tags = local.common_tags
}

resource "aws_iam_role_policy_attachment" "attach_sqs_access_policy" {
  policy_arn = aws_iam_policy.sqs_access_policy.arn
  role       = aws_iam_role.lightsail_role.name

  depends_on = [
    aws_iam_role.lightsail_role,
    aws_iam_policy.sqs_access_policy,
  ]
}

resource "aws_iam_policy" "ssm_access_policy" {
  name        = "sampay-ssm-access-${var.env}"
  description = "Allows read and write actions for SSM"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "ssm:GetParameter",
          "ssm:GetParameters",
          "ssm:GetParametersByPath",
        ],
        Resource = "arn:aws:ssm:${var.aws_region}:${data.aws_caller_identity.default.account_id}:parameter/sampay/app/${var.env}/*",
      },
    ],
  })

  tags = local.common_tags
}

resource "aws_iam_role_policy_attachment" "attach_ssm_access_policy" {
  policy_arn = aws_iam_policy.ssm_access_policy.arn
  role       = aws_iam_role.lightsail_role.name

  depends_on = [
    aws_iam_role.lightsail_role,
    aws_iam_policy.ssm_access_policy,
  ]
}
