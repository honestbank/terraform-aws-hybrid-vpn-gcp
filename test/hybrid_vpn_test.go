package test

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

func TestTerraformGcpHybridVpnAws(t *testing.T) {
	t.Parallel()
	runId, _ := runID()
	awsRegion := "ap-southeast-3"
	awsAssumeRoleArn := "arn:aws:iam::106256755710:role/OrganizationAccountAccessRole"
	gcpRegion := "asia-southeast2"
	gcpProject := "test-terraform-project-01"

	googleCredentials := getGoogleCredentials()

	//
	//
	// AWS VPC
	//
	//

	awsVpcBootstrapTerraformOptions := &terraform.Options{}
	awsVpcDir := test_structure.CopyTerraformFolderToTemp(t, ".", "modules/terraform-aws-vpc/aws-vpc")

	test_structure.RunTestStage(t, "create_aws_vpc_terraform_options", func() {
		awsVpcInputs := map[string]interface{}{
			"aws_region":          awsRegion,
			"aws_assume_role_arn": awsAssumeRoleArn,
			"name":                "vpc-" + runId,
			"azs":                 []string{awsRegion + "a", awsRegion + "b", awsRegion + "c"},
			"cidr":                "10.133.7.0/25",
			"public_subnets":      []string{"10.133.7.64/26"},
			"private_subnets":     []string{"10.133.7.0/26"},
			"enable_flow_log":     false,
			"flow_log_cloudwatch_log_group_retention_in_days": 0,
		}

		awsVpcBootstrapTerraformOptions = terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: awsVpcDir,
			Vars:         awsVpcInputs,
		})

		copySupportingFiles(t, []string{"provider-aws.tf"}, awsVpcDir)
	})

	defer test_structure.RunTestStage(t, "cleanup_aws_vpc_support_files", func() {
		_ = cleanupSupportingFiles([]string{"provider-aws.tf"}, awsVpcDir)
	})

	defer test_structure.RunTestStage(t, "cleanup_aws_vpc", func() {
		terraform.Destroy(t, awsVpcBootstrapTerraformOptions)
	})

	test_structure.RunTestStage(t, "create_aws_vpc", func() {
		terraform.InitAndApply(t, awsVpcBootstrapTerraformOptions)
	})

	//
	//
	// GCP VPC
	//
	//

	gcpVpcBootstrapTerraformOptions := &terraform.Options{}
	gcpVpcDir := test_structure.CopyTerraformFolderToTemp(t, ".", "modules/terraform-gcp-vpc/vpc")

	test_structure.RunTestStage(t, "create_gcp_vpc_terraform_options", func() {
		gcpVpcInputs := map[string]interface{}{
			"google_project":                       gcpProject,
			"google_region":                        gcpRegion,
			"google_credentials":                   googleCredentials,
			"network_name":                         "network-" + runId,
			"vpc_primary_subnet_name":              "subnet-" + runId,
			"vpc_routing_mode":                     "REGIONAL",
			"vpc_primary_subnet_ip_range_cidr":     "10.133.7.128/27",
			"vpc_secondary_ip_range_pods_name":     "pods",
			"vpc_secondary_ip_range_pods_cidr":     "10.133.7.160/27",
			"vpc_secondary_ip_range_services_name": "services",
			"vpc_secondary_ip_range_services_cidr": "10.133.7.192/27",
		}

		gcpVpcBootstrapTerraformOptions = terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: gcpVpcDir,
			Vars:         gcpVpcInputs,
		})

		copySupportingFiles(t, []string{"provider-gcp.tf"}, gcpVpcDir)
	})

	defer test_structure.RunTestStage(t, "cleanup_gcp_vpc_support_files", func() {
		_ = cleanupSupportingFiles([]string{"provider-gcp.tf"}, gcpVpcDir)
	})

	defer test_structure.RunTestStage(t, "cleanup_gcp_vpc", func() {
		terraform.Destroy(t, gcpVpcBootstrapTerraformOptions)
	})

	test_structure.RunTestStage(t, "create_gcp_vpc", func() {
		terraform.InitAndApply(t, gcpVpcBootstrapTerraformOptions)
	})

	//
	//
	// Hybrid VPN
	//
	//

	awsHybridVpnGcpTerraformOptions := &terraform.Options{}
	awsHybridVpnGcpDir := test_structure.CopyTerraformFolderToTemp(t, "..", "aws-hybrid-vpn-gcp")

	test_structure.RunTestStage(t, "create_aws_hybrid_vpn_gcp_terraform_options", func() {
		awsHybridVpnInputs := map[string]interface{}{
			"aws_region":          awsRegion,
			"aws_assume_role_arn": awsAssumeRoleArn,
			"google_region":       gcpRegion,
			"google_project":      gcpProject,
			"google_credentials":  googleCredentials,
			"resource_suffix":     runId,
			"aws_vpc_id":          terraform.Output(t, awsVpcBootstrapTerraformOptions, "vpc_id"),
			"aws_subnet_ids":      terraform.OutputList(t, awsVpcBootstrapTerraformOptions, "private_subnets"),
			"gcp_network_name":    terraform.Output(t, gcpVpcBootstrapTerraformOptions, "network_name"),
			"gcp_network_id":      terraform.Output(t, gcpVpcBootstrapTerraformOptions, "network_id"),
			"gcp_subnetwork_name": terraform.Output(t, gcpVpcBootstrapTerraformOptions, "primary_subnet_name"),
		}

		awsHybridVpnGcpTerraformOptions = terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: awsHybridVpnGcpDir,
			Vars:         awsHybridVpnInputs,
		})

		copySupportingFiles(t, []string{"provider-aws.tf", "provider-gcp.tf", "extra-inputs.tf"}, awsHybridVpnGcpDir)
	})

	defer test_structure.RunTestStage(t, "cleanup_aws_hybrid_vpn_gcp_support_files", func() {
		_ = cleanupSupportingFiles([]string{"provider-aws.tf", "provider-gcp.tf", "extra-inputs.tf"}, awsHybridVpnGcpDir)
	})

	defer test_structure.RunTestStage(t, "cleanup_aws_hybrid_vpn_gcp", func() {
		terraform.Destroy(t, awsHybridVpnGcpTerraformOptions)
	})

	test_structure.RunTestStage(t, "create_aws_hybrid_vpn_gcp", func() {
		terraform.InitAndApply(t, awsHybridVpnGcpTerraformOptions)
	})

	//
	//
	// Connectivity Validation
	//
	//

	validationResourcesTerraformOptions := &terraform.Options{}
	validationResourcesDir := test_structure.CopyTerraformFolderToTemp(t, ".", "modules/validation-resources")

	test_structure.RunTestStage(t, "create_validation_resources_options", func() {
		validationResourcesInputs := map[string]interface{}{
			"aws_region":          awsRegion,
			"aws_assume_role_arn": awsAssumeRoleArn,
			"google_region":       gcpRegion,
			"google_project":      gcpProject,
			"google_credentials":  googleCredentials,
			"name":                runId,
			"aws_vpc_id":          terraform.Output(t, awsVpcBootstrapTerraformOptions, "vpc_id"),
			"aws_subnet_id":       terraform.OutputList(t, awsVpcBootstrapTerraformOptions, "private_subnets")[0],
			"gcp_network_name":    terraform.Output(t, gcpVpcBootstrapTerraformOptions, "network_name"),
			"gcp_subnetwork_name": terraform.Output(t, gcpVpcBootstrapTerraformOptions, "primary_subnet_name"),
			"gcp_zone_name":       gcpRegion + "-a",
		}

		validationResourcesTerraformOptions = terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: validationResourcesDir,
			Vars:         validationResourcesInputs,
		})

		copySupportingFiles(t, []string{"provider-aws.tf", "provider-gcp.tf", "extra-inputs.tf"}, validationResourcesDir)
	})

	defer test_structure.RunTestStage(t, "cleanup_validation_resources_support_files", func() {
		_ = cleanupSupportingFiles([]string{"provider-aws.tf", "provider-gcp.tf", "extra-inputs.tf"}, validationResourcesDir)
	})

	defer test_structure.RunTestStage(t, "cleanup_validation_resources", func() {
		terraform.Destroy(t, validationResourcesTerraformOptions)
	})

	test_structure.RunTestStage(t, "create_validation_resources", func() {
		terraform.InitAndApply(t, validationResourcesTerraformOptions)
	})

	test_structure.RunTestStage(t, "validate_connectivity", func() {
		sshKeyPair := ssh.KeyPair{
			PrivateKey: terraform.Output(t, validationResourcesTerraformOptions, "ssh_private_key"),
		}

		sshAgent := ssh.SshAgentWithKeyPair(t, &sshKeyPair)
		fmt.Println(sshKeyPair.PrivateKey)
		defer sshAgent.Stop()

		awsPrivateIp := terraform.Output(t, validationResourcesTerraformOptions, "ec2_instance_private_ip")
		//awsPublicIp := terraform.Output(t, validationResourcesTerraformOptions, "ec2_instance_public_ip")
		//gcpPrivateIp := terraform.Output(t, validationResourcesTerraformOptions, "gce_instance_private_ip")
		gcpPublicIp := terraform.Output(t, validationResourcesTerraformOptions, "gce_instance_public_ip")

		expectedText := "4 packets transmitted, 4 received, 0% packet loss"
		//awsPingGcpCommand := "ping -c4 " + gcpPrivateIp
		gcpPingAwsCommand := "ping -c4 " + awsPrivateIp

		//awsSshHost := ssh.Host{
		//	Hostname:         awsPublicIp,
		//	SshUserName:      "ubuntu",
		//	OverrideSshAgent: sshAgent,
		//}
		gcpSshHost := ssh.Host{
			Hostname:         gcpPublicIp,
			SshUserName:      "terratest",
			OverrideSshAgent: sshAgent,
		}

		gcpToAwsOut := ssh.CheckSshCommandWithRetry(t, gcpSshHost, gcpPingAwsCommand, 10, 5*time.Second, ssh.CheckSshCommandE)
		assert.True(t, strings.Contains(gcpToAwsOut, expectedText))

		// SSH to AWS not working for the moment
		//awsToGcpOut := ssh.CheckSshCommandWithRetry(t, awsSshHost, awsPingGcpCommand, 10, 5*time.Second, ssh.CheckSshCommandE)
		//assert.True(t, strings.Contains(awsToGcpOut, expectedText))
	})
}
