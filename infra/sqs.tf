resource "aws_sqs_queue" "worker_dlq" {
  name                      = "${local.name_prefix}-worker-dlq"
  message_retention_seconds = 1209600 # 14 days
  sqs_managed_sse_enabled   = true

  tags = {
    Name = "${local.name_prefix}-worker-dlq"
  }
}

resource "aws_sqs_queue" "worker" {
  name                       = "${local.name_prefix}-worker"
  visibility_timeout_seconds = 300
  message_retention_seconds  = 345600 # 4 days
  sqs_managed_sse_enabled    = true

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.worker_dlq.arn
    maxReceiveCount     = 3
  })

  tags = {
    Name = "${local.name_prefix}-worker"
  }
}
