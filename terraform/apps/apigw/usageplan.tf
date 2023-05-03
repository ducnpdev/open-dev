resource "aws_api_gateway_usage_plan" "tfUsagePlan" {
  name         = "tf-my-usage-plan"
  description  = "my description"
  product_code = "MYCODE"

  api_stages {
    api_id = aws_api_gateway_rest_api.tfapi.id
    stage  = aws_api_gateway_stage.dev.stage_name
  }

#   api_stages {
#     api_id = aws_api_gateway_rest_api.example.id
#     stage  = aws_api_gateway_stage.production.stage_name
#   }

  quota_settings {
    limit  = 20
    # offset = 2
    period = "DAY" // WEEK
  }

  throttle_settings {
    burst_limit = 5
    rate_limit  = 10
  }
}