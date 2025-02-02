########################################################################################################################
# Deploy key
########################################################################################################################
resource "tls_private_key" "deploy_key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "github_repository_deploy_key" "main" {
  repository = var.github_repo
  title      = "Terraform Deploy Key"
  key        = tls_private_key.deploy_key.public_key_openssh
  read_only  = true
}

output "deploy_key_private" {
  value     = tls_private_key.deploy_key.private_key_pem
  sensitive = true
}

########################################################################################################################
# EC2
########################################################################################################################
locals {
  instance_name = "sampay-${var.env}"

  common_tags = {
    Project     = "sampay"
    Environment = var.env
    ManagedBy   = "Terraform"
  }
}

data "aws_caller_identity" "default" {}

resource "aws_key_pair" "ssh" {
  key_name   = "${local.instance_name}-key"
  public_key = var.ssh_public_key

  tags = local.common_tags
}

data "aws_ami" "amazon_linux_2023" {
  most_recent = true
  owners = ["amazon"]

  filter {
    name = "name"
    values = ["al2023-ami-*-kernel-6.1-x86_64"]
  }
}

resource "aws_instance" "main" {
  ami                         = data.aws_ami.amazon_linux_2023.id
  associate_public_ip_address = true
  iam_instance_profile        = aws_iam_instance_profile.ec2_instance_profile.name
  instance_type               = var.instance_type
  key_name                    = aws_key_pair.ssh.key_name
  subnet_id                   = var.subnet_id
  vpc_security_group_ids      = var.vpc_security_group_ids

  root_block_device {
    volume_size = var.volume_size
    volume_type = var.volume_type
  }

  user_data = templatefile("${path.module}/user_data.sh.tpl", {
    deploy_key : tls_private_key.deploy_key.private_key_openssh,
    ssh_port = var.ssh_port,
  })

  tags = merge(local.common_tags, {
    Name = local.instance_name
  })
}

resource "aws_eip" "web_eip" {
  instance = aws_instance.main.id

  tags = local.common_tags
}

output "public_ip" {
  value = aws_eip.web_eip.public_ip
}

data "aws_route53_zone" "main" {
  name = var.domain
}

locals {
  base_domain = var.env == "prod" ? var.domain : "${var.env}.${var.domain}"
  subdomains = {
    "api" = "api.${local.base_domain}"
    "web" = local.base_domain
  }
}

resource "aws_route53_record" "records" {
  for_each = local.subdomains

  zone_id = data.aws_route53_zone.main.zone_id
  name    = each.value
  type    = "A"
  ttl     = var.route53_record_ttl
  records = [aws_eip.web_eip.public_ip]

  lifecycle {
    create_before_destroy = true
  }
}

########################################################################################################################
# Github Actions secrets
########################################################################################################################
resource "github_actions_secret" "ec2_public_ip" {
  repository      = var.github_repo
  secret_name     = "EC2_PUBLIC_IP_${upper(var.env)}"
  plaintext_value = aws_instance.main.public_ip
}

resource "github_actions_secret" "ec2_ssh_key" {
  depends_on = [aws_key_pair.ssh]

  repository  = var.github_repo
  secret_name = "EC2_SSH_KEY_${upper(var.env)}"
  plaintext_value = base64encode(var.ssh_private_key)
}

########################################################################################################################
# IAM
########################################################################################################################
resource "aws_iam_role" "ec2_role" {
  name = "${local.instance_name}-ec2-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Service = "ec2.amazonaws.com"
        },
        Action = "sts:AssumeRole"
      }
    ]
  })

  tags = local.common_tags
}

resource "aws_iam_instance_profile" "ec2_instance_profile" {
  name = "${local.instance_name}-ec2-instance-profile"
  role = aws_iam_role.ec2_role.name

  tags = local.common_tags
}

resource "aws_iam_policy" "s3_access_policy" {
  name        = "sampay-${var.env}-s3-access"
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
  role       = aws_iam_role.ec2_role.name

  depends_on = [
    aws_iam_role.ec2_role,
    aws_iam_policy.s3_access_policy,
  ]
}

resource "aws_iam_policy" "ses_access_policy" {
  name        = "sampay-${var.env}-ses-access"
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
          "arn:aws:ses:${var.aws_region}:${data.aws_caller_identity.default.account_id}:identity/${var.email_from}",
        ],
      },
    ],
  })

  tags = local.common_tags
}

resource "aws_iam_role_policy_attachment" "attach_ses_access_policy" {
  policy_arn = aws_iam_policy.ses_access_policy.arn
  role       = aws_iam_role.ec2_role.name

  depends_on = [
    aws_iam_role.ec2_role,
    aws_iam_policy.ses_access_policy,
  ]
}

resource "aws_iam_policy" "sqs_access_policy" {
  name        = "sampay-${var.env}-sqs-access"
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
        ],
        Resource = [
          var.sqs_worker_dlq_queue_arn,
          var.sqs_worker_queue_arn,
        ],
      },
    ],
  })

  tags = local.common_tags
}

resource "aws_iam_role_policy_attachment" "attach_sqs_access_policy" {
  policy_arn = aws_iam_policy.sqs_access_policy.arn
  role       = aws_iam_role.ec2_role.name

  depends_on = [
    aws_iam_role.ec2_role,
    aws_iam_policy.sqs_access_policy,
  ]
}

resource "aws_iam_policy" "ssm_access_policy" {
  name        = "sampay-${var.env}-ssm-access"
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
  role       = aws_iam_role.ec2_role.name

  depends_on = [
    aws_iam_role.ec2_role,
    aws_iam_policy.ssm_access_policy,
  ]
}
