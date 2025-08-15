package examples_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type directory string

type testcase struct {
	name            string
	hrid            string
	directory       directory
	resourceAddress string
	updateValue     string
	updateField     string
	importIgnore    []string
}

func (d directory) get(_ config.TestStepConfigRequest) string {
	return string(examplesUseCasesPath + "/" + d)
}

const examplesUseCasesDir = "examples/use-cases"
const examplesUseCasesPath = "../../" + examplesUseCasesDir

func TestApplicationResource_Examples(t *testing.T) {

	directories := listTestDirectories(examplesUseCasesPath)

	cases := createTestCases(t, directories)

	cleanupTerraformStateFiles(directories)

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"

	providers := testProviders()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Running test case: %s", tc.name)

			resource.Test(t, resource.TestCase{
				Steps: []resource.TestStep{
					// Verifies resource create and read.
					{
						ProtoV6ProviderFactories: providers,
						ConfigDirectory:          tc.directory.get,
						ConfigVariables: config.Variables{
							"hrid": config.StringVariable(tc.hrid),
						},
					},
					/* TODO Verifies resource import.
					{
						ProtoV6ProviderFactories: providers,
						ConfigDirectory:          tc.directory.get,
						ConfigVariables: config.Variables{
							"environment_id":  config.StringVariable(environmentId),
							"hrid":            config.StringVariable(tc.hrid),
							"organization_id": config.StringVariable(organizationId),
						},
						ResourceName: tc.resourceAddress,
						ImportState:  true,
						ImportStateIdFunc: func(s *terraform.State) (string, error) {
							importIDBytes, err := json.Marshal(struct {
								EnvironmentId  string `json:"environment_id"`
								Hrid           string `json:"hrid"`
								OrganizationId string `json:"organization_id"`
							}{
								EnvironmentId:  s.RootModule().Resources[tc.resourceAddress].Primary.Attributes["environment_id"],
								Hrid:           s.RootModule().Resources[tc.resourceAddress].Primary.Attributes["hrid"],
								OrganizationId: s.RootModule().Resources[tc.resourceAddress].Primary.Attributes["organization_id"],
							})

							return string(importIDBytes), err
						},
						ImportStateVerify: true,
						ImportStateVerifyIgnore: tc.importIgnore,
					}, */
					// Verifies resource update.
					{
						ProtoV6ProviderFactories: providers,
						ConfigDirectory:          tc.directory.get,
						ConfigVariables: config.Variables{
							"environment_id":  config.StringVariable(environmentId),
							"hrid":            config.StringVariable(tc.hrid),
							tc.updateField:    config.StringVariable(tc.updateValue),
							"organization_id": config.StringVariable(organizationId),
						},
					},
				},
			})

		})
	}
}

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

func createTestCases(t *testing.T, directories []string) []testcase {
	cases := make([]testcase, 0)
	for _, dir := range directories {

		if strings.HasSuffix(dir, examplesUseCasesDir) {
			continue
		}
		testDir := filepath.Base(dir) // Get the parent directory name

		// Extract resource type and id from the directory name
		typeAndId := strings.Split(testDir, "-")
		if len(typeAndId) != 2 {
			panic("Invalid directory name: " + dir + ". Should be of the form <type>-<id> where id does not contain spaces hyphens")
		}

		// Setup test case
		hrid := typeAndId[1]
		resourceAddress := "apim_" + typeAndId[0] + "." + hrid
		updateField := "name"
		updateValue := "updated-" + hrid
		importIgnore := []string{}
		if typeAndId[0] == "subscription" {
			updateField = "ending_at"
			updateValue = "2050-12-25T09:12:28Z"
			importIgnore = []string{"starting_at"}
		}

		cases = append(cases, testcase{
			name:            testDir,
			directory:       directory(testDir),
			resourceAddress: resourceAddress,
			hrid:            hrid,
			updateField:     updateField,
			updateValue:     updateValue,
			importIgnore:    importIgnore,
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
