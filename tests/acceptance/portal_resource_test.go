package acceptance_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPortalResource_minimal(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_portal.test"

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
				},
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
				},
				ResourceName:      resourceAddress,
				ImportState:       true,
				ImportStateIdFunc: importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id"}, nil),
				ImportStateVerify: true,
			},
		},
	})
}

func portalNavigationConfig(hrid, name, navigation string) string {
	return fmt.Sprintf(`
provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}
resource "apim_portal" "test" {
  hrid = %q
  name = %q
  %s
}
`, hrid, name, navigation)
}

func TestPortalResource_navigationLifecycle(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()
	hrid := "test-" + acctest.RandString(10)
	address := "apim_portal.test"
	newPaths := fmt.Sprintf(`[{ path = %q, display_name = "Guides" }]`, "/guides-"+hrid)
	legacyPaths := fmt.Sprintf(`[{ path = %q, display_name = "Reference" }]`, "/reference-"+hrid)
	newConfig := "structure = { top_navbar = " + newPaths + " }"
	legacyConfig := "navigation = " + legacyPaths
	checkPaths := func(count, firstPath string) resource.TestCheckFunc {
		checks := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(address, "navigation.#", count),
			resource.TestCheckResourceAttr(address, "structure.top_navbar.#", count),
		}
		if firstPath != "" {
			checks = append(checks,
				resource.TestCheckResourceAttr(address, "navigation.0.path", firstPath),
				resource.TestCheckResourceAttr(address, "structure.top_navbar.0.path", firstPath),
			)
		}
		return resource.ComposeTestCheckFunc(checks...)
	}
	step := func(name, navigation, count, firstPath string) resource.TestStep {
		return resource.TestStep{
			ProtoV6ProviderFactories: testProviders(),
			Config:                   portalNavigationConfig(hrid, name, navigation),
			Check:                    checkPaths(count, firstPath),
		}
	}
	retainedConfig := portalNavigationConfig(hrid, "Unmanaged navigation rename", "")
	importedConfig := portalNavigationConfig(hrid, "Imported unmanaged navigation rename", "")
	resource.Test(t, resource.TestCase{Steps: []resource.TestStep{
		step("New portal", newConfig, "1", "/guides-"+hrid),
		step("Renamed new portal", newConfig, "1", "/guides-"+hrid),
		step("Legacy migration", legacyConfig, "1", "/reference-"+hrid),
		step("Clear legacy", "navigation = []", "0", ""),
		step("Migrate back to structure", newConfig, "1", "/guides-"+hrid),
		step("Clear structure", "structure = { top_navbar = [] }", "0", ""),
		step("Restore structure", newConfig, "1", "/guides-"+hrid),
		step("Null navigation inputs rename", "navigation = null\nstructure = null", "1", "/guides-"+hrid),
		step("Unmanaged navigation rename", "", "1", "/guides-"+hrid),
		{ProtoV6ProviderFactories: testProviders(), Config: retainedConfig, PlanOnly: true},
		{
			ProtoV6ProviderFactories: testProviders(), Config: retainedConfig,
			ResourceName: address, ImportState: true, ImportStateVerify: true,
			ImportStatePersist: true,
			ImportStateIdFunc:  importStateIDFunc(address, []string{"environment_id", "hrid", "organization_id"}, nil),
		},
		step("Imported unmanaged navigation rename", "", "1", "/guides-"+hrid),
		{ProtoV6ProviderFactories: testProviders(), Config: importedConfig, PlanOnly: true},
	}})
}

func TestPortalResource_emptyStructure(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()
	hrid := "test-" + acctest.RandString(10)
	validConfig := portalNavigationConfig(hrid, "Preserved navigation", `structure = { top_navbar = [{ path = "/guides" }] }`)
	resource.Test(t, resource.TestCase{Steps: []resource.TestStep{
		{
			ProtoV6ProviderFactories: testProviders(),
			Config:                   portalNavigationConfig(hrid, "Empty structure", "structure = {}"),
			ExpectError:              regexp.MustCompile("Missing portal top navbar"),
		},
		{
			ProtoV6ProviderFactories: testProviders(),
			Config:                   validConfig,
		},
		{
			ProtoV6ProviderFactories: testProviders(),
			Config:                   portalNavigationConfig(hrid, "Empty structure update", "structure = {}"),
			ExpectError:              regexp.MustCompile("Missing portal top navbar"),
		},
		{
			ProtoV6ProviderFactories: testProviders(),
			Config:                   portalNavigationConfig(hrid, "Null navbar", "structure = { top_navbar = null }"),
			ExpectError:              regexp.MustCompile("Missing portal top navbar"),
		},
		{
			ProtoV6ProviderFactories: testProviders(),
			Config:                   validConfig,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("apim_portal.test", "structure.top_navbar.#", "1"),
				resource.TestCheckResourceAttr("apim_portal.test", "structure.top_navbar.0.path", "/guides"),
				resource.TestCheckResourceAttr("apim_portal.test", "navigation.0.path", "/guides"),
			),
		},
		{ProtoV6ProviderFactories: testProviders(), Config: validConfig, PlanOnly: true},
	}})
}

func TestPortalResource_navigationConflict(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()
	resource.Test(t, resource.TestCase{Steps: []resource.TestStep{{
		ProtoV6ProviderFactories: testProviders(),
		Config:                   portalNavigationConfig("test-"+acctest.RandString(10), "Conflict", "navigation = []\nstructure = { top_navbar = [] }"),
		ExpectError:              regexp.MustCompile("navigation and structure cannot be configured together"),
	}}})
}

func TestPortalResource_update(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)

	navigationInitial := config.ListVariable(
		config.ObjectVariable(config.Variables{
			"path":         config.StringVariable("/guides" + randomId),
			"display_name": config.StringVariable("Guides"),
		}),
		config.ObjectVariable(config.Variables{
			"path":         config.StringVariable("/guides" + randomId + "/getting-started" + randomId),
			"display_name": config.StringVariable("Getting Started"),
		}),
	)
	navigationUpdated := config.ListVariable(
		config.ObjectVariable(config.Variables{
			"path":         config.StringVariable("/reference" + randomId),
			"display_name": config.StringVariable("Reference"),
		}),
		config.ObjectVariable(config.Variables{
			"path":         config.StringVariable("/guides" + randomId),
			"display_name": config.StringVariable("Guides"),
		}),
		config.ObjectVariable(config.Variables{
			"path":         config.StringVariable("/guides" + randomId + "/getting-started" + randomId),
			"display_name": config.StringVariable("Getting Started"),
		}),
	)
	navigationReduced := config.ListVariable(
		config.ObjectVariable(config.Variables{
			"path":         config.StringVariable("/guides" + randomId),
			"display_name": config.StringVariable("Guides"),
		}),
		config.ObjectVariable(config.Variables{
			"path":         config.StringVariable("/reference" + randomId),
			"display_name": config.StringVariable("Reference"),
		}),
	)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":       config.StringVariable(randomId),
					"name":       config.StringVariable("Portal One"),
					"navigation": navigationInitial,
				},
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":       config.StringVariable(randomId),
					"name":       config.StringVariable("Portal Updated"),
					"navigation": navigationUpdated,
				},
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":       config.StringVariable(randomId),
					"name":       config.StringVariable("Portal Updated (order idempotent)"),
					"navigation": navigationUpdated,
				},
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":       config.StringVariable(randomId),
					"name":       config.StringVariable("Portal Updated"),
					"navigation": navigationReduced,
				},
			},
		},
	})
}
