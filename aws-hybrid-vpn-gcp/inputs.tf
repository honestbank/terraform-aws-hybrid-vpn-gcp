variable "resource_suffix" {
  description = "String used in various resources to suffix resource names related to this VPN setup"
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
