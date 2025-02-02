locals {
  common_tags = {
    Project     = "sampay"
    Environment = var.env
    ManagedBy   = "Terraform"
  }
}

resource "aws_security_group" "eic" {
  name   = "sampay-${var.env}-eic"
  vpc_id = var.vpc_id

  ingress {
    description = "SSH access from EC2 Instance Connect"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["3.112.23.0/29"]
  }

  egress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.common_tags, {
    Name = "sampay-${var.env}-ec2-instance-connect"
  })
}

resource "aws_security_group" "ssh" {
  name        = "sampay-${var.env}-ssh"
  description = "Allow SSH traffic"
  vpc_id      = var.vpc_id

  ingress {
    description = "Allow SSH from trusted IP"
    from_port   = var.ssh_port
    to_port     = var.ssh_port
    protocol    = "tcp"
    cidr_blocks = [var.trusted_ip]
  }

  egress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "web" {
  name        = "sampay-${var.env}-web"
  description = "Allow HTTP and HTTPS traffic"
  vpc_id      = var.vpc_id

  ingress {
    description = "Allow HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "Allow HTTPS"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.common_tags, {
    Name = "sampay-${var.env}-web"
  })
}

output "sg_eic_id" {
  value = aws_security_group.eic.id
}

output "sg_ssh_id" {
  value = aws_security_group.ssh.id
}

output "sg_web_id" {
  value = aws_security_group.web.id
}

resource "github_actions_secret" "ssh_sg_id" {
  repository      = var.github_repo
  secret_name     = "SECURITY_GROUP_ID_${upper(var.env)}"
  plaintext_value = aws_security_group.ssh.id
}

resource "github_actions_secret" "ssh_port" {
  repository      = var.github_repo
  secret_name     = "SSH_PORT_${upper(var.env)}"
  plaintext_value = var.ssh_port
}
