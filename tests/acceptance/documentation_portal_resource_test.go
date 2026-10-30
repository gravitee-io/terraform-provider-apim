package acceptance_test

import (
	"fmt"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDocumentationPortalResource_nextGen(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()
	hrid := "test-" + acctest.RandString(10)
	address := "apim_documentation_portal.test"
	pageConfig := func(area string) string {
		placement := fmt.Sprintf("location = %q\norder = 1", "/docs-"+hrid)
		if area == "HOMEPAGE" {
			placement = `area = "HOMEPAGE"`
		}
		return fmt.Sprintf(`
provider "apim" {
  organization_id = "DEFAULT"
  environment_id = "DEFAULT"
}
resource "apim_portal" "parent" {
  hrid = %q
  name = "Documentation portal"
  structure = { top_navbar = [{ path = %q, visibility = "PRIVATE" }] }
}
resource "apim_documentation_portal" "test" {
  portal_hrid = apim_portal.parent.hrid
  hrid = %q
  name = "Next-gen documentation"
  type = "GRAVITEE_MARKDOWN"
  content = "# Portal documentation"
  %s
}
`, "portal-"+hrid, "/docs-"+hrid, hrid, placement)
	}
	check := func(area string) resource.TestCheckFunc {
		checks := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(address, "area", area),
			resource.TestCheckResourceAttr(address, "hrid", hrid),
			resource.TestCheckResourceAttr(address, "name", "Next-gen documentation"),
			resource.TestCheckResourceAttr(address, "type", "GRAVITEE_MARKDOWN"),
			resource.TestCheckResourceAttr(address, "content", "# Portal documentation"),
		}
		if area == "TOP_NAVBAR" {
			checks = append(checks,
				resource.TestCheckResourceAttr(address, "visibility", "PRIVATE"),
				resource.TestCheckResourceAttr(address, "location", "/docs-"+hrid),
				resource.TestCheckResourceAttr(address, "order", "1"),
			)
		}
		return resource.ComposeTestCheckFunc(checks...)
	}
	resource.Test(t, resource.TestCase{Steps: []resource.TestStep{
		{ProtoV6ProviderFactories: testProviders(), Config: pageConfig("TOP_NAVBAR"), Check: check("TOP_NAVBAR")},
		{ProtoV6ProviderFactories: testProviders(), Config: pageConfig("HOMEPAGE"), Check: check("HOMEPAGE")},
		{
			ProtoV6ProviderFactories: testProviders(), Config: pageConfig("HOMEPAGE"),
			ResourceName: address, ImportState: true, ImportStateVerify: true,
			ImportStateVerifyIdentifierAttribute: "hrid",
			ImportStateIdFunc:                    importStateIDFunc(address, []string{"environment_id", "hrid", "organization_id", "portal_hrid"}, nil),
		},
		{ProtoV6ProviderFactories: testProviders(), Config: pageConfig("HOMEPAGE"), PlanOnly: true},
		{ProtoV6ProviderFactories: testProviders(), Config: pageConfig("TOP_NAVBAR"), Check: check("TOP_NAVBAR")},
	}})
}

func TestDocumentationPortalResource_minimal(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
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
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
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
