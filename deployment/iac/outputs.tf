output "subnets" {
  value = aws_subnet.subnets.*.id
}

output "vpc_id" {
  value = aws_vpc.main_vpc.id
}

output "orders_queue_arn" {
  description = "The ARN of the main SQS queue"
  value       = aws_sqs_queue.orders_queue.arn
}

output "orders_queue_dlq_arn" {
  description = "The ARN of the DLQ SQS queue"
  value       = aws_sqs_queue.orders_queue_dlq.arn
}

output "dynamodb_table_orders_arn" {
  description = "The ARN of the DynamoDB table for webhooks events"
  value       = aws_dynamodb_table.dynamodb_table_orders.arn
}

output "security_group_name" {
  value = aws_security_group.security_group.name
}

output "security_group_id" {
  value = aws_security_group.security_group.id
}

output "api_gateway_id" {
  value = aws_api_gateway_rest_api.orders_api.id
}