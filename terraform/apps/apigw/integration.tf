resource "aws_api_gateway_integration" "tfintegration" {
   rest_api_id             = aws_api_gateway_rest_api.tfapi.id
   resource_id             = aws_api_gateway_resource.tfdmsv1vendorlogin.id
   http_method             = aws_api_gateway_method.tflogintoken.http_method
   integration_http_method = "POST"
   type                    = "HTTP_PROXY"
   timeout_milliseconds = 29000
   uri                     = "https://ipxxx-ranges.amazonaws.com/ip-ranges.json"
   # request_parameters = {
   #  "integration.request.header.X-Authorization" = "'static'"
   #  "integration.request.header.x-api-key" = "'static'"
#   }
#     depends_on = [
#    aws_api_gateway_method.tflogintoken
#   ]
}

## lambda integration
# resource "aws_api_gateway_integration" "tfintegration" {
#   rest_api_id             = aws_api_gateway_rest_api.tfapi.id
#   resource_id             = aws_api_gateway_resource.tfdmsv1vendor.id
#   http_method             = aws_api_gateway_method.tflogintoken.http_method
#   integration_http_method = "POST"
#   type                    = "AWS_PROXY"
#   uri                     = aws_lambda_function.lambda.invoke_arn
# }

## vpc-link integration
resource "aws_api_gateway_integration" "tfintegrationvpclink" {
#   rest_api_id = aws_api_gateway_rest_api.test.id
#   resource_id = aws_api_gateway_resource.test.id
#   http_method = aws_api_gateway_method.test.http_method
   rest_api_id             = aws_api_gateway_rest_api.tfapi.id
   resource_id             = aws_api_gateway_resource.tfdmsv1vendorlogin.id
   http_method             = aws_api_gateway_method.tflogintoken.http_method
   integration_http_method = "POST"
#   request_templates = {
#     "application/json" = ""
#     "application/xml"  = "#set($inputRoot = $input.path('$'))\n{ }"
#   }

#   request_parameters = {
#     "integration.request.header.X-Authorization" = "'static'"
#     "integration.request.header.X-Foo"           = "'Bar'"
#   }

  type                    = "HTTP_PROXY"
  uri                     = "https://www.google.de"
#   integration_http_method = "GET"
#   passthrough_behavior    = "WHEN_NO_MATCH"
#   content_handling        = "CONVERT_TO_TEXT"

  connection_type = "VPC_LINK"
  connection_id   = "https://ipxxx-ranges.amazonaws.com/ip-ranges.json"
}