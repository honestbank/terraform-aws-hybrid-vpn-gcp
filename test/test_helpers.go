package test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/random"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"os"
	"os/user"
	"strings"
)

// runID() Creates a unique ID suitable for including in the name of any cloud resource such that tests can be
// parallelized.  However, if there is an ENV set with the prefix "SKIP_", then assume we are working on individual tests
// locally and use the current username as the id.
func runID() (string, error) {
	if test_structure.SkipStageEnvVarSet() {
		currentUser, err := user.Current()
		if err != nil {
			return "", err
		}
		return strings.ToLower(currentUser.Username), nil
	}

	return strings.ToLower(random.UniqueId()), nil
}

// copySupportingFiles copies one or more files from the test/ dir into a destination dir.
// This is done to configure providers when using modules without them explicitly defined.
func copySupportingFiles(t *testing.T, fileNames []string, destination string) {
	testFileSourceDir, getTestDirErr := os.Getwd()
	if getTestDirErr != nil {
		fmt.Println("Calling t.FailNow(): could not execute os.Getwd(): ", getTestDirErr)
		t.FailNow()
	}

	fmt.Println("Test working directory is: ", testFileSourceDir)

	fmt.Println("Copying files: ", fileNames, " to temporary test dir: ", destination)
	for _, file := range fileNames {
		src := testFileSourceDir + "/" + file
		dest := destination + "/" + file
		copyErr := files.CopyFile(src, dest)
		if copyErr != nil {
			fmt.Println("😩 Calling t.FailNow(): failed copying from: ", src, " to: ", dest, " with error: ", copyErr)
			t.FailNow()
		} else {
			fmt.Println("✌️ Success! Copied from: ", src, " to: ", dest)
		}
	}
}

// cleanupSupportingFiles deletes one or more files from a directory, intended to be called after copySupportingFiles
func cleanupSupportingFiles(fileNames []string, destination string) error {
	fmt.Println("Deleting files: ", fileNames, "from dir: ", destination)
	for _, file := range fileNames {
		fullPath := destination + "/" + file
		removeErr := os.Remove(fullPath)
		if removeErr != nil {
			fmt.Println("😩 Failed deleting file ", fullPath, " with error: ", removeErr)
			return removeErr
		} else {
			fmt.Println("✌️ Success! Deleted file: ", fullPath)
		}
	}
	return nil
}

// validateGCPCredentials performs a quick preflight check to verify that GCP credentials are valid
// and have access to the target project with Compute Engine API enabled.
func validateGCPCredentials(t *testing.T, credentialsJSON string, project string) {
	t.Helper()

	// 1. Validate JSON structure and extract metadata
	var creds map[string]interface{}
	if err := json.Unmarshal([]byte(credentialsJSON), &creds); err != nil {
		t.Fatalf("GCP credentials preflight: invalid JSON: %v", err)
	}

	credType, _ := creds["type"].(string)
	if credType == "" {
		t.Fatalf("GCP credentials preflight: missing 'type' field in credentials JSON")
	}

	switch credType {
	case "service_account":
		fmt.Printf("GCP credentials preflight: type=%s, client_email=%s, project_in_key=%s, target_project=%s\n",
			credType, creds["client_email"], creds["project_id"], project)
	case "authorized_user":
		fmt.Printf("GCP credentials preflight: type=%s, client_id=%s, target_project=%s\n",
			credType, creds["client_id"], project)
	default:
		fmt.Printf("GCP credentials preflight: type=%s, target_project=%s\n", credType, project)
	}

	// 2. Verify we can obtain an access token
	ctx := context.Background()
	googleCreds, err := google.CredentialsFromJSON(ctx, []byte(credentialsJSON), compute.ComputeScope)
	if err != nil {
		t.Fatalf("GCP credentials preflight: failed to parse credentials: %v", err)
	}
	token, err := googleCreds.TokenSource.Token()
	if err != nil {
		t.Fatalf("GCP credentials preflight: failed to obtain access token: %v", err)
	}
	fmt.Printf("GCP credentials preflight: successfully obtained access token (expires: %s)\n", token.Expiry)

	// 3. Verify Compute Engine API access on target project
	svc, err := compute.NewService(ctx, option.WithCredentialsJSON([]byte(credentialsJSON)))
	if err != nil {
		t.Fatalf("GCP credentials preflight: failed to create Compute client: %v", err)
	}
	_, err = svc.Projects.Get(project).Do()
	if err != nil {
		t.Fatalf("GCP credentials preflight: cannot access project %q via Compute API (is the API enabled? does the SA have access?): %v", project, err)
	}
	fmt.Printf("GCP credentials preflight: successfully accessed Compute API for project %q\n", project)
}

// getGoogleCredentials reads GCP credentials from (in order):
// 1. GOOGLE_CREDENTIALS environment variable
// 2. gcp-creds.json file in the test/ folder
// 3. gcloud application default credentials
func getGoogleCredentials() string {
	envGoogleCredentials, envPresent := os.LookupEnv("GOOGLE_CREDENTIALS")
	if envPresent {
		fmt.Println("Using Google credentials from GOOGLE_CREDENTIALS environment variable")
		return envGoogleCredentials
	}

	fileGoogleCredentials, errReadingGCredsFromFile := os.ReadFile("gcp-creds.json")
	if errReadingGCredsFromFile == nil {
		fmt.Println("Using Google credentials from gcp-creds.json file")
		return string(fileGoogleCredentials)
	}

	// Fall back to gcloud application default credentials for local development
	home, err := os.UserHomeDir()
	if err == nil {
		adcPath := home + "/.config/gcloud/application_default_credentials.json"
		adcCredentials, adcErr := os.ReadFile(adcPath)
		if adcErr == nil {
			fmt.Println("Using Google credentials from gcloud application default credentials")
			return string(adcCredentials)
		}
	}

	panic("No Google credentials available")
}

// getAWSAssumeRoleARN reads optional AWS role ARN for provider assume_role.
// If unset, tests will use ambient AWS credentials directly.
func getAWSAssumeRoleARN() string {
	return os.Getenv("AWS_ASSUME_ROLE_ARN")
}
