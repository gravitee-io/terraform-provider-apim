package acceptance_test

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Verifies the create, read, import, and delete lifecycle of the
// `apim_subscription` resource.
func TestSubscriptionResource_minimal(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_subscription.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviders(),
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"organization_id": config.StringVariable(organizationId),
				},
			},
			// Verifies resource import.
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"organization_id": config.StringVariable(organizationId),
				},
				ResourceName: resourceAddress,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					importIDBytes, err := json.Marshal(struct {
						EnvironmentId  string `json:"environment_id"`
						Hrid           string `json:"hrid"`
						OrganizationId string `json:"organization_id"`
					}{
						EnvironmentId:  s.RootModule().Resources[resourceAddress].Primary.Attributes["environment_id"],
						Hrid:           s.RootModule().Resources[resourceAddress].Primary.Attributes["hrid"],
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
	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_subscription.test"

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"ending_at":       config.StringVariable("2040-12-25T09:12:28Z"),
					"organization_id": config.StringVariable(organizationId),
				},
			},
			// Verifies resource import.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"ending_at":       config.StringVariable("2040-12-25T09:12:28Z"),
					"organization_id": config.StringVariable(organizationId),
				},
				ResourceName: resourceAddress,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					importIDBytes, err := json.Marshal(struct {
						EnvironmentId  string `json:"environment_id"`
						Hrid           string `json:"hrid"`
						OrganizationId string `json:"organization_id"`
					}{
						EnvironmentId:  s.RootModule().Resources[resourceAddress].Primary.Attributes["environment_id"],
						Hrid:           s.RootModule().Resources[resourceAddress].Primary.Attributes["hrid"],
						OrganizationId: s.RootModule().Resources[resourceAddress].Primary.Attributes["organization_id"],
					})

					return string(importIDBytes), err
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"starting_at",
				},
			},
			// Verifies resource update.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"ending_at":       config.StringVariable("2042-12-25T09:12:28Z"),
					"organization_id": config.StringVariable(organizationId),
				},
			},
			// Verifies resource update using an ending_at in a different format.
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"ending_at":       config.StringVariable("2043-12-25T10:12:28+00:00"),
					"organization_id": config.StringVariable(organizationId),
				},
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}
