variable "aws_assume_role_arn" {
  description = "AWS role to assume for this run"
}

variable "aws_region" {
  description = "GCP Region to use for this run"
}

provider "aws" {
  dynamic "assume_role" {
    for_each = var.aws_assume_role_arn != "" ? [] : [1]
    content {
      role_arn = var.aws_assume_role_arn
      session_name = "terratest"
    }
  }

  region = var.aws_region

  # Jakarta not in provider whitelist at time of writing
  skip_region_validation = true
}