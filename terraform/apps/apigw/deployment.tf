resource "aws_api_gateway_deployment" "tfdeployment" {
  rest_api_id = aws_api_gateway_rest_api.tfapi.id

  # triggers = {
  #   redeployment = sha1(jsonencode(aws_api_gateway_rest_api.tfapi.body))
  # }

  lifecycle {
    create_before_destroy = true
  }
}

# resource "aws_api_gateway_stage" "example" {
#   deployment_id = aws_api_gateway_deployment.example.id
#   rest_api_id   = aws_api_gateway_rest_api.example.id
#   stage_name    = "example"
# }
