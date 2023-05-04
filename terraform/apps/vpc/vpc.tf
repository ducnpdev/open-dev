# resource "aws_lb" "example" {
#   name               = "example"
#   internal           = true
#   load_balancer_type = "network"

#   subnet_mapping {
#     subnet_id = "12345"
#   }
# }

# resource "aws_api_gateway_vpc_link" "example" {
#   name        = "example"
#   description = "example description"
#   target_arns = [aws_lb.example.arn]
# }

resource "aws_default_vpc" "default" {
  tags = {
    Name = "Default VPC"
  }
}
resource "aws_subnet" "main" {
  vpc_id     = aws_default_vpc.default.id
  # cidr_block = "172.31.0.0/16"

  tags = {
    Name = "Main"
  }
}