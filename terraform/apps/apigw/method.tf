resource "aws_api_gateway_method" "tflogintoken" {
  rest_api_id   = aws_api_gateway_rest_api.tfapi.id
  resource_id   = aws_api_gateway_resource.tfdmsv1vendorlogin.id
  http_method   = "POST"
  authorization = "NONE"
  api_key_required = true
  request_parameters = {
    "method.request.header.x-api-key" = true
  }
}


resource "aws_api_gateway_method_settings" "all" {
  rest_api_id = aws_api_gateway_rest_api.tfapi.id
  stage_name  = aws_api_gateway_stage.dev.stage_name
  method_path = "*/*"

  settings {
    metrics_enabled = true
    logging_level   = "INFO" // OFF ERROR  INFO
  }
}

resource "aws_api_gateway_method_response" "response_200" {
  rest_api_id = aws_api_gateway_rest_api.tfapi.id
  resource_id = aws_api_gateway_resource.tfdmsv1vendorlogin.id
  http_method = aws_api_gateway_method.tflogintoken.http_method
  status_code = "200"
}

# resource "aws_api_gateway_method_settings" "path_specific" {
#   rest_api_id = aws_api_gateway_rest_api.tfapi.id
#   stage_name  = aws_api_gateway_stage.dev.stage_name
#   method_path = "path1/GET"

#   settings {
#     metrics_enabled = true
#     logging_level   = "INFO"
#   }
# }