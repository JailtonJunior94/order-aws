resource "aws_iam_role" "lambda_exec" {
  name = "example-lambda-exec"
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action = "sts:AssumeRole",
      Effect = "Allow",
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })

  tags = {
    Name        = "${var.prefix}-${var.environment}"
    Environment = var.environment
  }
}

resource "aws_iam_role_policy" "lambda_policy" {
  name = "example-lambda-policy"
  role = aws_iam_role.lambda_exec.id
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        Resource = "arn:aws:logs:*:*:*"
      }
    ]
  })
}

resource "aws_lambda_function" "main" {
  filename         = "../../bin/presign.zip"
  function_name    = "presign"
  role             = aws_iam_role.lambda_exec.arn
  handler          = "bootstrap"
  description      = "Lambda function to generate pre-signed URLs for S3"
  architectures    = ["x86_64"]
  runtime          = "provided.al2"
  timeout          = 30
  memory_size      = 128
  source_code_hash = filebase64sha256("../../bin/presign.zip")

  environment {
    variables = {
      REGION          = "us-east-1"
      BUCKET_NAME     = aws_s3_bucket.orders_bucket.bucket
      BUCKET_ENDPOINT = "http://s3.localhost.localstack.cloud:4566"
    }
  }

  tags = {
    Name        = "${var.prefix}-${var.environment}"
    Environment = var.environment
  }
}
