package acceptance_test

import (
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestDocumentationAPIResource_minimal(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	apiHrid := "api-" + randomId
	resourceAddress := "apim_documentation_api.test"

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":     config.StringVariable(randomId),
					"api_hrid": config.StringVariable(apiHrid),
				},
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":     config.StringVariable(randomId),
					"api_hrid": config.StringVariable(apiHrid),
				},
				ResourceName:                         resourceAddress,
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIDFunc(resourceAddress, []string{"api_hrid", "environment_id", "hrid", "organization_id"}, nil),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "hrid",
			},
		},
	})
}

func TestDocumentationAPIResource_update(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_documentation_api.test"
	apiHrid := "api-" + randomId

	contentInitial := "# API Docs\n\nInitial content.\n"
	contentUpdated := "# API Reference\n\nUpdated content with **markdown**.\n"

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"api_hrid": config.StringVariable(apiHrid),
					"content":  config.StringVariable(contentInitial),
					"hrid":     config.StringVariable(randomId),
					"name":     config.StringVariable("API Docs"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "api_hrid", apiHrid),
					resource.TestCheckResourceAttr(resourceAddress, "content", contentInitial),
				),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"api_hrid": config.StringVariable(apiHrid),
					"content":  config.StringVariable(contentUpdated),
					"hrid":     config.StringVariable(randomId),
					"name":     config.StringVariable("API Reference"),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "name", "API Reference"),
					resource.TestCheckResourceAttr(resourceAddress, "content", contentUpdated),
				),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"api_hrid": config.StringVariable(apiHrid),
					"content":  config.StringVariable(contentUpdated),
					"hrid":     config.StringVariable(randomId),
					"name":     config.StringVariable("API Reference"),
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
