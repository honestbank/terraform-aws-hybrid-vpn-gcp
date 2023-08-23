# TODO: remove
resource "aws_vpn_gateway" "hybrid_vpn_gateway" {
  amazon_side_asn = local.aws_bgp_asn

  vpc_id = var.aws_vpc_id

  tags = {
    Name = "hybrid-vpn-gateway-${var.resource_suffix}"
  }
}

resource "aws_ec2_transit_gateway" "hybrid_vpn_transit_gateway" {
  amazon_side_asn = local.aws_bgp_asn
}

resource "aws_customer_gateway" "hybrid_customer_gateway" {
  for_each = local.gcp_public_ips_set

  bgp_asn    = local.gcp_bgp_asn
  ip_address = google_compute_ha_vpn_gateway.hybrid_ha_vpn_gateway.vpn_interfaces[index(local.gcp_public_ips, each.key)].ip_address
  type       = "ipsec.1"

  tags = {
    Name = "hybrid-customer-gateway-${var.resource_suffix}-${each.key}"
  }
}

resource "aws_vpn_connection" "hybrid_vpn_connection" {
  for_each = local.gcp_public_ips_set

  transit_gateway_id = aws_ec2_transit_gateway.hybrid_vpn_transit_gateway.id

  customer_gateway_id                  = aws_customer_gateway.hybrid_customer_gateway[each.value].id
  type                                 = aws_customer_gateway.hybrid_customer_gateway.type
  tunnel1_inside_cidr                  = local.vpn_tunnel_inside_cidrs[tonumber(each.value)][0]
  tunnel2_inside_cidr                  = local.vpn_tunnel_inside_cidrs[tonumber(each.value)][1]
  tunnel1_preshared_key                = random_password.psk[local.aws_vpn_tunnel_psk_index_map[tonumber(each.value)][0]].result
  tunnel2_preshared_key                = random_password.psk[local.aws_vpn_tunnel_psk_index_map[tonumber(each.value)][1]].result
  tunnel1_phase1_encryption_algorithms = ["AES256"]
  tunnel1_phase2_encryption_algorithms = ["AES256"]
  tunnel2_phase1_encryption_algorithms = ["AES256"]
  tunnel2_phase2_encryption_algorithms = ["AES256"]
  tunnel1_phase1_integrity_algorithms  = ["SHA2-256"]
  tunnel1_phase2_integrity_algorithms  = ["SHA2-256"]
  tunnel2_phase1_integrity_algorithms  = ["SHA2-256"]
  tunnel2_phase2_integrity_algorithms  = ["SHA2-256"]
  tunnel1_phase1_dh_group_numbers      = [18]
  tunnel1_phase2_dh_group_numbers      = [18]
  tunnel2_phase1_dh_group_numbers      = [18]
  tunnel2_phase2_dh_group_numbers      = [18]

  tags = {
    Name = "hybrid-vpn-connection-${var.resource_suffix}-${each.key}"
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "vpc_transit_gateway_attachment" {
  subnet_ids         = var.aws_subnet_ids
  transit_gateway_id = aws_ec2_transit_gateway.hybrid_vpn_transit_gateway.id
  vpc_id             = var.aws_vpc_id
}

# TODO: update
#resource "aws_ec2_transit_gateway_route" "example" {
#  destination_cidr_block         = "0.0.0.0/0"
#  transit_gateway_attachment_id  = aws_ec2_transit_gateway_vpc_attachment.example.id
#  transit_gateway_route_table_id = aws_ec2_transit_gateway.example.association_default_route_table_id
#}

resource "aws_ec2_transit_gateway_route_table" "vpn_transit_gateway_route_table" {
  transit_gateway_id = aws_ec2_transit_gateway.hybrid_vpn_transit_gateway.id
}

# TODO: `aws_ec2_transit_gateway_route_table_association` and `aws_ec2_transit_gateway_route_table_propagation`
# See example: https://github.com/mikemowgli/terraform-aws-transit-gateway/blob/master/deploy.tf (the code is 4 years old though)
# Very helpful for understanding route tables: https://docs.aws.amazon.com/vpc/latest/tgw/tgw-transit-gateways.html
data "aws_route_tables" "vpc_route_tables" {
  vpc_id = var.aws_vpc_id
}

resource "aws_vpn_gateway_route_propagation" "gcp_routes" {
  for_each = toset(data.aws_route_tables.vpc_route_tables.ids)

  vpn_gateway_id = aws_vpn_gateway.hybrid_vpn_gateway.id
  route_table_id = each.value
}
