resource "aws_api_gateway_usage_plan" "ekyc_service" {
  name         = "ekyc_service"
  description  = ""
  product_code = ""

  api_stages {
    api_id = var.apigw_id
    stage  = var.apigw_stage
    throttle {
      path = var.throttle_method_path_image_get
      burst_limit = 3
      rate_limit  = 6
    }
  }

  quota_settings {
    limit  = 10000
    # offset = 2
    period = "DAY" // WEEK
  }
  
  throttle_settings {
    burst_limit = 5
    rate_limit  = 10
  }



}