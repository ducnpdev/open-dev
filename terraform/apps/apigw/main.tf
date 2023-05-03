terraform {
  required_version = ">= 1.4.6"
}

provider "aws" {
  access_key = ""
  secret_key = ""
  region = var.aws_region
}

resource "aws_s3_bucket" "first_bucket" {
    bucket = "example-first-bucket"
}

