locals {
  common_tags = {
    Project     = "sampay"
    Environment = var.env
    ManagedBy   = "Terraform"
  }
}

resource "aws_vpc" "main" {
  cidr_block = var.vpc_cidr

  tags = merge(local.common_tags, {
    Name = "sampay-${var.env}-vpc"
  })
}

resource "aws_subnet" "public" {
  vpc_id = aws_vpc.main.id
  cidr_block = cidrsubnet(var.vpc_cidr, 8, 1)

  tags = merge(local.common_tags, {
    Name = "sampay-${var.env}-subnet"
  })
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = merge(local.common_tags, {
    Name = "sampay-${var.env}-subnet"
  })
}

resource "aws_route" "default_route" {
  route_table_id         = aws_vpc.main.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.main.id
}


output "vpc_id" {
  value = aws_vpc.main.id
}

output "public_subnet_id" {
  value = aws_subnet.public.id
}
