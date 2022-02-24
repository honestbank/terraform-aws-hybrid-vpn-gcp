variable "name" {
  description = "Name to use on resources"
  type        = string
  default     = "test"
}

variable "aws_vpc_id" {
  description = "AWS VPC to create EC2 instance in"
  type        = string
}

variable "aws_subnet_id" {
  description = "AWS subnet to create EC2 instance in"
  type        = string
}

variable "gcp_network_name" {
  description = "GCP network to create compute instance in"
  type        = string
}

variable "gcp_subnetwork_name" {
  description = "GCP subnet to create the compute instance in"
  type        = string
}

variable "gcp_zone_name" {
  description = "GCP availability zone to create the compute instance in"
  type        = string
}
