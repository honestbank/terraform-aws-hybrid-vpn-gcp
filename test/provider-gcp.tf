variable "google_credentials" {
  description = "GCP Service Account JSON keyfile contents."
  sensitive   = true
}

provider "google" {
  credentials = var.google_credentials
  project     = var.google_project
  region      = var.google_region
}
