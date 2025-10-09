package acceptance_test

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Verifies the create, read, import, and delete lifecycle of a minimal
// `apim_apiv4` resource.
func TestAPIV4Resource_minimal(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_apiv4.test"

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
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
					"notify_members",
					"version",
				},
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

// Verifies the create, read, import, and delete lifecycle of the
// `apim_apiv4` resource.
// TODO fix bugs before enabling this one
//func TestAPIV4Resource_all(t *testing.T) {
//	t.Parallel()
//
//	environmentId := "DEFAULT"
//	organizationId := "DEFAULT"
//	randomId := "test-" + acctest.RandString(10)
//	resourceAddress := "apim_apiv4.test"
//
//	resource.Test(t, resource.TestCase{
//		Steps: []resource.TestStep{
//			// Verifies resource create and read.
//			{
//				ProtoV6ProviderFactories: testProviders(),
//				ConfigDirectory:          config.TestNameDirectory(),
//				ConfigVariables: config.Variables{
//					"environment_id":  config.StringVariable(environmentId),
//					"hrid":            config.StringVariable(randomId),
//					"organization_id": config.StringVariable(organizationId),
//				},
//			},
//			// Verifies resource import.
//			{
//				ProtoV6ProviderFactories: testProviders(),
//				ConfigDirectory:          config.TestNameDirectory(),
//				ConfigVariables: config.Variables{
//					"environment_id":  config.StringVariable(environmentId),
//					"hrid":            config.StringVariable(randomId),
//					"organization_id": config.StringVariable(organizationId),
//				},
//				ResourceName: resourceAddress,
//				ImportState:  true,
//				ImportStateIdFunc: func(s *terraform.State) (string, error) {
//					importIDBytes, err := json.Marshal(struct {
//						EnvironmentId  string `json:"environment_id"`
//						Hrid           string `json:"hrid"`
//						OrganizationId string `json:"organization_id"`
//					}{
//						EnvironmentId:  s.RootModule().Resources[resourceAddress].Primary.Attributes["environment_id"],
//						Hrid:           s.RootModule().Resources[resourceAddress].Primary.Attributes["hrid"],
//						OrganizationId: s.RootModule().Resources[resourceAddress].Primary.Attributes["organization_id"],
//					})
//
//					return string(importIDBytes), err
//				},
//				ImportStateVerify: true,
//				ImportStateVerifyIgnore: []string{
//					"notify_members",
//					"version",
//				},
//			},
//			// Testing framework implicitly verifies resource delete.
//		},
//	})
//}

// Verifies the update of the name attribute of the `apim_apiv4` resource.
func TestAPIV4Resource_update(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_apiv4.test"

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"name":            config.StringVariable(randomId + "-original"),
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
				ImportStateVerifyIgnore: []string{
					"notify_members",
					"version",
				},
			},
			// Verifies resource update.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
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

func TestAPIV4Resource_apikey(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"name":            config.StringVariable(randomId + "-original"),
					"organization_id": config.StringVariable(organizationId),
				},
				ExpectError: regexp.MustCompile("got: \"API_KEY\""),
			},
		},
	})
}

func TestAPIV4Resource_plan_general_conditions(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	pageHrid := "general_conditions"
	updatedPageHrid := "homepage"

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization_id":         config.StringVariable(organizationId),
					"environment_id":          config.StringVariable(environmentId),
					"hrid":                    config.StringVariable(randomId),
					"general_conditions_hrid": config.StringVariable(pageHrid),
				},
			},
			// Verifies resource update.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization_id":         config.StringVariable(organizationId),
					"environment_id":          config.StringVariable(environmentId),
					"hrid":                    config.StringVariable(randomId),
					"general_conditions_hrid": config.StringVariable(updatedPageHrid),
				},
			},
		},
	})
}

func TestAPIV4Resource_cors(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_apiv4.test"

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"organization_id": config.StringVariable(organizationId),
					"name":            config.StringVariable(randomId + "-original"),
				},
			},
			// Verifies resource import.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
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
					"notify_members",
					"version",
				},
			},
			// Update `allow_headers`.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":             config.StringVariable(randomId),
					"resource_headers": config.ListVariable(config.StringVariable("accept"), config.StringVariable("content-type")),
				},
			},
			// Update `allow_headers` order.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
					"resource_headers": config.ListVariable(
						config.StringVariable("content-type"),
						config.StringVariable("accept"),
					),
				},
			},
			// Update `allow_methods`.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":          config.StringVariable(randomId),
					"allow_methods": config.ListVariable(config.StringVariable("DELETE")),
				},
			},
			// Update `allow_methods`.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
					"allow_methods": config.ListVariable(
						config.StringVariable("GET"),
						config.StringVariable("DELETE"),
						config.StringVariable("OPTIONS"),
						config.StringVariable("PUT"),
						config.StringVariable("POST"),
					),
				},
			},
			// Update `allow_methods` and order.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
					"allow_methods": config.ListVariable(
						config.StringVariable("POST"),
						config.StringVariable("GET"),
						config.StringVariable("PUT"),
						config.StringVariable("DELETE"),
					),
				},
			},
			// Update `allow_origin`.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
					"allow_origin": config.ListVariable(
						config.StringVariable(`.*\\.acme\\.com`),
						config.StringVariable(`.*\\.example\\.com`),
						config.StringVariable(`.*\\.simple\\.com`),
					),
				},
			},
			// Update `allow_origin` and order.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
					"allow_origin": config.ListVariable(
						config.StringVariable(`.*\\.simple\\.com`),
						config.StringVariable(`.*\\.acme\\.com`),
					),
				},
			},
			// Update `expose_headers`.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
					"expose_headers": config.ListVariable(
						config.StringVariable("accept"),
						config.StringVariable("x-custom-header"),
						config.StringVariable("content-type"),
					),
				},
			},
			// Update `expose_headers` and order.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
					"expose_headers": config.ListVariable(
						config.StringVariable("content-type"),
						config.StringVariable("accept"),
					),
				},
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

func TestAPIV4Resource_dyn_props(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"organization_id": config.StringVariable(organizationId),
				},
			},
		},
	})

}
