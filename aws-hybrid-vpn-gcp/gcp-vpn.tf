resource "google_compute_router" "hybrid_vpn_router" {
  name    = "hybrid-vpn-router-${var.resource_suffix}"
  network = var.gcp_network_name

  bgp {
    asn               = local.gcp_bgp_asn
    advertised_groups = ["ALL_SUBNETS"]
    advertise_mode    = "CUSTOM"
  }
}

# VPN Gateway
resource "google_compute_ha_vpn_gateway" "hybrid_ha_vpn_gateway" {
  name    = "hybrid-vpn-gateway-${var.resource_suffix}"
  network = var.gcp_network_id
}

locals {
  hybrid_vpn_connection_helper = [
    aws_vpn_connection.hybrid_vpn_connection[element(local.gcp_public_ips, 0)].tunnel1_address,
    aws_vpn_connection.hybrid_vpn_connection[element(local.gcp_public_ips, 0)].tunnel2_address,
    aws_vpn_connection.hybrid_vpn_connection[element(local.gcp_public_ips, 1)].tunnel1_address,
    aws_vpn_connection.hybrid_vpn_connection[element(local.gcp_public_ips, 1)].tunnel2_address
  ]
}

# External VPN Gateway
resource "google_compute_external_vpn_gateway" "hybrid_external_vpn_gateway" {
  name            = "hybrid-external-vpn-gateway-${var.resource_suffix}"
  redundancy_type = "FOUR_IPS_REDUNDANCY"

  dynamic "interface" {
    for_each = local.hybrid_vpn_connection_helper
    iterator = vpn_connection
    content {
      id         = vpn_connection.key
      ip_address = vpn_connection.value
    }
  }
}

resource "google_compute_vpn_tunnel" "hybrid_vpn_tunnel" {
  count = 4

  name = "hybrid-vpn-tunnel-${var.resource_suffix}-${count.index}"

  peer_external_gateway           = google_compute_external_vpn_gateway.hybrid_external_vpn_gateway.id
  peer_external_gateway_interface = count.index

  ike_version   = 2
  shared_secret = random_password.psk[count.index].result

  router                = google_compute_router.hybrid_vpn_router.id
  vpn_gateway           = google_compute_ha_vpn_gateway.hybrid_ha_vpn_gateway.id
  vpn_gateway_interface = count.index < 2 ? 0 : 1
}

resource "google_compute_router_interface" "tunnel_interface" {
  count = 4

  name       = "tunnel-${var.resource_suffix}-${count.index}"
  router     = google_compute_router.hybrid_vpn_router.name
  ip_range   = "169.254.1${count.index}.2/30"
  vpn_tunnel = google_compute_vpn_tunnel.hybrid_vpn_tunnel[count.index].name
}


resource "google_compute_router_peer" "hybrid_vpn_router_peers" {
  count = 4

  name            = "aws-connection-${count.index < 2 ? 0 : 1}-tunnel-${local.aws_vpn_connection_tunnel_index_map[count.index]}"
  peer_asn        = local.aws_bgp_asn
  peer_ip_address = "169.254.1${count.index}.1"
  interface       = "tunnel-${var.resource_suffix}-${count.index}"
  router          = google_compute_router.hybrid_vpn_router.name
}
