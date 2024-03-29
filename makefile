lint:
	terraform fmt --recursive

validate:
	terraform init
	terraform validate
	terraform fmt --recursive

docs:
	terraform-docs .
	terraform-docs -c .terraform-docs.yml ./aws-hybrid-vpn-gcp
