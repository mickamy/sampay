data "aws_caller_identity" "default" {}

data "aws_iam_openid_connect_provider" "github_actions_oidc" {
  url = "https://token.actions.githubusercontent.com"
}

resource "aws_iam_role" "github_actions_role" {
  name = "GitHubActionsRole-${data.aws_caller_identity.default.account_id}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Federated = data.aws_iam_openid_connect_provider.github_actions_oidc.arn
        },
        Action = "sts:AssumeRoleWithWebIdentity",
        Condition = {
          StringLike = {
            "token.actions.githubusercontent.com:sub" = "repo:${var.github_repo}:*"
          }
        }
      }
    ]
  })
}

resource "aws_iam_policy" "ec2_management_policy" {
  name        = "EC2ManagementPolicy-${data.aws_caller_identity.default.account_id}"
  description = "Policy to manage EC2 Security Groups for SSH access"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "ec2:AuthorizeSecurityGroupIngress",
          "ec2:RevokeSecurityGroupIngress",
          "ec2:DescribeSecurityGroups"
        ],
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ec2_management_attachment" {
  role       = aws_iam_role.github_actions_role.name
  policy_arn = aws_iam_policy.ec2_management_policy.arn
}
