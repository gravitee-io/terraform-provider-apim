package acceptance_test

import (
	"fmt"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const portalThemeTestImage = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+jT1sAAAAASUVORK5CYII="

func portalThemeConfig(hrid, name, primary string, full, attached bool) string {
	definition := "definition = {}"
	images := ""
	if full {
		definition = fmt.Sprintf(`definition = {
      custom_css = "body { line-height: 1.5; }"
      font = { font_family = "Arial, sans-serif" }
      color = {
        primary = %q
        secondary = "#445566"
        tertiary = "#778899"
        error = "#bb1122"
        background = { page = "#ffffff", card = "#eeeeee" }
      }
    }`, primary)
		images = fmt.Sprintf("logo = %[1]q\noptional_logo = %[1]q\nfavicon = %[1]q\nbackground_image = %[1]q", portalThemeTestImage)
	}
	parent := ""
	if attached {
		parent = fmt.Sprintf(`
resource "apim_portal" "attached" {
  hrid = %q
  name = "Themed portal"
  active_theme_hrid = apim_portal_theme.test.hrid
}
`, "portal-"+hrid)
	}
	return fmt.Sprintf(`
provider "apim" {
  organization_id = "DEFAULT"
  environment_id = "DEFAULT"
}
resource "apim_portal_theme" "test" {
  hrid = %q
  name = %q
  %s
  %s
}
%s
`, hrid, name, definition, images, parent)
}

// Applying a theme changes the environment's active theme, so these tests run serially.
func TestPortalThemeResource_minimal(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	hrid := "test-" + acctest.RandString(10)
	address := "apim_portal_theme.test"
	configuration := portalThemeConfig(hrid, "Minimal theme", "", false, true)
	resource.Test(t, resource.TestCase{Steps: []resource.TestStep{
		{
			ProtoV6ProviderFactories: testProviders(), Config: configuration,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(address, "hrid", hrid),
				resource.TestCheckResourceAttr(address, "name", "Minimal theme"),
				resource.TestCheckNoResourceAttr(address, "definition.color.primary"),
				resource.TestCheckResourceAttrPair("apim_portal.attached", "active_theme_hrid", address, "hrid"),
			),
		},
		{
			ProtoV6ProviderFactories: testProviders(), Config: configuration,
			ResourceName: address, ImportState: true, ImportStateVerify: true,
			ImportStateVerifyIdentifierAttribute: "hrid",
			ImportStateIdFunc:                    importStateIDFunc(address, []string{"environment_id", "hrid", "organization_id"}, nil),
		},
		{ProtoV6ProviderFactories: testProviders(), Config: configuration, PlanOnly: true},
	}})
}

func TestPortalThemeResource_allProps(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	hrid := "test-" + acctest.RandString(10)
	address := "apim_portal_theme.test"
	initial := portalThemeConfig(hrid, "Brand theme", "#112233", true, true)
	updated := portalThemeConfig(hrid, "Updated brand theme", "#334455", true, true)
	detached := portalThemeConfig(hrid, "Updated brand theme", "#334455", true, false)
	check := func(name, primary string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(address, "hrid", hrid),
			resource.TestCheckResourceAttr(address, "name", name),
			resource.TestCheckResourceAttr(address, "definition.custom_css", "body { line-height: 1.5; }"),
			resource.TestCheckResourceAttr(address, "definition.font.font_family", "Arial, sans-serif"),
			resource.TestCheckResourceAttr(address, "definition.color.primary", primary),
			resource.TestCheckResourceAttr(address, "definition.color.secondary", "#445566"),
			resource.TestCheckResourceAttr(address, "definition.color.tertiary", "#778899"),
			resource.TestCheckResourceAttr(address, "definition.color.error", "#bb1122"),
			resource.TestCheckResourceAttr(address, "definition.color.background.page", "#ffffff"),
			resource.TestCheckResourceAttr(address, "definition.color.background.card", "#eeeeee"),
			resource.TestCheckResourceAttr(address, "logo", portalThemeTestImage),
			resource.TestCheckResourceAttr(address, "optional_logo", portalThemeTestImage),
			resource.TestCheckResourceAttr(address, "favicon", portalThemeTestImage),
			resource.TestCheckResourceAttr(address, "background_image", portalThemeTestImage),
		)
	}
	resource.Test(t, resource.TestCase{Steps: []resource.TestStep{
		{ProtoV6ProviderFactories: testProviders(), Config: initial, Check: check("Brand theme", "#112233")},
		{ProtoV6ProviderFactories: testProviders(), Config: updated, Check: check("Updated brand theme", "#334455")},
		{ProtoV6ProviderFactories: testProviders(), Config: updated, PlanOnly: true},
		{
			ProtoV6ProviderFactories: testProviders(), Config: updated,
			ResourceName: address, ImportState: true, ImportStateVerify: true,
			ImportStateVerifyIdentifierAttribute: "hrid",
			ImportStateIdFunc:                    importStateIDFunc(address, []string{"environment_id", "hrid", "organization_id"}, nil),
		},
		{ProtoV6ProviderFactories: testProviders(), Config: detached, Check: check("Updated brand theme", "#334455")},
		{ProtoV6ProviderFactories: testProviders(), Config: detached, PlanOnly: true},
	}})
}
