variable "aws_region" {
  default     = "ap-southeast-1"
  description = "AWS Region to deploy example API Gateway REST API"
  type        = string
}

variable "profile" {
  default     = "awsuat"
  # description = "AWS Region to deploy example API Gateway REST API"
  type        = string
}
# apigw
variable "apigw_id" {
  default     = "xxx"
  description = "id of apigw"
  type        = string
}

variable "api_key_value" {
  default     = "ekyc-service-xxxx"
  description = "id of apigw"
  type        = string
}
variable "api_key_name" {
  default     = "ekyc-service"
  description = "id of apigw"
  type        = string
}

variable "apigw_root_resource" {
  default     = "xxx"
  description = "root resource"
  type        = string
}

# usage plan
variable "apigw_stage" {
  default     = "uat"
  description = "stage of apigw"
  type        = string
}

variable "throttle_method_path_image_get" {
  default     = "/ekycservice/path/POST"
  description = "path of usage"
  type        = string
}


# 
variable "vpc_link_id" {
  default     = "yyy"
  description = "id of vpc link"
  type        = string
}

variable "vpc_link_url" {
  default     = "https://google.com"
  description = "url of vpc link"
  type        = string
}