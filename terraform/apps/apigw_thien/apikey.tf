resource "aws_api_gateway_api_key" "ekyc_service" {
  name = var.api_key_name
  description = "ekyc service"
  enabled = true
  value = var.api_key_value
}

resource "aws_api_gateway_usage_plan_key" "ekyc_service_key" {
  key_id        = aws_api_gateway_api_key.ekyc_service.id
  key_type      = "API_KEY"
  usage_plan_id = aws_api_gateway_usage_plan.ekyc_service.id
}