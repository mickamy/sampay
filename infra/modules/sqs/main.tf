locals {
  common_tags = {
    Project     = "sampay"
    Environment = var.env
    ManagedBy   = "Terraform"
  }
}

resource "aws_sqs_queue" "worker_queue" {
  name                       = "${var.env}-sampay-worker-queue"
  visibility_timeout_seconds = 43200

  tags = local.common_tags
}

resource "aws_sqs_queue" "worker_dlq" {
  name = "${var.env}-sampay-worker-dlq"

  tags = local.common_tags
}

output "worker_arn" {
  value = aws_sqs_queue.worker_queue.arn
}

output "worker_dlq_arn" {
  value = aws_sqs_queue.worker_dlq.arn
}
