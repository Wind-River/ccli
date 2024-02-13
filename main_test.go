package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

var logical_part_id string
var paths []string
var fvc []string

// TestPing is used to check the connection to the software part catalog's server
func TestPing(tester *testing.T) {
	// ccli ping
	cmd := exec.Command("ccli", "ping")
	// capturing command line output
	output, err := cmd.Output()
	if err != nil {
		tester.Error("failed to capture command line output", err)
	}
	// extracting the output message to be checked
	result := string(output)
	expected := "Ping Result: Success\n"
	if result != expected {
		tester.Errorf("Expected %s but got %s", expected, result)
	} else {
		paths = []string{"testdir/packages/openid-client-4.9.1_test.zip", "testdir/packages/busybox-1.35.0_test.zip"}
		fvc = []string{"46564332008de01dcc150bcf6673a576d4c438b442afbb61d2cc98017234e44d9e338f19e8", "4656433200fe9e83718aa74cf8e3af8517fbaa0c273ed31d32363278759f86ec45251b5599"}
	}
}

// TestUpload uploads the package present in the given path using the command line
// and checks if the command line output is as expected
func TestUpload(tester *testing.T) {
	// set timer for upload to 20 seconds
	duration := 180 * time.Second
	for i := 0; i < 1; i = i + 1 {
		// ccli upload testdir/packages/openid-client-4.9.1.zip
		cmd := exec.Command("ccli", "upload", paths[i])
		// capturing command line output
		output, err := cmd.Output()
		if err != nil {
			tester.Error("failed to capture command line output", err)
		}
		// extracting the output message to be checked
		result := string(output)
		expected := "Successfully uploaded package: " + paths[i] + "\n"
		if result != expected {
			tester.Errorf("Expected %s but got %s", expected, result)
		}

		var flag bool
		flag = false
		// Loop to check if the upload has been reflected on the catalog and
		for start := time.Now(); time.Since(start) < duration; {
			// ccli query 'query{part(file_verification_code:\"465643320044e55d9adee108307f5c274ecf14c4dd1442c43a66fc8955dcf7e40d6f8a50d1\"){size}}'
			cmd := exec.Command("ccli", "query", "query{part(file_verification_code:\""+fvc[i]+"\"){size}}")
			// capturing command line output
			output, err = cmd.Output()
			if err == nil {
				flag = true
				break
			}
		}
		if !flag {
			tester.Errorf("Timed out before upload could complete for %s", paths[i])
		}

	}
}

// TestAddPart adds a part based on the yml file present in the given path using the
// command line and checks if the command line output is as expected
func TestAddPart(tester *testing.T) {
	// ccli add --part testdir/yml/busybox-1.35.0.yml
	cmd := exec.Command("ccli", "add", "part", "testdir/yml/busybox-1.35.0.yml")
	// capturing command line output
	output, err := cmd.Output()
	if err != nil {
		tester.Error("failed to capture command line output", err)
	}
	// splitting and extracting the output message to be checked
	result := strings.Split(string(output), "\n")[0]
	// saving the part id outputted for it to be deleted later
	logical_part_id = strings.Fields(strings.Split(string(output), "\n")[2])[1]
	logical_part_id = logical_part_id[1 : len(logical_part_id)-2]
	expected := "Successfully added part from: testdir/yml/busybox-1.35.0.yml"
	if result != expected {
		tester.Errorf("Expected %s but got %s", expected, result)
	}
}

// TestQuery runs a qraphql query using the command line
// and checks if the command line output is as expected
func TestQuery(tester *testing.T) {
	// ccli query 'query{part(file_verification_code:\"465643320044e55d9adee108307f5c274ecf14c4dd1442c43a66fc8955dcf7e40d6f8a50d1\"){size}}'
	cmd := exec.Command("ccli", "query", "query{part(file_verification_code:\""+fvc[0]+"\"){size}}")
	// capturing command line output
	output, err := cmd.Output()
	if err != nil {
		tester.Error("failed to capture command line output", err)
	}
	// splitting and extracting the output message to be checked. Fields func is used to divide the string by whitespaces to extract specific value
	result := strings.Fields(string(output))[4]
	expected := "164900"
	if result != expected {
		tester.Errorf("Expected %s but got %s", expected, result)
	}

}

// TestUpdate updates a part's information based on the yml file present in the given path using the
// command line and checks if the command line output is as expected
func TestUpdate(tester *testing.T) {
	// ccli update --part testdir/yml/openid-client-4.9.1.yml
	cmd := exec.Command("ccli", "update", "testdir/yml/openid-client-4.9.1.yml")
	// capturing command line output
	output, err := cmd.Output()
	if err != nil {
		tester.Error("failed to capture command line output", err)
	}
	// splitting and extracting the output message to be checked
	result := strings.Split(string(output), "\n")[0]
	expected := "Part successfully updated"
	if result != expected {
		tester.Errorf("Expected %s but got %s", expected, result)
	}
}

// TestAddLicenseProfile adds a part's licensing profile based on the yml file present in the given path using the
// command line and checks if the command line output is as expected
func TestAddLicenseProfile(tester *testing.T) {
	// ccli add --profile testdir/yml/openid_license.yml
	cmd := exec.Command("ccli", "add", "profile", "testdir/yml/openid_license.yml")
	// capturing command line output
	output, err := cmd.Output()
	if err != nil {
		tester.Error("failed to capture command line output", err)
	}
	// splitting and extracting the output message to be checked. Fields func is used to divide the string by whitespaces to extract specific value
	result := strings.Fields(string(output))[0]
	expected := "Successfully"
	if result != expected {
		tester.Errorf("Expected %s but got %s", expected, result)
	}
}

// TestAddSecurityProfile adds a part's security profile based on the yml file present in the given path using the
// command line and checks if the command line output is as expected
func TestAddSecurityProfile(tester *testing.T) {
	// ccli add --profile testdir/yml/openid_security.yml
	cmd := exec.Command("ccli", "add", "profile", "testdir/yml/openid_security.yml")
	// capturing command line output
	output, err := cmd.Output()
	if err != nil {
		tester.Error("failed to capture command line output", err)
	}
	// splitting and extracting the output message to be checked. Fields func is used to divide the string by whitespaces to extract specific value
	result := strings.Fields(string(output))[0]
	expected := "Successfully"
	if result != expected {
		tester.Errorf("Expected %s but got %s", expected, result)
	}
}

// TestAddQualityProfile adds a part's quality profile based on the yml file present in the given path using the
// command line and checks if the command line output is as expected
func TestAddQualityProfile(tester *testing.T) {
	// ccli add --profile testdir/yml/openid_quality.yml
	cmd := exec.Command("ccli", "add", "profile", "testdir/yml/openid_quality.yml")
	// capturing command line output
	output, err := cmd.Output()
	if err != nil {
		tester.Error("failed to capture command line output", err)
	}
	// splitting and extracting the output message to be checked. Fields func is used to divide the string by whitespaces to extract specific value
	result := strings.Fields(string(output))[0]
	expected := "Successfully"
	if result != expected {
		tester.Errorf("Expected %s but got %s", expected, result)
	}
}

// TestExportPart exports a part to the given path in the form of a yml file using the
// command line and checks if the command line output is as expected
func TestExportPart(tester *testing.T) {
	// ccli export -fvc 465643320044e55d9adee108307f5c274ecf14c4dd1442c43a66fc8955dcf7e40d6f8a50d1 -o testdir/testpart.yml
	cmd := exec.Command("ccli", "export", "part", "fvc", fvc[0], "-o", "testdir/testpart.yml")
	// capturing command line output
	output, err := cmd.Output()
	if err != nil {
		tester.Error("failed to capture command line output", err)
	}
	// extracting the output message to be checked
	result := string(output)
	expected := "Part successfully exported to path: testdir/testpart.yml\n"
	if result != expected {
		tester.Errorf("Expected %s but got %s", expected, result)
	}
	// remove the exported test yml files
	os.RemoveAll("testdir/testpart.yml")
}

// TestExportTemplate exports a part or profile template to the given path in the form of a yml file using the
// command line and checks if the command line output is as expected
func TestExportTemplate(tester *testing.T) {
	// ccli export --template licensing -o testdir/testlicense.yml
	cmd := exec.Command("ccli", "export", "template", "license", "-o", "testdir/testlicense.yml")
	// capturing command line output
	output, err := cmd.Output()
	if err != nil {
		tester.Error("failed to capture command line output", err)
	}
	// extracting the output message to be checked
	result := string(output)
	expected := "Profile template successfully output to: testdir/testlicense.yml\n"
	if result != expected {
		tester.Errorf("Expected %s but got %s", expected, result)
	}
	// remove the exported test yml files
	os.RemoveAll("testdir/testlicense.yml")
}

// TestDelete first finds out a part's unique part-id using the file verification code and then deletes the part using the
// part-id and command line. Finally checks if the command line output is as expected
func TestDelete(tester *testing.T) {
	// ccli find -fvc 465643320044e55d9adee108307f5c274ecf14c4dd1442c43a66fc8955dcf7e40d6f8a50d1
	cmd := exec.Command("ccli", "find", "fvc", fvc[0])
	// capturing command line output
	outputfind, err := cmd.Output()
	if err != nil {
		tester.Error("failed to capture command line output", err)
	}
	// splitting and extracting the output message to be checked. Fields func is used to divide the string by whitespaces to extract specific value
	part_Id := strings.Fields(string(outputfind))[2]
	// ccli delete --id <part-id>
	cmd = exec.Command("ccli", "delete", part_Id)
	// capturing command line output
	outputDeletePart1, err := cmd.Output()
	if err != nil {
		tester.Error("failed to capture command line output", err)
	}
	// extracting the output message to be checked
	// e.g., Hello World\nHi world
	// split - Hi world
	result := string(outputDeletePart1)
	expected := "Successfully deleted id: " + part_Id + " from catalog\n"
	if result != expected {
		tester.Errorf("Expected %s but got %s", expected, result)
	}
}

/*
	func TestForceDelete(tester *testing.T) {
		// ccli find -fvc 465643320044e55d9adee108307f5c274ecf14c4dd1442c43a66fc8955dcf7e40d6f8a50d1
		cmd := exec.Command("ccli", "find", "-fvc", fvc[1])
		// capturing command line output
		outputfind, err := cmd.Output()
		if err != nil {
			tester.Error("failed to capture command line output", err)
		}
		// splitting and extracting the output message to be checked. Fields func is used to divide the string by whitespaces to extract specific value
		part_Id := strings.Fields(string(outputfind))[2]
		// ccli delete --id <part-id>
		cmd = exec.Command("ccli", "delete", "--id", part_Id, "--force")
		// capturing command line output
		output, err := cmd.Output()
		if err != nil {
			tester.Error("failed to capture command line output", err)
		}
		// extracting the output message to be checked
		// e.g., Hello World\nHi world
		// split - Hi world
		result := string(output)
		expected := "Successfully deleted id: " + part_Id + " from catalog\n"
		if result != expected {
			tester.Errorf("Expected %s but got %s", expected, result)
		}
	}
func TestRecursiveDelete(tester *testing.T) {
	// ccli find -fvc 465643320044e55d9adee108307f5c274ecf14c4dd1442c43a66fc8955dcf7e40d6f8a50d1
	cmd := exec.Command("ccli", "find", "fvc", fvc[1])
	// capturing command line output
	outputfind, err := cmd.Output()
	if err != nil {
		tester.Error("failed to capture command line output", err)
	}
	// splitting and extracting the output message to be checked. Fields func is used to divide the string by whitespaces to extract specific value
	part_Id := strings.Fields(string(outputfind))[2]
	// ccli delete --id <part-id>
	cmd = exec.Command("ccli", "delete", part_Id, "--recursive")
	// capturing command line output
	output, err := cmd.Output()
	if err != nil {
		tester.Error("failed to capture command line output", err)
	}
	// extracting the output message to be checked
	// e.g., Hello World\nHi world
	// split - Hi world
	result := string(output)
	expected := "Successfully deleted id: " + part_Id + " from catalog\n"
	if result != expected {
		tester.Errorf("Expected %s but got %s", expected, result)
	}
}
*/

/*
	func TestSharedFileDeletion(tester *testing.T) {
		cmd := exec.Command("ccli", "query", "query{file(sha256:\"f29bc64a9d3732b4b9035125fdb3285f5b6455778edca72414671e0ca3b2e0de\"){names}}")
		// capturing command line output
		_, err := cmd.Output()
		if err == nil {
			tester.Error("shared file could not be deleted", err)
		}
	}
*/
func TestArchiveDeletion(tester *testing.T) {
	cmd := exec.Command("ccli", "query", "query{archive(name:\"openid-client-4.9.1_test.zip\"){insert_date}}")
	// capturing command line output
	_, err := cmd.Output()
	if err == nil {
		tester.Error("archive could not be deleted", err)
	}
	cmd = exec.Command("ccli", "query", "query{archive(name:\"busybox-1.35.0_test.zip\"){insert_date}}")
	// capturing command line output
	_, err = cmd.Output()
	if err == nil {
		tester.Error("archive could not be deleted", err)
	}
}

func TestLogicalPartDelete(tester *testing.T) {
	// ccli delete --id <part-id>
	cmd := exec.Command("ccli", "delete", logical_part_id)
	// capturing command line output
	outputDeletePart2, err := cmd.Output()
	if err != nil {
		tester.Error("failed to capture command line output", err)
	}
	// extracting the output message to be checked
	result := string(outputDeletePart2)
	expected := "Successfully deleted id: " + logical_part_id + " from catalog\n"
	if result != expected {
		tester.Errorf("Expected %s but got %s", expected, result)
	}
}

// TestMain runs all the tests
func TestMain(mainTester *testing.M) {
	exitCode := mainTester.Run()
	os.Exit(exitCode)
}
