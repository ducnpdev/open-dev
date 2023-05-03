resource "aws_api_gateway_resource" "tfdms" {
  rest_api_id = aws_api_gateway_rest_api.tfapi.id
  parent_id   = aws_api_gateway_rest_api.tfapi.root_resource_id
  path_part   = "dms"
}

resource "aws_api_gateway_resource" "tfdmsv1" {
  rest_api_id = aws_api_gateway_rest_api.tfapi.id
  parent_id   = aws_api_gateway_resource.tfdms.id
  path_part   = "v1"
}

resource "aws_api_gateway_resource" "tfdmsv1vendor" {
  rest_api_id = aws_api_gateway_rest_api.tfapi.id
  parent_id   = aws_api_gateway_resource.tfdmsv1.id
  path_part   = "vendor"
}

resource "aws_api_gateway_resource" "tfdmsv1vendorlogin" {
  rest_api_id = aws_api_gateway_rest_api.tfapi.id
  parent_id   = aws_api_gateway_resource.tfdmsv1vendor.id
  path_part   = "login"
}
