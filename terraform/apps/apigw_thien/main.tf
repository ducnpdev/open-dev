terraform {
  required_version = ">= 1.4.6"
}

provider "aws" {
  region = var.aws_region
  profile = var.profile
  token = ""
  endpoints {
    sts = ""
  }
}
