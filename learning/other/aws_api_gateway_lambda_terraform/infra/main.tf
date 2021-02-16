provider "aws" {
  region     = var.aws_region
  access_key = var.aws_access_key
  secret_key = var.aws_secret_key
}

variable "aws_access_key" {
  type = string
  description = "AWS access key"
}
variable "aws_secret_key" {
  type = string
  description = "AWS secret key"
}
variable "aws_region" {
  type = string
  description = "AWS region"
}

variable "app_name" {
  description = "Application name"
  default = "sample-lambda-api"
}

variable "app_env" {
  description = "Apllication environment tag"
  default = "dev"
}

locals {
  app_id = "${lower(var.app_name)}-${lower(var.app_env)}-${random_id.unique_suffix.hex}"
}

resource "random_id" "unique_suffix" {
  byte_length = 2
}

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = "build/bin/app"
  output_path = "build/bin/app.zip"
}

output "api_url" {
  value = aws_api_gateway_deployment.api_deployment.invoke_url
}
