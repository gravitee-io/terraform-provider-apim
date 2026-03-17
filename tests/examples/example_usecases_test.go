package examples_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUseCases_Examples(t *testing.T) {

	directories := listTestDirectories(examplesUseCasesPath)

	cases := createTestCases(directories)

	cleanupTerraformStateFiles(directories)
	t.Cleanup(func() { cleanupTerraformStateFiles(directories) })

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

func TestTutorials(t *testing.T) {
	cases := []struct {
		name      string
		directory func(config.TestStepConfigRequest) string
	}{{
		name: "application-b2b_multi_certs",
		directory: func(config.TestStepConfigRequest) string {
			return "../../examples/tutorials/application-b2b_multi_certs"
		},
	}, {
		name: "v4api-shared_resource_local",
		directory: func(config.TestStepConfigRequest) string {
			return "../../examples/tutorials/v4api-shared_resource_local"
		},
	}}

	emptyReq := config.TestStepConfigRequest{}

	t.Cleanup(func() {
		cleanupTerraformStateFiles([]string{cases[0].directory(emptyReq), cases[1].directory(emptyReq)})
	})

	providers := testProviders()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Running test case: %s", tc.name)

			resource.Test(t, resource.TestCase{
				Steps: []resource.TestStep{
					// Verifies resource create and read.
					{
						ProtoV6ProviderFactories: providers,
						ConfigDirectory:          tc.directory,
					},
				},
			})
			// Testing framework implicitly verifies resource delete.
		})
	}
}
