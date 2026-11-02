package acceptance_test

import (
	"fmt"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func apiLinkConfig(hrid, name, href, extra string) string {
	return fmt.Sprintf(`
provider "apim" {
  organization_id = "DEFAULT"
  environment_id = "DEFAULT"
}
resource "apim_apiv4" "parent" {
  hrid = "api-%[1]s"
  name = "API link parent"
  type = "PROXY"
  version = "1"
  state = "STOPPED"
  lifecycle_state = "UNPUBLISHED"
  visibility = "PRIVATE"
  portal_navigation = [
    { path = "/private-%[1]s", visibility = "PRIVATE" },
    { path = "/public-%[1]s", visibility = "PUBLIC" }
  ]
  endpoint_groups = [{
    name = "Default HTTP proxy group"
    type = "http-proxy"
    services = {}
    endpoints = [{
      name = "Default HTTP proxy"
      type = "http-proxy"
      configuration = jsonencode({ target = "https://example.com" })
      services = {}
    }]
  }]
  listeners = [{ http = {
    type = "HTTP"
    entrypoints = [{ type = "http-proxy" }]
    paths = [{ path = "/api-link-%[1]s/" }]
  } }]
}

resource "apim_api_link" "test" {
  api_hrid = apim_apiv4.parent.hrid
  hrid = "%[1]s"
  name = %[2]q
  href = %[3]q
  %[4]s
}
`, hrid, name, href, extra)
}

func TestAPILinkResource_minimal(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()
	hrid := "test-" + acctest.RandString(10)
	address := "apim_api_link.test"
	configuration := apiLinkConfig(hrid, "External docs", "https://example.com/docs", "")
	resource.Test(t, resource.TestCase{Steps: []resource.TestStep{
		{
			ProtoV6ProviderFactories: testProviders(), Config: configuration,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(address, "hrid", hrid),
				resource.TestCheckResourceAttr(address, "name", "External docs"),
				resource.TestCheckResourceAttr(address, "href", "https://example.com/docs"),
				resource.TestCheckResourceAttr(address, "api_hrid", "api-"+hrid),
				resource.TestCheckResourceAttr(address, "visibility", "PUBLIC"),
			),
		},
		{
			ProtoV6ProviderFactories: testProviders(), Config: configuration,
			ResourceName: address, ImportState: true, ImportStateVerify: true,
			ImportStateVerifyIdentifierAttribute: "hrid",
			ImportStateIdFunc:                    importStateIDFunc(address, []string{"environment_id", "hrid", "organization_id", "api_hrid"}, nil),
		},
		{ProtoV6ProviderFactories: testProviders(), Config: configuration, PlanOnly: true},
	}})
}

func TestAPILinkResource_allProps(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()
	hrid := "test-" + acctest.RandString(10)
	address := "apim_api_link.test"
	initial := apiLinkConfig(hrid, "Private docs", "https://example.com/private", fmt.Sprintf("location = %q\norder = 1", "/private-"+hrid))
	updated := apiLinkConfig(hrid, "Public reference", "https://example.com/reference", fmt.Sprintf("location = %q\norder = 2\nvisibility = %q", "/public-"+hrid, "PUBLIC"))
	check := func(name, href, location, order, visibility string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(address, "hrid", hrid),
			resource.TestCheckResourceAttr(address, "api_hrid", "api-"+hrid),
			resource.TestCheckResourceAttr(address, "name", name),
			resource.TestCheckResourceAttr(address, "href", href),
			resource.TestCheckResourceAttr(address, "location", location),
			resource.TestCheckResourceAttr(address, "order", order),
			resource.TestCheckResourceAttr(address, "visibility", visibility),
		)
	}
	resource.Test(t, resource.TestCase{Steps: []resource.TestStep{
		{ProtoV6ProviderFactories: testProviders(), Config: initial, Check: check("Private docs", "https://example.com/private", "/private-"+hrid, "1", "PRIVATE")},
		{ProtoV6ProviderFactories: testProviders(), Config: updated, Check: check("Public reference", "https://example.com/reference", "/public-"+hrid, "2", "PUBLIC")},
		{ProtoV6ProviderFactories: testProviders(), Config: updated, PlanOnly: true},
		{
			ProtoV6ProviderFactories: testProviders(), Config: updated,
			ResourceName: address, ImportState: true, ImportStateVerify: true,
			ImportStateVerifyIdentifierAttribute: "hrid",
			ImportStateIdFunc:                    importStateIDFunc(address, []string{"environment_id", "hrid", "organization_id", "api_hrid"}, nil),
		},
		{ProtoV6ProviderFactories: testProviders(), Config: initial, Check: check("Private docs", "https://example.com/private", "/private-"+hrid, "1", "PRIVATE")},
	}})
}
