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
// `apim_application` resource.
func TestApplicationResource_minimal(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_application.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviders(),
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
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
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

// Verifies the create, read, import, and delete lifecycle of the
// `apim_application` resource with as many fields as possible
func TestApplicationResource_all(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_application.test"

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
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

// Verifies the update the name of `apim_application` resource.
func TestApplicationResource_update(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_application.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviders(),
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
					"name": config.StringVariable(randomId + "-original"),
				},
			},
			// Verifies resource import.
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"name":            config.StringVariable(randomId + "-original"),
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
			},
			// Verifies resource update.
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"name":            config.StringVariable(randomId + "-updated"),
					"organization_id": config.StringVariable(organizationId),
				},
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}
