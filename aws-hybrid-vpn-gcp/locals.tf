locals {
  gcp_bgp_asn = 16550
  aws_bgp_asn = 65479
}

locals {
  # Helper local to get correct internal VPN CIDR when creating AWS VPN connections in pairs
  vpn_tunnel_inside_cidrs = [
    {
      0 = "169.254.10.0/30",
      1 = "169.254.11.0/30",
    }, {
      0 = "169.254.12.0/30",
      1 = "169.254.13.0/30",
    }
  ]

  # Helper to get correct tunnel values when creating AWS VPN connections in pairs
  aws_vpn_tunnel_psk_index_map = [
    {
      0 = 0,
      1 = 1,
    },
    {
      0 = 2,
      1 = 3,
    }
  ]

  # Helper local to get correct AWS connections when looping over GCP tunnel peers
  aws_vpn_connection_tunnel_index_map = {
    0 = 0
    1 = 1
    2 = 0
    3 = 1
  }
}

locals {
  gcp_public_ips     = ["0", "1"] # Must be ["0", "1"], used to calculate VPN tunnel inside CIDRs
  gcp_public_ips_set = toset(local.gcp_public_ips)
}
