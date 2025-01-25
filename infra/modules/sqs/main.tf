locals {
  common_tags = {
    Project     = "sampay"
    Environment = var.env
    ManagedBy   = "Terraform"
  }
}

resource "aws_sqs_queue" "worker_queue" {
  name                       = "sampay-${var.env}-worker-queue"
  visibility_timeout_seconds = 43200

  tags = local.common_tags
}

resource "aws_sqs_queue" "worker_dlq" {
  name = "sampay-${var.env}-worker-dlq"

  tags = local.common_tags
}

output "worker_arn" {
  value = aws_sqs_queue.worker_queue.arn
}

output "worker_url" {
  value = aws_sqs_queue.worker_queue.url
}

output "worker_dlq_arn" {
  value = aws_sqs_queue.worker_dlq.arn
}

output "worker_dlq_url" {
  value = aws_sqs_queue.worker_dlq.url
}
