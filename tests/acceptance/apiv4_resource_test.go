package acceptance_test

import (
	"regexp"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
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
				ResourceName:      resourceAddress,
				ImportState:       true,
				ImportStateIdFunc: importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id"}, nil),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"notify_members",
					"version",
					"listeners.0.http.paths.0.override_access",
				},
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

// Verifies explicit clearing of HTTP listener settings converges after refresh.
func TestAPIV4Resource_clearComputedListenerSettings(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_apiv4.test"
	clearCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceAddress, "listeners.0.http.path_mappings.#", "0"),
		resource.TestCheckResourceAttr(resourceAddress, "listeners.0.http.request_validation.reject_null_byte", "false"),
	)
	resource.Test(t, resource.TestCase{Steps: []resource.TestStep{
		{
			ProtoV6ProviderFactories: testProviders(),
			ConfigDirectory:          config.TestNameDirectory(),
			ConfigVariables: config.Variables{
				"environment_id":  config.StringVariable(environmentId),
				"hrid":            config.StringVariable(randomId),
				"organization_id": config.StringVariable(organizationId),
				"clear":           config.BoolVariable(false),
			},
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(resourceAddress, "listeners.0.http.path_mappings.#", "1"),
				resource.TestCheckResourceAttr(resourceAddress, "listeners.0.http.request_validation.reject_null_byte", "true"),
			),
		},
		{
			ProtoV6ProviderFactories: testProviders(),
			ConfigDirectory:          config.TestNameDirectory(),
			ConfigVariables: config.Variables{
				"environment_id":  config.StringVariable(environmentId),
				"hrid":            config.StringVariable(randomId),
				"organization_id": config.StringVariable(organizationId),
				"clear":           config.BoolVariable(true),
			},
			Check: clearCheck,
		},
		{
			ProtoV6ProviderFactories: testProviders(),
			ConfigDirectory:          config.TestNameDirectory(),
			ConfigVariables: config.Variables{
				"environment_id":  config.StringVariable(environmentId),
				"hrid":            config.StringVariable(randomId),
				"organization_id": config.StringVariable(organizationId),
				"clear":           config.BoolVariable(true),
			},
			PlanOnly: true,
			Check:    clearCheck,
		},
	}})
}

// Verifies an explicit empty selection clears it and remains stable when omitted.
func TestAPIV4Resource_clearConnectionEvents(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()

	resourceAddress := "apim_apiv4.test"
	randomId := "test-" + acctest.RandString(10)
	emptyEventsCheck := resource.TestCheckResourceAttr(resourceAddress, "analytics.connection_events.#", "0")
	resource.Test(t, resource.TestCase{Steps: []resource.TestStep{
		{
			ProtoV6ProviderFactories: testProviders(),
			ConfigDirectory:          config.TestNameDirectory(),
			ConfigVariables: config.Variables{
				"hrid":         config.StringVariable(randomId),
				"omit_events":  config.BoolVariable(false),
				"empty_events": config.BoolVariable(false),
			},
			Check: resource.TestCheckResourceAttr(resourceAddress, "analytics.connection_events.#", "1"),
		},
		{
			ProtoV6ProviderFactories: testProviders(),
			ConfigDirectory:          config.TestNameDirectory(),
			ConfigVariables: config.Variables{
				"hrid":         config.StringVariable(randomId),
				"omit_events":  config.BoolVariable(false),
				"empty_events": config.BoolVariable(true),
			},
			Check: emptyEventsCheck,
		},
		{
			ProtoV6ProviderFactories: testProviders(),
			ConfigDirectory:          config.TestNameDirectory(),
			ConfigVariables: config.Variables{
				"hrid":         config.StringVariable(randomId),
				"omit_events":  config.BoolVariable(true),
				"empty_events": config.BoolVariable(false),
			},
			PlanOnly: true,
			Check:    emptyEventsCheck,
		},
	}})
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
				ResourceName:      resourceAddress,
				ImportState:       true,
				ImportStateIdFunc: importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id"}, nil),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"notify_members",
					"version",
					"listeners.0.http.paths.0.override_access",
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
				ResourceName:      resourceAddress,
				ImportState:       true,
				ImportStateIdFunc: importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id"}, nil),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"notify_members",
					"version",
					"listeners.0.http.paths.0.override_access",
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
				ResourceName:      resourceAddress,
				ImportState:       true,
				ImportStateIdFunc: importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id"}, nil),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"notify_members",
					"version",
					"listeners.0.http.paths.0.override_access",
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
				ResourceName:      resourceAddress,
				ImportState:       true,
				ImportStateIdFunc: importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id"}, nil),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"notify_members",
					"version",
					"listeners.0.http.paths.0.override_access",
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

// Verifies that setting an empty map (`{}`) for
// `response_templates.<template>.<content-type>.headers` converges (no perpetual plan drift).
// This is the behavior described in GKO-3100.
func TestAPIV4Resource_responseTemplates_emptyHeaders_map_idempotency(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Create.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"organization_id": config.StringVariable(organizationId),
				},
			},
			// Re-apply the same config: plan must be empty (idempotency check).
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"organization_id": config.StringVariable(organizationId),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// Verifies that a listener path without trailing slash does not cause permanent drift.
// APIM normalizes paths by adding a trailing slash, and the provider must suppress
// the resulting diff to maintain idempotency (GKO-2621).
func TestAPIV4Resource_path_idempotency(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Create with path without trailing slash.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"organization_id": config.StringVariable(organizationId),
				},
			},
			// Re-apply the same config: plan must be empty (idempotency check).
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"organization_id": config.StringVariable(organizationId),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAPIV4Resource_portal_navigation(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_apiv4.test"

	navigationInitial := config.ListVariable(
		config.ObjectVariable(config.Variables{
			"path":         config.StringVariable("/guides" + randomId),
			"display_name": config.StringVariable("Guides"),
			"order":        config.IntegerVariable(1),
			"visibility":   config.StringVariable("PRIVATE"),
		}),
		config.ObjectVariable(config.Variables{
			"path":         config.StringVariable("/guides" + randomId + "/getting-started"),
			"display_name": config.StringVariable("Getting Started"),
			"order":        config.IntegerVariable(1),
		}),
	)
	navigationUpdated := config.ListVariable(
		config.ObjectVariable(config.Variables{
			"path":         config.StringVariable("/reference" + randomId),
			"display_name": config.StringVariable("Reference"),
			"order":        config.IntegerVariable(1),
			"visibility":   config.StringVariable("PUBLIC"),
		}),
		config.ObjectVariable(config.Variables{
			"path":         config.StringVariable("/guides" + randomId),
			"display_name": config.StringVariable("Guides"),
			"order":        config.IntegerVariable(3),
			"visibility":   config.StringVariable("PRIVATE"),
		}),
	)
	navigationReduced := config.ListVariable(
		config.ObjectVariable(config.Variables{
			"path":         config.StringVariable("/guides" + randomId),
			"display_name": config.StringVariable("Guides"),
			"order":        config.IntegerVariable(1),
			"visibility":   config.StringVariable("PUBLIC"),
		}),
	)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":    config.StringVariable(environmentId),
					"hrid":              config.StringVariable(randomId),
					"organization_id":   config.StringVariable(organizationId),
					"portal_navigation": navigationInitial,
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "portal_navigation.#", "2"),
					resource.TestCheckResourceAttr(resourceAddress, "portal_navigation.0.visibility", "PRIVATE"),
					resource.TestCheckNoResourceAttr(resourceAddress, "portal_navigation.1.visibility"),
					checkNavigationVisibility(resourceAddress, true, "Getting Started", "PRIVATE"),
					resource.TestCheckResourceAttr(resourceAddress, "portal_navigation.1.path", "/guides"+randomId+"/getting-started"),
				),
			},
			{

				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":    config.StringVariable(environmentId),
					"hrid":              config.StringVariable(randomId),
					"organization_id":   config.StringVariable(organizationId),
					"portal_navigation": navigationInitial,
				},
				ResourceName:      resourceAddress,
				ImportState:       true,
				ImportStateIdFunc: importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id"}, nil),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"notify_members",
					"listeners.0.http.paths.0.override_access",
				},
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":    config.StringVariable(environmentId),
					"hrid":              config.StringVariable(randomId),
					"organization_id":   config.StringVariable(organizationId),
					"portal_navigation": navigationUpdated,
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "portal_navigation.#", "2"),
					resource.TestCheckResourceAttr(resourceAddress, "portal_navigation.0.path", "/reference"+randomId),
					resource.TestCheckResourceAttr(resourceAddress, "portal_navigation.0.order", "1"),
					resource.TestCheckResourceAttr(resourceAddress, "portal_navigation.0.visibility", "PUBLIC"),
					resource.TestCheckResourceAttr(resourceAddress, "portal_navigation.1.order", "3"),
					resource.TestCheckResourceAttr(resourceAddress, "portal_navigation.1.visibility", "PRIVATE"),
				),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":    config.StringVariable(environmentId),
					"hrid":              config.StringVariable(randomId),
					"organization_id":   config.StringVariable(organizationId),
					"portal_navigation": navigationReduced,
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "portal_navigation.#", "1"),
					resource.TestCheckResourceAttr(resourceAddress, "portal_navigation.0.visibility", "PUBLIC"),
				),
			},
		},
	})
}

func TestAPIV4Resource_console_notification(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	groupHrid := randomId + "-notif-group"

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Create with notification group and events.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":                config.StringVariable(randomId),
					"notification_groups": config.ListVariable(config.StringVariable(groupHrid)),
					"events": config.ListVariable(
						config.StringVariable("API_STARTED"),
						config.StringVariable("API_STOPPED"),
					),
				},
			},
			// Update: swap one event and change order.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":                config.StringVariable(randomId),
					"notification_groups": config.ListVariable(config.StringVariable(groupHrid)),
					"events": config.ListVariable(
						config.StringVariable("API_DEPLOYED"),
						config.StringVariable("API_STARTED"),
					),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
			},
			// Update: remove the notification group.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":                config.StringVariable(randomId),
					"notification_groups": config.ListVariable(),
					"events": config.ListVariable(
						config.StringVariable("API_DEPLOYED"),
						config.StringVariable("API_STARTED"),
					),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
			},
			// Update: turn console notifications off.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":                 config.StringVariable(randomId),
					"notification_groups":  config.ListVariable(),
					"notification_enabled": config.BoolVariable(false),
					"events": config.ListVariable(
						config.StringVariable("API_DEPLOYED"),
						config.StringVariable("API_STARTED"),
					),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"apim_apiv4.test",
						tfjsonpath.New("console_notification"),
						knownvalue.Null(),
					),
				},
			},
			// An empty events list is rejected at plan time.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":                config.StringVariable(randomId),
					"notification_groups": config.ListVariable(),
					// ListVariable() without arguments renders null, not [].
					"events": config.ListVariable([]config.Variable{}...),
				},
				ExpectError: regexp.MustCompile("at least 1"),
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

func TestAPIV4Resource_webhook(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)

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
			// Testing framework implicitly verifies resource delete.
		},
	})
}

func TestAPIV4Resource_kafka_native(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)

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
			// Testing framework implicitly verifies resource delete.
		},
	})
}

func TestAPIV4Resource_allow_fields(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)

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
			// Testing framework implicitly verifies resource delete.
		},
	})
}
