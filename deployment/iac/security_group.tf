resource "aws_security_group" "security_group" {
  vpc_id = aws_vpc.main_vpc.id

  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_blocks     = ["0.0.0.0/0"]
    prefix_list_ids = []
  }

  depends_on = [aws_vpc.main_vpc]

  tags = {
    Name        = "${var.prefix}-${var.environment}"
    Environment = var.environment
  }
}
