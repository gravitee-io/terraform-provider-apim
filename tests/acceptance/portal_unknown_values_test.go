package acceptance_test

import (
	"fmt"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestPortalResource_unknownValuesAtPlan(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	hrid := "probe-" + acctest.RandString(10)
	path := "/guides-" + hrid
	config := fmt.Sprintf(`
provider "apim" {
  organization_id = "DEFAULT"
  environment_id = "DEFAULT"
}
resource "terraform_data" "navigation" {
  input = { path = %q, visibility = "PRIVATE" }
}
resource "apim_portal" "test" {
  hrid = %q
  name = "Unknown input probe"
  structure = { top_navbar = [{ path = terraform_data.navigation.output.path, visibility = terraform_data.navigation.output.visibility }] }
}
`, path, hrid)
	resource.Test(t, resource.TestCase{Steps: []resource.TestStep{{
		ProtoV6ProviderFactories: testProviders(), Config: config,
		ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
			plancheck.ExpectUnknownValue("apim_portal.test", tfjsonpath.New("structure").AtMapKey("top_navbar").AtSliceIndex(0).AtMapKey("path")),
			plancheck.ExpectUnknownValue("apim_portal.test", tfjsonpath.New("structure").AtMapKey("top_navbar").AtSliceIndex(0).AtMapKey("visibility")),
		}},
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("apim_portal.test", "structure.top_navbar.0.path", path),
			resource.TestCheckResourceAttr("apim_portal.test", "structure.top_navbar.0.visibility", "PRIVATE"),
		),
	}}})
}
