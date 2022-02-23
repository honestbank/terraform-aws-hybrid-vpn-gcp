variable "name" {
  description = "Name of the Hybrid VPN connection, used in various resources to uniquely identify them"
  type        = string
  default     = "test"
}

variable "google_project" {
  description = "The GCP project to use for this run"
  type        = string
}

variable "google_credentials" {
  description = "Contents of a JSON keyfile of an account with write access to the project"
  type        = string
}

variable "google_region" {
  description = "GCP region used to create all resources in this run"
  type        = string
}

variable "aws_vpc_id" {
  description = "AWS VPC to connect the VPN to"
  type        = string
}

variable "aws_subnet_ids" {
  description = "AWS subnets to connect the VPN to"
  type        = list(string)
}

variable "gcp_network_name" {
  description = "Name of the GCP VPC network to connect the VPN to"
  type        = string
}

variable "gcp_network_id" {
  description = "ID of the GCP VPC network to connect the VPN to (must be the same network as gcp_network_name)"
  type        = string
}

variable "gcp_subnetwork_name" {
  description = "GCP subnet to connect the VPN to"
  type        = string
}
