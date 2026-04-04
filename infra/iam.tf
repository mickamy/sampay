resource "aws_iam_policy" "pi_secret_policy" {
  name        = "PiSecretsReadOnly"
  description = "Allow Pi to read secrets"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action   = ["secretsmanager:GetSecretValue", "secretsmanager:DescribeSecret"]
        Effect   = "Allow"
        Resource = [
          aws_secretsmanager_secret.app.arn,
          aws_secretsmanager_secret.db.arn,
          aws_secretsmanager_secret.kvs.arn,
        ]
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