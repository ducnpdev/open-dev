variable "aws_region" {
  default     = "ap-southeast-1"
  description = "AWS Region to deploy sqs"
  type        = string
}
variable "sqs_hdss_name" {
  default     = "queue_name.fifo"
  description = "sqs name"
  type        = string
}
