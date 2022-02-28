resource "random_password" "psk" {
  count = length(local.gcp_public_ips) * 2

  length  = 63
  special = true
  # These are the only supported special characters
  override_special = "_."

  # Cannot start with a 0, so exclude numbers
  number = false
}
