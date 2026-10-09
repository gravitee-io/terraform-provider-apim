package acceptance_test

import (
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDocumentationPortalResource_minimal(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	portalHrid := "api-" + randomId
	resourceAddress := "apim_documentation_portal.test"

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":        config.StringVariable(randomId),
					"portal_hrid": config.StringVariable(portalHrid),
				},
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":        config.StringVariable(randomId),
					"portal_hrid": config.StringVariable(portalHrid),
				},
				ResourceName:                         resourceAddress,
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id", "portal_hrid"}, nil),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "hrid",
			},
		},
	})
}

func TestDocumentationPortalResource_update(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	portalHrid := "portal-" + randomId

	contentInitial := "# Getting Started\n\nThis page has multiple lines.\n- first\n- second\n"
	contentUpdated := "# Reference\n\nThis page also has multiple lines.\n\n```yaml\nopenapi: 3.0.0\n```\n"

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"content":     config.StringVariable(contentInitial),
					"hrid":        config.StringVariable(randomId),
					"name":        config.StringVariable("Getting Started"),
					"order":       config.IntegerVariable(1),
					"portal_hrid": config.StringVariable(portalHrid),
				},
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"content":     config.StringVariable(contentUpdated),
					"hrid":        config.StringVariable(randomId),
					"name":        config.StringVariable("Reference"),
					"order":       config.IntegerVariable(2),
					"portal_hrid": config.StringVariable(portalHrid),
				},
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"content":     config.StringVariable(contentUpdated),
					"hrid":        config.StringVariable(randomId),
					"name":        config.StringVariable("Reference"),
					"order":       config.IntegerVariable(2),
					"portal_hrid": config.StringVariable(portalHrid),
				},
			},
		},
	})
}
