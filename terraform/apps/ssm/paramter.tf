resource "aws_ssm_parameter" "parameter_name" {
  name  = var.parameter_name_key
  type  = "String"
  value = "{\"secretKey\":\"123\",\"url\":\"https://google.com\",\"secretKeyExternal\":\"\",\"headers\":{\"key\":\"9Buvd48nJsMHe\",\"x-api-key\":\"9Buvd48p5p2HJsMHe\",\"Content-Type\":\"application/json\"}}"
}