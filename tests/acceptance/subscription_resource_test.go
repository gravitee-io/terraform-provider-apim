package acceptance_test

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
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

// Verifies that plan_hrid, api_hrid, and application_hrid cannot be changed on an existing subscription.
func TestSubscriptionResource_immutable_fields(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9)
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)

	cert := getClientTLSCert(t)

	baseVars := config.Variables{
		"environment_id":     config.StringVariable(environmentId),
		"hrid":               config.StringVariable(randomId),
		"organization_id":    config.StringVariable(organizationId),
		"api_hrid":           config.StringVariable("api-" + randomId),
		"plan_hrid":          config.StringVariable("jwt"),
		"application_hrid":   config.StringVariable("app-" + randomId),
		"client_certificate": config.StringVariable(cert),
	}

	changePlanVars := config.Variables{
		"environment_id":     config.StringVariable(environmentId),
		"hrid":               config.StringVariable(randomId),
		"organization_id":    config.StringVariable(organizationId),
		"api_hrid":           config.StringVariable("api-" + randomId),
		"plan_hrid":          config.StringVariable("mtls"),
		"application_hrid":   config.StringVariable("app-" + randomId),
		"client_certificate": config.StringVariable(cert),
	}

	changeAPIVars := config.Variables{
		"environment_id":     config.StringVariable(environmentId),
		"hrid":               config.StringVariable(randomId),
		"organization_id":    config.StringVariable(organizationId),
		"api_hrid":           config.StringVariable("api-other-" + randomId),
		"plan_hrid":          config.StringVariable("jwt"),
		"application_hrid":   config.StringVariable("app-" + randomId),
		"client_certificate": config.StringVariable(cert),
	}

	changeAppVars := config.Variables{
		"environment_id":     config.StringVariable(environmentId),
		"hrid":               config.StringVariable(randomId),
		"organization_id":    config.StringVariable(organizationId),
		"api_hrid":           config.StringVariable("api-" + randomId),
		"plan_hrid":          config.StringVariable("jwt"),
		"application_hrid":   config.StringVariable("app-other-" + randomId),
		"client_certificate": config.StringVariable(cert),
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Step 1: Create the subscription.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          baseVars,
			},
			// Step 2: Try to change the plan — expect plan error.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          changePlanVars,
				ExpectError:              regexp.MustCompile("This attribute cannot be changed on an existing resource"),
			},
			// Step 3: Try to change the application — expect plan error.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          changeAppVars,
				ExpectError:              regexp.MustCompile("This attribute cannot be changed on an existing resource"),
			},
			// Step 4: Try to change the API — expect plan error.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          changeAPIVars,
				ExpectError:              regexp.MustCompile("This attribute cannot be changed on an existing resource"),
			},

			// Testing framework implicitly verifies resource delete.
		},
	})
}
