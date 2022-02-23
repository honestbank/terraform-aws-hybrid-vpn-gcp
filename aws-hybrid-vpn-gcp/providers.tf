terraform {
  required_version = "~> 1.1"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.52"
    }

    google = {
      version = "~> 4.10"
    }

    random = {
      version = "~> 3.1"
    }
  }
}
