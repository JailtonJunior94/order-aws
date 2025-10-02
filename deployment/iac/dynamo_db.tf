resource "aws_dynamodb_table" "dynamodb_table_orders_sequence" {
  name           = "${var.environment}-orders-sequence"
  read_capacity  = 10
  write_capacity = 5

  hash_key       = "date"
  range_key      = "code"

  attribute {
    name = "date"
    type = "S"
  }

  attribute {
    name = "code"
    type = "S"
  }

  ttl {
    attribute_name = "expire_at"
    enabled        = true
  }

  tags = {
    Name        = "${var.prefix}-${var.environment}"
    Environment = var.environment
  }
}

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
    Name        = "${var.prefix}-${var.environment}"
    Environment = var.environment
  }
}
