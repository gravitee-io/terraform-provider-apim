package acceptance_test

import (
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Verifies the create, read, import, and delete lifecycle of the
// `apim_group` resource.
func TestGroupResource_minimal(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_group.test"

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
				},
			},
			// Verifies resource import.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
				},
				ResourceName:      resourceAddress,
				ImportState:       true,
				ImportStateIdFunc: importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id"}, nil),
				ImportStateVerify: true,
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

// Verifies the create, read, import, and delete lifecycle of the
// `apim_group` resource.
func TestGroupResource_all(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_group.test"

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
				},
			},
			// Verifies resource import.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
				},
				ResourceName:      resourceAddress,
				ImportState:       true,
				ImportStateIdFunc: importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id"}, nil),
				ImportStateVerify: true,
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

// Verifies the create, read, import, and delete lifecycle of the
// `apim_group` resource.
func TestGroupResource_update(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)

	name := config.StringVariable("The group")
	nameUpdated := config.StringVariable("The group updated")

	memberApi := config.ObjectVariable(config.Variables{
		"source":    config.StringVariable("memory"),
		"source_id": config.StringVariable("api1"),
		"roles": config.MapVariable(map[string]config.Variable{
			"API":         config.StringVariable("OWNER"),
			"APPLICATION": config.StringVariable("USER"),
			"INTEGRATION": config.StringVariable("USER"),
		}),
	})
	memberApp := config.ObjectVariable(config.Variables{
		"source":    config.StringVariable("memory"),
		"source_id": config.StringVariable("application1"),
		"roles": config.MapVariable(map[string]config.Variable{
			"API":         config.StringVariable("USER"),
			"APPLICATION": config.StringVariable("OWNER"),
			"INTEGRATION": config.StringVariable("USER"),
		}),
	})
	membersApi := config.ListVariable(memberApi)
	membersApp := config.ListVariable(memberApp)
	membersApiApp := config.ListVariable(memberApi, memberApp)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":    config.StringVariable(randomId),
					"name":    name,
					"members": membersApi,
				},
			},
			// Verifies resource name update.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":    config.StringVariable(randomId),
					"name":    nameUpdated,
					"members": membersApi,
				},
			},
			// Verifies resource members replace.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":    config.StringVariable(randomId),
					"name":    nameUpdated,
					"members": membersApp,
				},
			},
			// Verifies resource members re-add.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":    config.StringVariable(randomId),
					"name":    nameUpdated,
					"members": membersApiApp,
				},
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}
