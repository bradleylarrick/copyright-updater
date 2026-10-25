# Configure the global provider settings
provider "aws" {
  region = "us-west-2"
}

# Define input variables for reusability
variable "environment" {
  type        = string
  default     = "production"
  description = "The deployment environment name"
}

# Create a local network block (VPC)
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true

  tags = {
    Name        = "main-vpc"
    Environment = var.environment
  }
}

# Deploy a web server instance inside the network
resource "aws_instance" "web_server" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t3.micro"
  
  # Reference the VPC ID created in the block above
  subnet_id = aws_vpc.main.default_subnet_id

  # Nested configuration block for storage
  root_block_device {
    volume_size           = 30
    volume_type           = "gp3"
    encrypted             = true
    delete_on_termination = true
  }

  tags = {
    Name = "${var.environment}-web-server"
  }
}

# Output the public IP after deployment
output "server_public_ip" {
  value       = aws_instance.web_server.public_ip
  description = "The public IP address of the main web server."
}

