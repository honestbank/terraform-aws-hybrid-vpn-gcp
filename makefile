lint:
	terraform fmt --recursive

validate:
	terraform init
	terraform validate
	terraform fmt --recursive

docs:
	rm -rf aws-hybrid-vpn-gcp/*/.terraform aws-hybrid-vpn-gcp/*/.terraform.lock.hcl
	rm -rf modules/*/.terraform modules/*/.terraform.lock.hcl

	terraform-docs -c .terraform-docs.yml .
