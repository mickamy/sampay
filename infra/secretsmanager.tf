resource "random_password" "postgres" {
  length  = 32
  special = false
}

resource "random_password" "db_writer" {
  length  = 32
  special = false
}

resource "random_password" "db_reader" {
  length  = 32
  special = false
}

resource "random_password" "kvs" {
  length  = 32
  special = false
}

resource "random_password" "session" {
  length  = 64
  special = false
}

resource "random_password" "jwt" {
  length  = 64
  special = false
}

resource "aws_secretsmanager_secret" "app" {
  name = "${local.name_prefix}/app"

  tags = {
    Name = "${local.name_prefix}-app-secret"
  }
}

resource "aws_secretsmanager_secret_version" "app" {
  secret_id = aws_secretsmanager_secret.app.id

  secret_string = jsonencode({
    # Auto-generated
    SESSION_SECRET     = random_password.session.result
    JWT_SIGNING_SECRET = random_password.jwt.result

    # DB connection
    DB_HOST            = "postgres"
    DB_PORT            = "5432"
    DB_NAME            = "sampay"
    DB_TIMEZONE        = "Asia/Tokyo"
    DB_ADMIN_USER      = "postgres"
    DB_ADMIN_PASSWORD  = random_password.postgres.result
    DB_WRITER_USER     = "sampay_writer"
    DB_WRITER_PASSWORD = random_password.db_writer.result
    DB_READER_USER     = "sampay_reader"
    DB_READER_PASSWORD = random_password.db_reader.result

    # KVS connection
    KVS_HOST     = "valkey"
    KVS_PORT     = "6379"
    KVS_USERNAME = ""
    KVS_PASSWORD = random_password.kvs.result

    # AWS resources (from TF outputs)
    S3_PUBLIC_BUCKET_NAME  = aws_s3_bucket.public.id
    S3_PRIVATE_BUCKET_NAME = aws_s3_bucket.private.id
    CLOUDFRONT_DOMAIN      = local.cdn_domain
    SQS_WORKER_URL         = aws_sqs_queue.worker.url
    SQS_WORKER_DLQ_URL     = aws_sqs_queue.worker_dlq.url

    # OAuth
    OAUTH_REDIRECT_URL   = "https://${local.app_domain}/oauth/callback"
    LINE_CHANNEL_ID      = var.line_channel_id
    LINE_CHANNEL_SECRET  = var.line_channel_secret

    # Misc
    EMAIL_FROM = var.email_from != "" ? var.email_from : "noreply@${var.domain}"
  })

  lifecycle {
    ignore_changes = [secret_string]
  }
}

resource "aws_secretsmanager_secret" "db" {
  name = "${local.name_prefix}/db"

  tags = {
    Name = "${local.name_prefix}-db-secret"
  }
}

resource "aws_secretsmanager_secret_version" "db" {
  secret_id = aws_secretsmanager_secret.db.id

  secret_string = jsonencode({
    POSTGRES_USER     = "postgres"
    POSTGRES_PASSWORD = random_password.postgres.result
  })

  lifecycle {
    ignore_changes = [secret_string]
  }
}

resource "aws_secretsmanager_secret" "kvs" {
  name = "${local.name_prefix}/kvs"

  tags = {
    Name = "${local.name_prefix}-kvs-secret"
  }
}

resource "aws_secretsmanager_secret_version" "kvs" {
  secret_id = aws_secretsmanager_secret.kvs.id

  secret_string = jsonencode({
    KVS_PASSWORD = random_password.kvs.result
  })

  lifecycle {
    ignore_changes = [secret_string]
  }
}
