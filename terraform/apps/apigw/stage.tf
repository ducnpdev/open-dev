#
# Stage and Stage Settings
#

resource "aws_api_gateway_stage" "dev" {
  deployment_id = aws_api_gateway_deployment.tfdeployment.id
  rest_api_id   = aws_api_gateway_rest_api.tfapi.id
  stage_name    = "dev"
}

#resource "aws_api_gateway_stage" "uat" {
  # deployment_id = aws_api_gateway_deployment.example.id
  #  rest_api_id   = aws_api_gateway_rest_api.tfapi.id
  # stage_name    = "uat"
#}

# resource "aws_api_gateway_stage" "prod" {
  # deployment_id = aws_api_gateway_deployment.example.id
  # rest_api_id   = aws_api_gateway_rest_api.example.id
  # stage_name    = "prod"
#}




# resource "aws_api_gateway_method_settings" "example" {
#   rest_api_id = aws_api_gateway_rest_api.example.id
#   stage_name  = aws_api_gateway_stage.example.stage_name
#   method_path = "*/*"

#   settings {
#     metrics_enabled = true
#   }
# }