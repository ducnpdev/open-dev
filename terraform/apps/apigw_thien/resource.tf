resource "aws_api_gateway_resource" "ekyc_dev" {
  rest_api_id = var.apigw_id
  parent_id   = var.apigw_root_resource
  path_part   = "ekycservice"
}

resource "aws_api_gateway_resource" "ekyc_dev_v1" {
  rest_api_id = var.apigw_id
  parent_id   = aws_api_gateway_resource.ekyc_dev.id
  path_part   = "v1"
}

resource "aws_api_gateway_resource" "ekyc_dev_v1_get_image" {
  rest_api_id = var.apigw_id
  parent_id   = aws_api_gateway_resource.ekyc_dev_v1.id
  path_part   = "get-images"
}
