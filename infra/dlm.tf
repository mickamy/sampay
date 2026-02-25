resource "aws_iam_role" "dlm" {
  name = "${local.name_prefix}-dlm-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "dlm.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Name = "${local.name_prefix}-dlm-role"
  }
}

resource "aws_iam_role_policy" "dlm" {
  name = "${local.name_prefix}-dlm-policy"
  role = aws_iam_role.dlm.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ec2:CreateSnapshot",
          "ec2:CreateTags",
          "ec2:DeleteSnapshot",
          "ec2:DescribeInstances",
          "ec2:DescribeVolumes",
          "ec2:DescribeSnapshots",
        ]
        Resource = "*"
      }
    ]
  })
}

resource "aws_dlm_lifecycle_policy" "ebs_backup" {
  description        = "EBS daily backup for ${local.name_prefix}"
  execution_role_arn = aws_iam_role.dlm.arn
  state              = "ENABLED"

  policy_details {
    resource_types = ["INSTANCE"]

    schedule {
      name = "Daily snapshot"

      create_rule {
        interval      = 24
        interval_unit = "HOURS"
        times         = ["16:00"] # 01:00 JST
      }

      retain_rule {
        count = 7
      }

      tags_to_add = {
        SnapshotCreator = "DLM"
      }

      copy_tags = true
    }

    target_tags = {
      Name = "${local.name_prefix}-ec2"
    }
  }

  tags = {
    Name = "${local.name_prefix}-dlm"
  }
}
