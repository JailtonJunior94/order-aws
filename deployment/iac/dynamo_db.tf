resource "aws_dynamodb_table" "dynamodb_table_orders" {
  name           = "${var.environment}-orders"
  read_capacity  = 10
  write_capacity = 5
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }

  tags = {
    Name        = "${var.prefix}-orders-bucket"
    Environment = var.environment
  }
}
