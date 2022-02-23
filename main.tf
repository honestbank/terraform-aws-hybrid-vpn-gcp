module "aws-hybrid-vpn-gcp" {
  source = "./aws-hybrid-vpn-gcp"

  name = var.name

  # AWS Target VPC
  aws_vpc_id                 = var.aws_vpc_id
  aws_subnet_ids             = var.aws_subnet_ids

  # GCP Target VPC
  gcp_network_name = var.gcp_network_name
  gcp_network_id   = var.gcp_network_id
  gcp_subnetwork_name  = var.gcp_subnetwork_name
}
