package acceptance_test

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Verifies the create, read, import, and delete lifecycle of the
// `apim_subscription` resource.
func TestSubscriptionResource_minimal(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	apiRandomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_subscription.test"

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(apiRandomId),
					"organization_id": config.StringVariable(organizationId),
				},
			},
			// Verifies resource import.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(apiRandomId),
					"organization_id": config.StringVariable(organizationId),
				},
				ResourceName: resourceAddress,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					importIDBytes, err := json.Marshal(struct {
						EnvironmentId  string `json:"environment_id"`
						Hrid           string `json:"hrid"`
						ApiHrid        string `json:"api_hrid"`
						OrganizationId string `json:"organization_id"`
					}{
						EnvironmentId:  s.RootModule().Resources[resourceAddress].Primary.Attributes["environment_id"],
						Hrid:           "test",
						ApiHrid:        s.RootModule().Resources[resourceAddress].Primary.Attributes["api_hrid"],
						OrganizationId: s.RootModule().Resources[resourceAddress].Primary.Attributes["organization_id"],
					})

					return string(importIDBytes), err
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"starting_at",
				},
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

// Verifies the update ending_at of the name attribute of the `apim_subscription` resource.
func TestSubscriptionResource_update(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	apiRandomId := "test-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(apiRandomId),
					"ending_at":       config.StringVariable("2040-12-25T09:12:28Z"),
					"organization_id": config.StringVariable(organizationId),
				},
			},
			// Verifies resource update.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(apiRandomId),
					"ending_at":       config.StringVariable("2042-12-25T09:12:28Z"),
					"organization_id": config.StringVariable(organizationId),
				},
			},
			// Verifies resource update using an ending_at in a different format.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(apiRandomId),
					"ending_at":       config.StringVariable("2043-12-25T10:12:28+03:00"),
					"organization_id": config.StringVariable(organizationId),
				},
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

// Verifies create with metadata, update metadata, and plan stability of the `apim_subscription` resource.
func TestSubscriptionResource_metadata(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	apiRandomId := "test-" + acctest.RandString(10)

	metadataCreate := config.MapVariable(map[string]config.Variable{
		"key1": config.StringVariable("value1"),
		"key2": config.StringVariable("value2"),
	})

	metadataUpdate := config.MapVariable(map[string]config.Variable{
		"key1": config.StringVariable("updated"),
		"key3": config.StringVariable("value3"),
	})

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create with metadata.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(apiRandomId),
					"organization_id": config.StringVariable(organizationId),
					"metadata":        metadataCreate,
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Verifies resource update with different metadata.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(apiRandomId),
					"organization_id": config.StringVariable(organizationId),
					"metadata":        metadataUpdate,
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(apiRandomId),
					"organization_id": config.StringVariable(organizationId),
					"metadata":        config.MapVariable(map[string]config.Variable{}),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

// Verifies create, destroy, and recreate with metadata, checking plan stability.
func TestSubscriptionResource_metadata_recreate(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	apiRandomId := "test-" + acctest.RandString(10)

	metadata := config.MapVariable(map[string]config.Variable{
		"key1": config.StringVariable("value1"),
		"key2": config.StringVariable("value2"),
	})

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create with metadata.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(apiRandomId),
					"organization_id": config.StringVariable(organizationId),
					"metadata":        metadata,
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Destroy all resources.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(apiRandomId),
					"organization_id": config.StringVariable(organizationId),
					"metadata":        metadata,
				},
				Destroy: true,
			},
			// Recreate with same metadata and verify plan stability.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(apiRandomId),
					"organization_id": config.StringVariable(organizationId),
					"metadata":        metadata,
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}
