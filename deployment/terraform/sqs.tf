resource "aws_sqs_queue" "orders_queue" {
  name                       = "${var.enviroment}-orders"
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
    Environment = var.enviroment
  }
}

resource "aws_sqs_queue" "orders_queue_dlq" {
  name                       = "${var.enviroment}-orders-dlq"
  delay_seconds              = 0
  max_message_size           = 1024
  message_retention_seconds  = 86400
  visibility_timeout_seconds = 30

  tags = {
    Name        = "${var.prefix}-orders-queue-dlq"
    Environment = var.enviroment
  }
}

output "orders_queue_arn" {
  description = "The ARN of the main SQS queue"
  value       = aws_sqs_queue.orders_queue.arn
}

output "orders_queue_dlq_arn" {
  description = "The ARN of the DLQ SQS queue"
  value       = aws_sqs_queue.orders_queue_dlq.arn
}
