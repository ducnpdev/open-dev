
## vpc-link integration
resource "aws_api_gateway_integration" "ekyc_service_integration_nlb" {
   rest_api_id             = var.apigw_id
   resource_id             = aws_api_gateway_resource.ekyc_dev_v1_get_image.id
   http_method             = aws_api_gateway_method.get_image.http_method
   integration_http_method = "POST"
   type                    = "HTTP"
   uri                     = var.vpc_link_url
   connection_type = "VPC_LINK"
   connection_id   = var.vpc_link_id
}