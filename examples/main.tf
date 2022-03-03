terraform {
  required_providers {
     aws-download = {
      source = "nicolai86/aws-download"
      version = "0.0.5"
    }
  }
}

variable "role_arn" {
  type = string 
  description = "aws role arn to assume"
}

provider "aws-download" {
  role_arn     = var.role_arn
  session_name = "Terraform"
  region = "us-west-2"
}

variable "s3_uri" {
  type = string 
  description = "s3 uri of file to download"
}

locals {
  s3_parts  = split("/", var.s3_uri)
  s3_bucket = local.s3_parts[2]
  s3_key    = "/${join("/", slice(local.s3_parts, 3, length(local.s3_parts)))}"
}

data "aws-download_s3_object" "example" {
  bucket = local.s3_bucket
  key    = local.s3_key
  filename = "/tmp/test"
}
