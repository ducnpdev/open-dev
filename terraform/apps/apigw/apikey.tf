resource "aws_api_gateway_api_key" "tfAPIKey" {
  name = "tf-api-key"
  description = "terrafrom description"
  enabled = true
  value = "2af0b7bc-4375-4df2-a44d-2eff998a1193"
}

resource "aws_api_gateway_usage_plan_key" "tfUsagePlanKey" {
  key_id        = aws_api_gateway_api_key.tfAPIKey.id
  key_type      = "API_KEY"
  usage_plan_id = aws_api_gateway_usage_plan.tfUsagePlan.id
}