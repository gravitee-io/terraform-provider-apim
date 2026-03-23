package examples_test

import (
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUseCases_Examples(t *testing.T) {
	runExampleTests(t, examplesUseCasesPath)
}

func TestTutorials(t *testing.T) {
	runExampleTests(t, examplesTutorialsPath)
}

func runExampleTests(t *testing.T, basePath string) {
	directories := listTestDirectories(basePath)

	cases := createTestCases(directories)

	cleanupTerraformStateFiles(directories)
	t.Cleanup(func() { cleanupTerraformStateFiles(directories) })

	providers := testProviders()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.skipVersions) > 0 {
				utils.SkipFor(t, tc.skipVersions...)
			}

			t.Logf("Running test case: %s", tc.name)

			resource.Test(t, resource.TestCase{
				Steps: []resource.TestStep{
					{
						ProtoV6ProviderFactories: providers,
						ConfigDirectory:          tc.directory.get,
					},
				},
			})
		})
	}
}
