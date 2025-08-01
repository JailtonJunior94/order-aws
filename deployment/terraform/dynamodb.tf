resource "aws_dynamodb_table" "dynamodb_table_orders" {
  name           = "${var.enviroment}-orders"
  read_capacity  = 10
  write_capacity = 5
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }

  tags = {
    Name        = "${var.prefix}-orders-bucket"
    Environment = var.enviroment
  }
}

output "dynamodb_table_orders_arn" {
  description = "The ARN of the DynamoDB table for webhooks events"
  value       = aws_dynamodb_table.dynamodb_table_orders.arn
}
