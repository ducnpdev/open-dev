variable "aws_region" {
  default     = "ap-southeast-1"
  description = "AWS Region to deploy parameter"
  type        = string
}

variable "parameter_name_key" {
  default     = "sm-parameter-apse1-lab-name"
  description = "paramter store of name"
  type        = string
}
