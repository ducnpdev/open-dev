resource "aws_api_gateway_rest_api" "tfapi" {
  body = jsonencode({
    openapi = "3.0.1"
    info = {
      title   = var.rest_api_name
      version = "1.0"
    }
    paths = {
      (var.rest_api_path) = {
        get = {
          x-amazon-apigateway-integration = {
            httpMethod           = "GET"
            payloadFormatVersion = "1.0"
            type                 = "HTTP_PROXY"
            uri                  = "https://ip-ranges.amazonaws.com/ip-ranges.json"
          }
        }
      },
      # (var.rest_api_path_dms) = {
      #   post = {
      #     x-amazon-apigateway-integration = {
      #       httpMethod           = "POST"
      #       payloadFormatVersion = "1.0"
      #       type                 = "HTTP_PROXY"
      #       uri                  = "https://ip-ranges.amazonaws.com/ip-ranges.json"
      #     }
      #   }
      # },
      #  (var.rest_api_path_dms_1) = {
      #   post = {
      #     x-amazon-apigateway-integration = {
      #       httpMethod           = "POST"
      #       payloadFormatVersion = "1.0"
      #       type                 = "HTTP_PROXY"
      #       uri                  = "https://ip-ranges.amazonaws.com/ip-ranges.json"
      #     }
      #   }
      # },
    }
  })

  name = var.rest_api_name

  endpoint_configuration {
    types = ["EDGE"] // REGIONAL
  }
}


# resource "aws_api_gateway_method" "MyDemoMethod" {
#   rest_api_id   = aws_api_gateway_rest_api.tfapi.id
#   resource_id   = aws_api_gateway_resource.MyDemoResource.id
#   http_method   = "GET"
#   authorization = "NONE"
# }

# resource "aws_api_gateway_model" "MyDemoModel" {
#   rest_api_id  = aws_api_gateway_rest_api.tfapi.id
#   name         = "tfusermodel"
#   description  = "a JSON schema"
#   content_type = "application/json"

#   schema = jsonencode({
#     type = "object"
#   })
# }