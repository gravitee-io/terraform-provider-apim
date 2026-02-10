package examples_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestApplicationResource_Examples(t *testing.T) {

	directories := listTestDirectories(examplesUseCasesPath)

	cases := createTestCases(directories)

	cleanupTerraformStateFiles(directories)

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
					},
				},
			})
			// Testing framework implicitly verifies resource delete.
		})
	}
}
