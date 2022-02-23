# Terraform AWS Hybrid VPN GCP

HA IPsec managed cloud VPN between AWS and GCP with automatic route propagation into target VPCs.

<!-- BEGIN_TF_DOCS -->
## Requirements

No requirements.

## Providers

No providers.

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_aws-hybrid-vpn-gcp"></a> [aws-hybrid-vpn-gcp](#module\_aws-hybrid-vpn-gcp) | ./aws-hybrid-vpn-gcp | n/a |

## Resources

No resources.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_aws_subnet_ids"></a> [aws\_subnet\_ids](#input\_aws\_subnet\_ids) | AWS subnets to connect the VPN to | `list(string)` | n/a | yes |
| <a name="input_aws_vpc_id"></a> [aws\_vpc\_id](#input\_aws\_vpc\_id) | AWS VPC to connect the VPN to | `string` | n/a | yes |
| <a name="input_gcp_network_id"></a> [gcp\_network\_id](#input\_gcp\_network\_id) | ID of the GCP VPC network to connect the VPN to (must be the same network as gcp\_network\_name) | `string` | n/a | yes |
| <a name="input_gcp_network_name"></a> [gcp\_network\_name](#input\_gcp\_network\_name) | Name of the GCP VPC network to connect the VPN to | `string` | n/a | yes |
| <a name="input_gcp_subnetwork_name"></a> [gcp\_subnetwork\_name](#input\_gcp\_subnetwork\_name) | GCP subnet to connect the VPN to | `string` | n/a | yes |
| <a name="input_google_credentials"></a> [google\_credentials](#input\_google\_credentials) | Contents of a JSON keyfile of an account with write access to the project | `string` | n/a | yes |
| <a name="input_google_project"></a> [google\_project](#input\_google\_project) | The GCP project to use for this run | `string` | n/a | yes |
| <a name="input_google_region"></a> [google\_region](#input\_google\_region) | GCP region used to create all resources in this run | `string` | n/a | yes |
| <a name="input_name"></a> [name](#input\_name) | Name of the Hybrid VPN connection, used in various resources to uniquely identify them | `string` | `"test"` | no |

## Outputs

No outputs.
<!-- END_TF_DOCS -->
