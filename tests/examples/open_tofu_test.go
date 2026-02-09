package examples_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestApplicationResource_Examples_OpenTofu(t *testing.T) {

	directories := listTestDirectories(examplesUseCasesPath)

	cases := createTestCases(directories)

	cleanupTerraformStateFiles(directories)
	t.Cleanup(func() { cleanupTerraformStateFiles(directories) })

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Running OpenTofu test case: %s", tc.name)

			testDir := filepath.Join(examplesUseCasesPath, string(tc.directory))

			// Run OpenTofu apply
			if err := runOpenTofuCommand(testDir, "apply", "-auto-approve"); err != nil {
				t.Fatalf("OpenTofu apply failed: %v", err)
			}

			// Run OpenTofu to check plan
			if err := runOpenTofuCommand(testDir, "plan", "-detailed-exitcode"); err != nil {
				t.Fatalf("OpenTofu plan failed: %v", err)
			}
			// Clean up - destroy resources
			if err := runOpenTofuCommand(testDir, "destroy", "-auto-approve"); err != nil {
				t.Logf("OpenTofu destroy failed (non-fatal): %v", err)
			}
		})
	}
}

func runOpenTofuCommand(dir string, args ...string) error {
	ctx := context.Background()

	// Check if tofu is available, fallback to terraform
	binary := "tofu"
	if _, err := exec.LookPath("tofu"); err != nil {
		return fmt.Errorf("'tofu' binary not found in PATH")
	}

	cmd := exec.CommandContext(ctx, binary, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"TF_LOG=DEBUG", // Enable debug logging
		"TF_LOG_PATH="+filepath.Join(dir, "tofu.log"),
	)

	// Capture output for debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command '%s %s' failed: %v\nOutput: %s",
			binary, strings.Join(args, " "), err, string(output))
	}

	return nil
}
