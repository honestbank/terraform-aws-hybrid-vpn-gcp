resource "tls_private_key" "ssh_keypair" {
  algorithm = "RSA"
}

// AWS

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "namea"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Owner ID of Canonical (Ubuntu)
}

resource "aws_key_pair" "terratest_key" {
  key_name   = "terratest-${var.name}"
  public_key = tls_private_key.ssh_keypair.public_key_openssh
}

resource "aws_security_group" "security_group" {
  name   = "terratest-${var.name}"
  vpc_id = var.aws_vpc_id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "test" {
  ami                         = data.aws_ami.ubuntu.id
  instance_type               = "t3.micro"
  subnet_id                   = var.aws_subnet_id
  associate_public_ip_address = true
  key_name                    = aws_key_pair.terratest_key.key_name
  security_groups             = [aws_security_group.security_group.id]

  tags = {
    Name = "terratest-${var.name}"
  }
}

// GCP

resource "google_compute_instance" "test" {
  name         = "terratest-${var.name}"
  machine_type = "e2-micro"
  zone         = var.gcp_zone_name

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  network_interface {
    network    = var.gcp_network_name
    subnetwork = var.gcp_subnetwork_name

    // Associates an ephemeral public IP
    access_config {}
  }

  metadata = {
    ssh-keys = "terratest:${tls_private_key.ssh_keypair.public_key_openssh}"
  }
}

resource "google_compute_firewall" "external_access" {
  name    = "terratest-${var.name}-external-access"
  network = var.gcp_network_name

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["0.0.0.0/0"]
}

resource "google_compute_firewall" "internal_access" {
  name    = "terratest-${var.name}-internal-access"
  network = var.gcp_network_name

  allow {
    protocol = "all"
  }

  # Covers all of GCP and AWS subnets
  source_ranges = ["10.133.7.0/24"]
}
