output "ec2_instance_public_ip" {
  value = aws_instance.test.public_ip
}

output "ec2_instance_private_ip" {
  value = aws_instance.test.private_ip
}

output "gce_instance_public_ip" {
  value = google_compute_instance.test.network_interface[0].access_config[0].nat_ip
}

output "gce_instance_private_ip" {
  value = google_compute_instance.test.network_interface[0].network_ip
}

output "ssh_private_key" {
  value = tls_private_key.ssh_keypair.private_key_pem
  sensitive = true
}