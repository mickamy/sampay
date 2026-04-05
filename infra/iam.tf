resource "aws_iam_policy" "pi_secret_policy" {
  name        = "PiSecretsReadOnly"
  description = "Allow Pi to read secrets"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = "ecr:GetAuthorizationToken"
        Resource = "*"
      },
      {
        Effect = "Allow"
        Action = [
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
          "ecr:BatchCheckLayerAvailability",
        ]
        Resource = [
          aws_ecr_repository.backend.arn,
          aws_ecr_repository.frontend.arn,
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:ListBucket",
        ]
        Resource = [
          aws_s3_bucket.public.arn,
          "${aws_s3_bucket.public.arn}/*",
          aws_s3_bucket.private.arn,
          "${aws_s3_bucket.private.arn}/*",
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "sqs:SendMessage",
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes",
        ]
        Resource = [
          aws_sqs_queue.worker.arn,
          aws_sqs_queue.worker_dlq.arn,
        ]
      },
      {
        Action   = ["secretsmanager:GetSecretValue", "secretsmanager:DescribeSecret"]
        Effect   = "Allow"
        Resource = [
          aws_secretsmanager_secret.app.arn,
          aws_secretsmanager_secret.db.arn,
          aws_secretsmanager_secret.kvs.arn,
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "ses:SendEmail",
          "ses:SendRawEmail",
        ]
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_user" "pi_user" {
  name = "pi-server-user"
}

resource "aws_iam_user_policy_attachment" "pi_attach" {
  user       = aws_iam_user.pi_user.name
  policy_arn = aws_iam_policy.pi_secret_policy.arn
}