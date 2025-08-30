resource "aws_sqs_queue" "orders_queue" {
  name                       = "${var.environment}-orders"
  delay_seconds              = 0
  max_message_size           = 1024
  message_retention_seconds  = 345600
  visibility_timeout_seconds = 30
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.orders_queue_dlq.arn
    maxReceiveCount     = 10
  })

  tags = {
    Name        = "${var.prefix}-orders-queue"
    Environment = var.environment
  }
}

resource "aws_sqs_queue" "orders_queue_dlq" {
  name                       = "${var.environment}-orders-dlq"
  delay_seconds              = 0
  max_message_size           = 1024
  message_retention_seconds  = 86400
  visibility_timeout_seconds = 30

  tags = {
    Name        = "${var.prefix}-orders-queue-dlq"
    Environment = var.environment
  }
}
