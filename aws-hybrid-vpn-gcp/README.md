# Terraform Local-Backend Module for GCP (Google Cloud Platform) and AWS (Amazon Web Services) Hybrid VPN (Virtual Private Network) Build

Since Terratest doesn't yet support specifying a backend config file via command-line arguments,
this internal/external module structure is required to enable E2E testing using Terratest.

See our [How to structure a Terraform module Notion page](https://www.notion.so/honestbank/How-to-structure-a-Terraform-module-31374a1594f84ef7b185ef4e06b36619)
for more details on Terraform module structuring.

This folder can be init'ed and applied using Terraform to test functionality.

To run E2E tests, navigate to the [test folder](../test) and run `go test -v -timeout 30m`.

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.1 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 3.52 |
| <a name="requirement_google"></a> [google](#requirement\_google) | ~> 4.10 |
| <a name="requirement_random"></a> [random](#requirement\_random) | ~> 3.1 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | 3.74.2 |
| <a name="provider_google"></a> [google](#provider\_google) | 4.11.0 |
| <a name="provider_random"></a> [random](#provider\_random) | 3.1.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_customer_gateway.hybrid_customer_gateway](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/customer_gateway) | resource |
| [aws_vpn_connection.hybrid_vpn_connection](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/vpn_connection) | resource |
| [aws_vpn_gateway.hybrid_vpn_gateway](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/vpn_gateway) | resource |
| [aws_vpn_gateway_route_propagation.gcp_routes](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/vpn_gateway_route_propagation) | resource |
| [google_compute_external_vpn_gateway.hybrid_external_vpn_gateway](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_external_vpn_gateway) | resource |
| [google_compute_ha_vpn_gateway.hybrid_ha_vpn_gateway](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_ha_vpn_gateway) | resource |
| [google_compute_router.hybrid_vpn_router](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_router) | resource |
| [google_compute_router_interface.tunnel_interface](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_router_interface) | resource |
| [google_compute_router_peer.hybrid_vpn_router_peers](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_router_peer) | resource |
| [google_compute_vpn_tunnel.hybrid_vpn_tunnel](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_vpn_tunnel) | resource |
| [random_password.psk](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password) | resource |
| [aws_route_tables.vpc_route_tables](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/route_tables) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_aws_subnet_ids"></a> [aws\_subnet\_ids](#input\_aws\_subnet\_ids) | AWS subnets to connect the VPN to | `list(string)` | n/a | yes |
| <a name="input_aws_vpc_id"></a> [aws\_vpc\_id](#input\_aws\_vpc\_id) | AWS VPC to connect the VPN to | `string` | n/a | yes |
| <a name="input_gcp_network_id"></a> [gcp\_network\_id](#input\_gcp\_network\_id) | ID of the GCP VPC network to connect the VPN to (must be the same network as gcp\_network\_name) | `string` | n/a | yes |
| <a name="input_gcp_network_name"></a> [gcp\_network\_name](#input\_gcp\_network\_name) | Name of the GCP VPC network to connect the VPN to | `string` | n/a | yes |
| <a name="input_gcp_subnetwork_name"></a> [gcp\_subnet\_name](#input\_gcp\_subnet\_name) | GCP subnet to connect the VPN to | `string` | n/a | yes |
| <a name="input_name"></a> [name](#input\_name) | Name of the Hybrid VPN deployment, used in various resources to uniquely identify them | `string` | `"test"` | no |

## Outputs

No outputs.
<!-- END_TF_DOCS -->
