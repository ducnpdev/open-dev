resource "aws_sqs_queue" "terraform_queue" {
  name                  =  var.sqs_hdss_name
  fifo_queue            = true
  content_based_deduplication = true
  visibility_timeout_seconds = 120
  message_retention_seconds = 432000
}