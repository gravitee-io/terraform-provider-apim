package examples_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
)

func TestExportHCL(t *testing.T) {

	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10)

	// Create API in APIM
	id := importAPI(t)

	// Copy and setup file to perform HCL export
	dir := prepareTerraformFiles(t, id)

	// clean after test
	defer os.RemoveAll(dir)
	defer destroyAPI(t, dir)

	// Run HCL export
	exportHCL(t, dir)

	// Replace API path in exported file to be able to apply it
	replacePathInExportedFile(t, dir)

	// Apply changes to APIM (new API)
	applyChanges(t, dir)

	// No diffs after apply
	verifyTerraformPlan(t, dir)

}

func importAPI(t *testing.T) string {

	serverURL := os.Getenv("APIM_SERVER_URL")
	serverURL = strings.ReplaceAll(serverURL, "/automation", "")

	username := os.Getenv("APIM_USERNAME")
	password := os.Getenv("APIM_PASSWORD")

	payload, err := os.ReadFile("../../examples/tutorials/hcl-export/exported.json")
	if err != nil {
		t.Fatalf("Failed to read exported.json: %v", err)
	}

	req, err := http.NewRequest("POST", serverURL+"/management/v2/environments/DEFAULT/apis/_import/definition", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Fatalf("Request failed with status: %d, body: %s", resp.StatusCode, string(body))
	}

	type importResponse struct {
		ID string `json:"id"`
	}
	var response importResponse
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	return response.ID
}

func exportHCL(t *testing.T, dir string) {
	runTerraformCommand(t, dir, "plan", "-generate-config-out=exported.tf")
}

func replacePathInExportedFile(t *testing.T, dir string) {
	exportedFilePath := filepath.Join(dir, "exported.tf")

	content, err := os.ReadFile(exportedFilePath)
	if err != nil {
		t.Fatalf("Failed to read exported.tf: %v", err)
	}

	modifiedContent := strings.ReplaceAll(string(content), "/terraform/exported-example/", "/terraform/exported-example-reimported/")

	if err := os.WriteFile(exportedFilePath, []byte(modifiedContent), 0644); err != nil {
		t.Fatalf("Failed to write modified exported.tf: %v", err)
	}
}

func applyChanges(t *testing.T, dir string) {
	// Delete import.tf
	importFilePath := filepath.Join(dir, "import.tf")
	if err := os.Remove(importFilePath); err != nil {
		t.Fatalf("Failed to delete import.tf: %v", err)
	}

	// Modify main.tf to remove X-Gravitee-Set-Hrid header
	mainFilePath := filepath.Join(dir, "main.tf")
	content, err := os.ReadFile(mainFilePath)
	if err != nil {
		t.Fatalf("Failed to read main.tf: %v", err)
	}

	modifiedContent := strings.ReplaceAll(string(content), `"X-Gravitee-Set-Hrid" = "true"`, "")

	if err := os.WriteFile(mainFilePath, []byte(modifiedContent), 0644); err != nil {
		t.Fatalf("Failed to write modified main.tf: %v", err)
	}

	// Run terraform apply
	runTerraformCommand(t, dir, "apply", "-auto-approve")
}

func verifyTerraformPlan(t *testing.T, dir string) {
	cmd := exec.Command("terraform", "plan", "-detailed-exitcode")
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()

	// Exit code 0: No changes, Exit code 1: Error, Exit code 2: Changes present
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			if exitCode != 0 {
				t.Fatalf("terraform plan failed with exit code %d\nOutput: %s", exitCode, string(output))
			}
		} else {
			t.Fatalf("Failed to run terraform plan: %v\nOutput: %s", err, string(output))
		}
	}
}

func destroyAPI(t *testing.T, dir string) {
	runTerraformCommand(t, dir, "destroy", "-auto-approve")
}

func prepareTerraformFiles(t *testing.T, apiID string) string {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "terraform-import-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Copy main.tf
	mainTfContent, err := os.ReadFile("../../examples/tutorials/hcl-export/main.tf")
	if err != nil {
		t.Fatalf("Failed to read main.tf: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "main.tf"), mainTfContent, 0644); err != nil {
		t.Fatalf("Failed to write main.tf to temp directory: %v", err)
	}

	// Copy and modify import.tf
	importTfContent, err := os.ReadFile("../../examples/tutorials/hcl-export/import.tf")
	if err != nil {
		t.Fatalf("Failed to read import.tf: %v", err)
	}
	modifiedContent := strings.ReplaceAll(string(importTfContent), "<<API ID>>", apiID)
	if err := os.WriteFile(filepath.Join(tempDir, "import.tf"), []byte(modifiedContent), 0644); err != nil {
		t.Fatalf("Failed to write import.tf to temp directory: %v", err)
	}

	return tempDir
}

func runTerraformCommand(t *testing.T, dir string, args ...string) []byte {
	cmd := exec.Command("terraform", args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command '%s' failed: %v\nOutput: %s", cmd.String(), err, string(output))
	}
	return output
}
