terraform {
  required_version = ">= 1.4.6"
}

provider "aws" {
  access_key = ""
  secret_key = ""
  region = var.aws_region
}
