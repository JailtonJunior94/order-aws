resource "aws_s3_bucket" "orders_bucket" {
  bucket        = "${var.environment}-orders-bucket"
  force_destroy = true
  
  tags = {
    Name        = "${var.prefix}-${var.environment}"
    Environment = var.environment
  }
}

resource "aws_s3_bucket_ownership_controls" "orders_bucket_ownership_controls" {
  bucket = aws_s3_bucket.orders_bucket.id

  rule {
    object_ownership = "BucketOwnerEnforced"
  }
}

resource "aws_s3_bucket_versioning" "orders_bucket_versioning" {
  bucket = aws_s3_bucket.orders_bucket.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_acl" "orders_bucket_acl" {
  bucket = aws_s3_bucket.orders_bucket.id

  acl        = "private"
  depends_on = [aws_s3_bucket.orders_bucket, aws_s3_bucket_ownership_controls.orders_bucket_ownership_controls]
}

resource "aws_s3_bucket_public_access_block" "orders_bucket_access_block" {
  bucket = aws_s3_bucket.orders_bucket.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_policy" "orders_bucket_policy" {
  bucket = aws_s3_bucket.orders_bucket.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect    = "Allow"
        Principal = "*"
        Action    = ["s3:GetObject", "s3:PutObject", "s3:DeleteObject"]
        Resource  = "${aws_s3_bucket.orders_bucket.arn}/*"
      }
    ]
  })
}

resource "aws_s3_bucket_notification" "orders_bucket_notification" {
  bucket = aws_s3_bucket.orders_bucket.id

  queue {
    queue_arn     = aws_sqs_queue.orders_queue.arn
    events        = ["s3:ObjectCreated:*"]
    filter_prefix = "" # opcional: pode filtrar por prefixo
    filter_suffix = "" # opcional: pode filtrar por sufixo
  }

  depends_on = [aws_s3_bucket_policy.orders_bucket_policy]
}
