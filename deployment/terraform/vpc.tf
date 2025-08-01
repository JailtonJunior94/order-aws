resource "aws_vpc" "main_vpc" {
  cidr_block = "10.0.0.0/16"
  tags = {
    Name = "${var.prefix}-vpc"
  }
}

data "aws_availability_zones" "available_zones" {
  state = "available"
}

resource "aws_subnet" "subnets" {
  count                   = 2
  map_public_ip_on_launch = true
  vpc_id                  = aws_vpc.main_vpc.id
  cidr_block              = "10.0.${count.index}.0/24"
  availability_zone       = data.aws_availability_zones.available_zones.names[count.index]

  tags = {
    Name = "${var.prefix}-subnet-${count.index}"
  }
  depends_on = [aws_vpc.main_vpc]
}

output "subnets" {
  value = aws_subnet.subnets.*.id
}

output "vpc_id" {
  value = aws_vpc.main_vpc.id
}

resource "aws_internet_gateway" "internet_gateway" {
  vpc_id = aws_vpc.main_vpc.id
  tags = {
    Name = "${var.prefix}-internet-gateway"
  }
  depends_on = [aws_vpc.main_vpc]
}

resource "aws_route_table" "route_table" {
  vpc_id = aws_vpc.main_vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.internet_gateway.id
  }

  tags = {
    Name = "${var.prefix}-route-table"
  }
  depends_on = [
    aws_vpc.main_vpc,
    aws_internet_gateway.internet_gateway
  ]
}

resource "aws_route_table_association" "route_table_association" {
  count          = 2
  route_table_id = aws_route_table.route_table.id
  subnet_id      = aws_subnet.subnets.*.id[count.index]

  depends_on = [
    aws_subnet.subnets,
    aws_route_table.route_table
  ]
}
