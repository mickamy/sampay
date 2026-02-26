output "ec2_public_ip" {
  value       = aws_eip.main.public_ip
  description = "EC2 Elastic IP"
}

output "ecr_backend_url" {
  value       = local.ecr_backend_url
  description = "ECR backend repository URL"
}

output "ecr_frontend_url" {
  value       = local.ecr_frontend_url
  description = "ECR frontend repository URL"
}

output "s3_public_bucket" {
  value       = aws_s3_bucket.public.id
  description = "S3 public bucket name"
}

output "s3_private_bucket" {
  value       = aws_s3_bucket.private.id
  description = "S3 private bucket name"
}

output "cloudfront_distribution_domain" {
  value       = aws_cloudfront_distribution.cdn.domain_name
  description = "CloudFront distribution domain"
}

output "sqs_worker_url" {
  value       = aws_sqs_queue.worker.url
  description = "SQS worker queue URL"
}

output "app_domain" {
  value       = local.app_domain
  description = "Application domain"
}

output "api_domain" {
  value       = local.api_domain
  description = "API domain"
}

output "cdn_domain" {
  value       = local.cdn_domain
  description = "CDN domain"
}

output "sns_alerts_arn" {
  value       = aws_sns_topic.alerts.arn
  description = "SNS alerts topic ARN"
}

output "s3_logs_bucket" {
  value       = aws_s3_bucket.logs.id
  description = "S3 logs bucket name"
}

output "sm_app_secret_arn" {
  value       = aws_secretsmanager_secret.app.arn
  description = "Secrets Manager app secret ARN"
}

output "sm_db_secret_arn" {
  value       = aws_secretsmanager_secret.db.arn
  description = "Secrets Manager DB secret ARN"
}

output "sm_kvs_secret_arn" {
  value       = aws_secretsmanager_secret.kvs.arn
  description = "Secrets Manager KVS secret ARN"
}

output "ec2_security_group_id" {
  value       = aws_security_group.ec2.id
  description = "EC2 security group ID (set as EC2_SG_ID in GitHub environment secrets)"
}

output "github_actions_role_arn" {
  value       = aws_iam_role.github_actions.arn
  description = "IAM role ARN for GitHub Actions (set as AWS_ROLE_ARN in GitHub environment secrets)"
}
