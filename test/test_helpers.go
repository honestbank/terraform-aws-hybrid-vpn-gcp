package test

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/random"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"os"
	"os/user"
	"strings"
	"testing"
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

// getGoogleCredentials reads a static service account credentials JSON file named gcp-creds.json from the test/ folder
func getGoogleCredentials() string {
	envGoogleCredentials, envPresent := os.LookupEnv("GOOGLE_CREDENTIALS")
	if envPresent {
		return envGoogleCredentials
	}

	fileGoogleCredentials, errReadingGCredsFromFile := os.ReadFile("gcp-creds.json")
	if errReadingGCredsFromFile == nil {
		return string(fileGoogleCredentials)
	}
	panic("No Google credentials available")
}
