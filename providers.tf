provider "google" {
  credentials = var.google_credentials
  project     = var.google_project
  region      = var.google_region
}

provider "aws" {}
