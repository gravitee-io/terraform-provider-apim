package examples_test

import (
	"fmt"
	"github.com/gravitee-io/terraform-provider-apim/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"os"
	"path/filepath"
	"strings"
)

type directory string

type testcase struct {
	name            string
	hrid            string
	directory       directory
	resourceAddress string
}

func (d directory) get(config.TestStepConfigRequest) string {
	return string(examplesUseCasesPath + "/" + d)
}

const examplesUseCasesDir = "examples/use-cases"
const examplesUseCasesPath = "../../" + examplesUseCasesDir

func listTestDirectories(basePath string) []string {
	var dirs []string
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if _, err := os.Stat(filepath.Join(path, ".testignore")); os.IsNotExist(err) {
				dirs = append(dirs, path)
			}
		}
		return nil
	})
	if err != nil {
		panic(err.Error())
	}
	return dirs
}

func createTestCases(directories []string) []testcase {
	cases := make([]testcase, 0)
	for _, dir := range directories {

		if strings.HasSuffix(dir, examplesUseCasesDir) {
			continue
		}
		testDir := filepath.Base(dir) // Get the parent directory name

		// Extract resource type and id from the directory name
		typeAndId := strings.Split(testDir, "-")
		if len(typeAndId) != 2 {
			panic("Invalid directory name: " + testDir + ". Should be of the form <type>-<id> where id does not contain spaces hyphens")
		}

		// Setup test case
		hrid := typeAndId[1]
		resourceAddress := "apim_" + typeAndId[0] + "." + hrid

		cases = append(cases, testcase{
			name:            testDir,
			directory:       directory(testDir),
			resourceAddress: resourceAddress,
			hrid:            hrid,
		})

	}
	return cases
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
