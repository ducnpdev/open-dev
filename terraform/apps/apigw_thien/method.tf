resource "aws_api_gateway_method" "get_image" {
  rest_api_id   = var.apigw_id
  resource_id   = aws_api_gateway_resource.ekyc_dev_v1_get_image.id
  http_method   = "POST"

  authorization = "NONE"
  api_key_required = true
  request_parameters = {
    "method.request.header.x-api-key" = true
  }
}