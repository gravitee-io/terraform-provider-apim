package examples_test

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gravitee-io/terraform-provider-apim/internal/provider"
	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/config"
)

const testIgnoreFile = ".testignore"
const tofuIgnoreFile = ".tofuignore"

type directory string

type testcase struct {
	name         string
	directory    directory
	skipVersions []utils.ApimVersion
}

// Compile-time check to ensure directory.get implements config.TestStepConfigFunc
var _ config.TestStepConfigFunc = directory("").get

func (d directory) get(config.TestStepConfigRequest) string {
	return string(d)
}

const examplesUseCasesPath = "../../examples/use-cases"
const examplesTutorialsPath = "../../examples/tutorials"

func listTestDirectoriesSkipping(basePath string, dirShouldBeSkipped func(string) bool) []string {
	var dirs []string
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != basePath {
			if dirShouldBeSkipped(path) {
				fmt.Printf("Skipping directory %s", path)
				return filepath.SkipDir
			}
			dirs = append(dirs, path)
		}
		return nil
	})
	if err != nil {
		panic(err.Error())
	}
	return dirs
}

// hasTestIgnoreWithoutVersions returns true if the directory has a .testignore
// file that contains no valid APIM version entries, meaning "skip entirely".
func hasTestIgnoreWithoutVersions(dir string) bool {
	path := filepath.Join(dir, testIgnoreFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return parseTestIgnore(path) == nil
}

// hasTestIgnoreWithoutVersions returns true if the directory has a .testignore .tofuignore
// file that contains no valid APIM version entries, meaning "skip entirely".
func hasTestIgnoreOrTofuIgnoreWithoutVersion(dir string) bool {
	if hasTestIgnoreWithoutVersions(dir) {
		return true
	}
	path := filepath.Join(dir, tofuIgnoreFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return parseTestIgnore(filepath.Join(dir, tofuIgnoreFile)) == nil
}

func createTestCases(directories []string) []testcase {
	cases := make([]testcase, 0)
	for _, dir := range directories {
		testDir := filepath.Base(dir)
		skipVersions := parseTestIgnore(filepath.Join(dir, testIgnoreFile))

		cases = append(cases, testcase{
			name:         testDir,
			directory:    directory(dir),
			skipVersions: skipVersions,
		})
	}
	return cases
}

// parseTestIgnore reads a .testignore/.tofuignore file and returns the APIM versions to skip.
// If the file does not exist or contains no valid versions, it returns nil (skip none).
// If the file contains only unknown version strings, they are ignored.
func parseTestIgnore(path string) []utils.ApimVersion {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var versions []utils.ApimVersion
	scanner := bufio.NewScanner(f)
	var ln int
	for scanner.Scan() {
		ln++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		version := utils.ParseApimVersion(line)
		if version == utils.ApimUnknown {
			panic(fmt.Sprintf("Invalid APIM version '%s' in %s:%d", line, path, ln))
		}
		versions = append(versions, version)
	}
	return versions
}

func cleanupTerraformStateFiles(directories []string) {
	for _, dir := range directories {
		tfFiles := []string{
			filepath.Join(dir, ".terraform"),
			filepath.Join(dir, "terraform.tfstate"),
			filepath.Join(dir, "terraform.tfstate.backup"),
		}

		for _, file := range tfFiles {
			if err := os.RemoveAll(file); err != nil {
				panic(fmt.Sprintf("Failed to remove terraform file %s: %v", file, err))
			}
		}
	}
}

// Returns a mapping of provider type names to provider server implementations,
// suitable for acceptance testing via the
// resource.TestCase.ProtoV6ProtocolFactories field.
func testProviders() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"apim": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
}
