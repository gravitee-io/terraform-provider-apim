package acceptance_test

import (
	"regexp"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// What an import cannot be compared on: the platform reports discovered tools,
// prompts and resources in no stable order, and the provider does not read the
// authentication back, so an imported server has none until it is applied.
var catalogMcpServerNotImported = []string{"tools", "prompts", "resources", "server_connection.auth"}

func TestCatalogMcpServerResource_minimal(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	endpoint := utils.McpServerURL(t)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_catalog_mcp_server.test"
	variables := config.Variables{
		"hrid":     config.StringVariable(randomId),
		"endpoint": config.StringVariable(endpoint),
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// The platform discovers the server on create: what it found is in the state.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceAddress, "id"),
					resource.TestCheckResourceAttr(resourceAddress, "entity_id", "mcp-server."+randomId),
					resource.TestCheckResourceAttr(resourceAddress, "server_connection.transport", "HTTP"),
					resource.TestCheckResourceAttrSet(resourceAddress, "last_synced_at"),
					resource.TestCheckResourceAttrSet(resourceAddress, "tools.0.name"),
					resource.TestCheckResourceAttrSet(resourceAddress, "tools.0.entity_id"),
				),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables,
				ResourceName:             resourceAddress,
				ImportState:              true,
				ImportStateIdFunc:        importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id"}, nil),
				ImportStateVerify:        true,
				ImportStateVerifyIgnore:  catalogMcpServerNotImported,
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

func TestCatalogMcpServerResource_update(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	endpoint := utils.McpServerURL(t)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_catalog_mcp_server.test"
	variables := func(description string) config.Variables {
		return config.Variables{
			"hrid":        config.StringVariable(randomId),
			"endpoint":    config.StringVariable(endpoint),
			"description": config.StringVariable(description),
		}
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables("Registered by the acceptance tests"),
				Check:                    resource.TestCheckResourceAttr(resourceAddress, "description", "Registered by the acceptance tests"),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables("Described again"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "description", "Described again"),
					resource.TestCheckResourceAttrSet(resourceAddress, "tools.0.name"),
				),
			},
			// Re-applying the same configuration plans nothing.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables("Described again"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

func TestCatalogMcpServerResource_header_auth(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	endpoint := utils.McpServerURL(t)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_catalog_mcp_server.test"
	variables := config.Variables{
		"hrid":     config.StringVariable(randomId),
		"endpoint": config.StringVariable(endpoint),
		"value":    config.StringVariable("Bearer acceptance-test"),
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "server_connection.auth.header.name", "Authorization"),
					resource.TestCheckResourceAttr(resourceAddress, "server_connection.auth.header.value", "Bearer acceptance-test"),
				),
			},
			// The platform never returns the header value. The state keeps the
			// declared one across a read, so the same configuration plans nothing.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.TestCheckResourceAttr(resourceAddress, "server_connection.auth.header.value", "Bearer acceptance-test"),
			},
		},
	})
}

func TestCatalogMcpServerResource_immutable_fields(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	endpoint := utils.McpServerURL(t)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	variables := func(entityId string) config.Variables {
		return config.Variables{
			"hrid":      config.StringVariable(randomId),
			"endpoint":  config.StringVariable(endpoint),
			"entity_id": config.StringVariable(entityId),
		}
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables("mcp-server." + randomId),
			},
			// Authorization policies name the entity id: another one is refused when planning.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables("mcp-server." + randomId + "-renamed"),
				ExpectError:              regexp.MustCompile("This attribute cannot be changed on an existing resource"),
			},
		},
	})
}

func TestCatalogMcpServerResource_unreachable_upstream(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	variables := func(lookup bool) config.Variables {
		return config.Variables{
			"hrid":   config.StringVariable(randomId),
			"lookup": config.BoolVariable(lookup),
		}
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// The platform discovers the server before registering it: an upstream it
			// cannot reach fails the apply, and the finding names the discovery.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(false),
				ExpectError:              regexp.MustCompile(`(?i)discovery\s+failed`),
			},
			// Nothing was registered: reading the same hrid finds no server.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(true),
				ExpectError:              regexp.MustCompile(`404`),
			},
		},
	})
}
