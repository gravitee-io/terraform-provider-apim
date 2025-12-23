package acceptance_test

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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
func TestAPIV4Resource_all(t *testing.T) {
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
					"pages.1.homepage",
					"pages.2.homepage",
					"pages.3.homepage",
				},
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

// Verifies the update of the name attribute of the `apim_apiv4` resource.
func TestAPIV4Resource_plans(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_apiv4.test"

	plan1 := config.ObjectVariable(config.Variables{
		"hrid":        config.StringVariable("plan1"),
		"description": config.StringVariable("plan 1"),
		"mode":        config.StringVariable("STANDARD"),
		"name":        config.StringVariable("plan 1"),
		"status":      config.StringVariable("PUBLISHED"),
		"type":        config.StringVariable("API"),
		"validation":  config.StringVariable("AUTO"),
		"security": config.ObjectVariable(config.Variables{
			"type": config.StringVariable("KEY_LESS"),
		}),
	})
	plan2 := config.ObjectVariable(config.Variables{
		"hrid":        config.StringVariable("plan2"),
		"description": config.StringVariable("plan 2"),
		"mode":        config.StringVariable("STANDARD"),
		"name":        config.StringVariable("plan 2"),
		"status":      config.StringVariable("PUBLISHED"),
		"type":        config.StringVariable("API"),
		"validation":  config.StringVariable("AUTO"),
		"security": config.ObjectVariable(config.Variables{
			"type": config.StringVariable("MTLS"),
		}),
	})
	plan3 := config.ObjectVariable(config.Variables{
		"hrid":        config.StringVariable("plan3"),
		"description": config.StringVariable("plan 3"),
		"mode":        config.StringVariable("STANDARD"),
		"name":        config.StringVariable("plan 3"),
		"status":      config.StringVariable("PUBLISHED"),
		"type":        config.StringVariable("API"),
		"validation":  config.StringVariable("AUTO"),
		"security": config.ObjectVariable(config.Variables{
			"type": config.StringVariable("MTLS"),
		}),
	})
	plan2Deprecated := config.ObjectVariable(config.Variables{
		"hrid":        config.StringVariable("plan2"),
		"description": config.StringVariable("plan 2"),
		"mode":        config.StringVariable("STANDARD"),
		"name":        config.StringVariable("plan 2"),
		"status":      config.StringVariable("DEPRECATED"),
		"type":        config.StringVariable("API"),
		"validation":  config.StringVariable("AUTO"),
		"security": config.ObjectVariable(config.Variables{
			"type": config.StringVariable("MTLS"),
		}),
	})

	plans := config.ListVariable(plan1, plan2, plan3)
	plansDeprecated := config.ListVariable(plan1, plan3, plan2Deprecated)
	plansShuffled := config.ListVariable(plan3, plan1, plan2)
	plansDelete := config.ListVariable(plan1, plan2)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization_id": config.StringVariable(organizationId),
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"plans":           plans,
				},
			},
			// Verifies resource import.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization_id": config.StringVariable(organizationId),
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"plans":           plans,
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
					"organization_id": config.StringVariable(organizationId),
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"plans":           plansDeprecated,
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
			// Verifies resource update again.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization_id": config.StringVariable(organizationId),
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"plans":           plansShuffled,
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
			// Verifies resource update again.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization_id": config.StringVariable(organizationId),
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"plans":           plansDelete,
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

			// Testing framework implicitly verifies resource delete.
		},
	})
}

// Verifies the update of the name attribute of the `apim_apiv4` resource.
func TestAPIV4Resource_pages(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_apiv4.test"

	homePage := config.ObjectVariable(config.Variables{
		"hrid":        config.StringVariable("homepage"),
		"name":        config.StringVariable("Homepage"),
		"content":     config.StringVariable("# Homepage"),
		"homepage":    config.BoolVariable(true),
		"type":        config.StringVariable("MARKDOWN"),
		"parent_hrid": config.StringVariable(""),
	})

	folder := config.ObjectVariable(config.Variables{
		"hrid":        config.StringVariable("folder"),
		"name":        config.StringVariable("Pages"),
		"content":     config.StringVariable(""),
		"homepage":    config.BoolVariable(false),
		"type":        config.StringVariable("FOLDER"),
		"parent_hrid": config.StringVariable(""),
	})

	page := config.ObjectVariable(config.Variables{
		"hrid":        config.StringVariable("page"),
		"name":        config.StringVariable("Page"),
		"content":     config.StringVariable("This is the content"),
		"homepage":    config.BoolVariable(false),
		"type":        config.StringVariable("MARKDOWN"),
		"parent_hrid": config.StringVariable("folder"),
	})

	pages := config.ListVariable(homePage, folder, page)
	pagesShuffled := config.ListVariable(folder, page, homePage)
	pagesDelete := config.ListVariable(folder, homePage)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization_id": config.StringVariable(organizationId),
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"pages":           pages,
				},
			},
			// Verifies resource import.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization_id": config.StringVariable(organizationId),
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"pages":           pages,
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
					"organization_id": config.StringVariable(organizationId),
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"pages":           pagesShuffled,
				},
			},
			// Verifies resource update again.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization_id": config.StringVariable(organizationId),
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"pages":           pagesShuffled,
				},
			},
			// Verifies resource update again.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization_id": config.StringVariable(organizationId),
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"pages":           pagesDelete,
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
	updatedPageHrid := "general_conditions_v2"
	nonExistingPageHrid := "foo"
	unpublishedExistingPageHrid := "unpublished_general_conditions"

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
			// Check error on page hrid.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization_id":         config.StringVariable(organizationId),
					"environment_id":          config.StringVariable(environmentId),
					"hrid":                    config.StringVariable(randomId),
					"general_conditions_hrid": config.StringVariable(nonExistingPageHrid),
				},
				ExpectError: regexp.MustCompile("non existing page as general conditions"),
			},
			// Check error on unpublished page.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization_id":         config.StringVariable(organizationId),
					"environment_id":          config.StringVariable(environmentId),
					"hrid":                    config.StringVariable(randomId),
					"general_conditions_hrid": config.StringVariable(unpublishedExistingPageHrid),
				},
				ExpectError: regexp.MustCompile("non published page as general conditions"),
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
			// Update `allow_methods` and set '*'.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":          config.StringVariable(randomId),
					"allow_methods": config.ListVariable(config.StringVariable("*")),
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
