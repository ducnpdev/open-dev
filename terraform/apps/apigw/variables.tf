variable "aws_region" {
  default     = "ap-southeast-1"
  description = "AWS Region to deploy example API Gateway REST API"
  type        = string
}
variable "rest_api_name" {
  default     = "tf-study-api-gateway-rest-api"
  description = "Name of the API Gateway REST API (can be used to trigger redeployments)"
  type        = string
}

variable "rest_api_path" {
  default     = "/path1"
  description = "Path to create in the API Gateway REST API (can be used to trigger redeployments)"
  type        = string
}

variable "rest_api_path_dms" {
  default     = "/dms/v1/vendor/login"
  description = "Path to create in the API Gateway REST API (can be used to trigger redeployments)"
  type        = string
}

variable "rest_api_path_dms_1" {
  default     = "/dms/v1/file/create"
  description = "Path to create in the API Gateway REST API (can be used to trigger redeployments)"
  type        = string
}
