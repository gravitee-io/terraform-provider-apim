package acceptance_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

func portalStructureConfig(hrid, name, structure string) string {
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
`, hrid, name, structure)
}

func TestPortalResource_structureLifecycle(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()
	hrid := "test-" + acctest.RandString(10)
	address := "apim_portal.test"
	newPaths := fmt.Sprintf(`[{ path = %q, display_name = "Guides" }]`, "/guides-"+hrid)
	updatedPaths := fmt.Sprintf(`[{ path = %q, display_name = "Reference" }, { path = %q, display_name = "Guides" }]`, "/reference-"+hrid, "/guides-"+hrid)
	newConfig := "structure = { top_navbar = " + newPaths + " }"
	updatedConfig := "structure = { top_navbar = " + updatedPaths + " }"
	checkPaths := func(count, firstPath string) resource.TestCheckFunc {
		checks := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(address, "structure.top_navbar.#", count),
		}
		if firstPath != "" {
			checks = append(checks,
				resource.TestCheckResourceAttr(address, "structure.top_navbar.0.path", firstPath),
			)
		}
		return resource.ComposeTestCheckFunc(checks...)
	}
	step := func(name, structure, count, firstPath string) resource.TestStep {
		return resource.TestStep{
			ProtoV6ProviderFactories: testProviders(),
			Config:                   portalStructureConfig(hrid, name, structure),
			Check:                    checkPaths(count, firstPath),
		}
	}
	retainedConfig := portalStructureConfig(hrid, "Unmanaged navigation rename", "")
	importedConfig := portalStructureConfig(hrid, "Imported unmanaged navigation rename", "")
	forgetConfig := `provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}
removed {
  from = apim_portal.test

  lifecycle {
    destroy = false
  }
}
`
	resource.Test(t, resource.TestCase{Steps: []resource.TestStep{
		step("New portal", newConfig, "1", "/guides-"+hrid),
		step("Renamed new portal", newConfig, "1", "/guides-"+hrid),
		step("Update navigation order", updatedConfig, "2", "/reference-"+hrid),
		step("Repeat navigation order", updatedConfig, "2", "/reference-"+hrid),
		step("Clear structure", "structure = { top_navbar = [] }", "0", ""),
		step("Restore structure", newConfig, "1", "/guides-"+hrid),
		step("Null structure rename", "structure = null", "1", "/guides-"+hrid),
		step("Unmanaged navigation rename", "", "1", "/guides-"+hrid),
		{ProtoV6ProviderFactories: testProviders(), Config: retainedConfig, PlanOnly: true},
		{ProtoV6ProviderFactories: testProviders(), Config: forgetConfig},
		{
			ProtoV6ProviderFactories: testProviders(), Config: importedConfig,
			ResourceName: address, ImportState: true,
			ImportStatePersist: true,
			ImportStateId:      fmt.Sprintf(`{"environment_id":"DEFAULT","hrid":%q,"organization_id":"DEFAULT"}`, hrid),
			ImportStateCheck: func(states []*terraform.InstanceState) error {
				if len(states) != 1 {
					return fmt.Errorf("expected one imported portal state, got %d", len(states))
				}
				attrs := states[0].Attributes
				for key, want := range map[string]string{
					"hrid":                        hrid,
					"environment_id":              "DEFAULT",
					"organization_id":             "DEFAULT",
					"structure.top_navbar.#":      "1",
					"structure.top_navbar.0.path": "/guides-" + hrid,
				} {
					if got := attrs[key]; got != want {
						return fmt.Errorf("imported portal attribute %s: expected %q, got %q", key, want, got)
					}
				}
				return nil
			},
		},
		step("Imported unmanaged navigation rename", "", "1", "/guides-"+hrid),
		{ProtoV6ProviderFactories: testProviders(), Config: importedConfig, PlanOnly: true},
	}})
}

func TestPortalResource_emptyStructure(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()
	hrid := "test-" + acctest.RandString(10)
	path := "/guides-" + hrid
	validConfig := portalStructureConfig(hrid, "Preserved navigation", fmt.Sprintf(`structure = { top_navbar = [{ path = %q }] }`, path))
	resource.Test(t, resource.TestCase{Steps: []resource.TestStep{
		{
			ProtoV6ProviderFactories: testProviders(),
			Config:                   portalStructureConfig(hrid, "Empty structure", "structure = {}"),
			ExpectError:              regexp.MustCompile("Missing portal top navbar"),
		},
		{
			ProtoV6ProviderFactories: testProviders(),
			Config:                   validConfig,
		},
		{
			ProtoV6ProviderFactories: testProviders(),
			Config:                   portalStructureConfig(hrid, "Empty structure update", "structure = {}"),
			ExpectError:              regexp.MustCompile("Missing portal top navbar"),
		},
		{
			ProtoV6ProviderFactories: testProviders(),
			Config:                   portalStructureConfig(hrid, "Null navbar", "structure = { top_navbar = null }"),
			ExpectError:              regexp.MustCompile("Missing portal top navbar"),
		},
		{
			ProtoV6ProviderFactories: testProviders(),
			Config:                   validConfig,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("apim_portal.test", "structure.top_navbar.#", "1"),
				resource.TestCheckResourceAttr("apim_portal.test", "structure.top_navbar.0.path", path),
			),
		},
		{ProtoV6ProviderFactories: testProviders(), Config: validConfig, PlanOnly: true},
	}})
}

func TestPortalResource_update(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	guidesPath := "/guides-" + randomId

	entry := func(path, name string) config.Variable {
		return config.ObjectVariable(config.Variables{
			"path":         config.StringVariable(path),
			"display_name": config.StringVariable(name),
		})
	}
	structure := func(entries ...config.Variable) config.Variable {
		return config.ObjectVariable(config.Variables{"top_navbar": config.ListVariable(entries...)})
	}
	structureInitial := structure(
		entry(guidesPath, "Guides"),
		entry(guidesPath+"/getting-started", "Getting Started"),
	)
	structureUpdated := structure(
		entry("/reference-"+randomId, "Reference"),
		entry(guidesPath, "Guides"),
		entry(guidesPath+"/getting-started", "Getting Started"),
	)
	structureReduced := structure(
		entry(guidesPath, "Guides"),
		entry("/reference-"+randomId, "Reference"),
	)
	checkStructure := func(name string, entries ...[2]string) resource.TestCheckFunc {
		checks := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr("apim_portal.test", "name", name),
			resource.TestCheckResourceAttr("apim_portal.test", "hrid", randomId),
			resource.TestCheckResourceAttr("apim_portal.test", "structure.top_navbar.#", fmt.Sprint(len(entries))),
		}
		for index, entry := range entries {
			checks = append(checks,
				resource.TestCheckResourceAttr("apim_portal.test", fmt.Sprintf("structure.top_navbar.%d.path", index), entry[0]),
				resource.TestCheckResourceAttr("apim_portal.test", fmt.Sprintf("structure.top_navbar.%d.display_name", index), entry[1]),
			)
		}
		return resource.ComposeTestCheckFunc(checks...)
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":      config.StringVariable(randomId),
					"name":      config.StringVariable("Portal One"),
					"structure": structureInitial,
				},
				Check: checkStructure("Portal One", [2]string{guidesPath, "Guides"}, [2]string{guidesPath + "/getting-started", "Getting Started"}),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":      config.StringVariable(randomId),
					"name":      config.StringVariable("Portal Updated"),
					"structure": structureUpdated,
				},
				Check: checkStructure("Portal Updated", [2]string{"/reference-" + randomId, "Reference"}, [2]string{guidesPath, "Guides"}, [2]string{guidesPath + "/getting-started", "Getting Started"}),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":      config.StringVariable(randomId),
					"name":      config.StringVariable("Portal Updated (order idempotent)"),
					"structure": structureUpdated,
				},
				Check: checkStructure("Portal Updated (order idempotent)", [2]string{"/reference-" + randomId, "Reference"}, [2]string{guidesPath, "Guides"}, [2]string{guidesPath + "/getting-started", "Getting Started"}),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":      config.StringVariable(randomId),
					"name":      config.StringVariable("Portal Updated"),
					"structure": structureReduced,
				},
				Check: checkStructure("Portal Updated", [2]string{guidesPath, "Guides"}, [2]string{"/reference-" + randomId, "Reference"}),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":      config.StringVariable(randomId),
					"name":      config.StringVariable("Portal Updated"),
					"structure": structureReduced,
				},
				PlanOnly: true,
				Check:    checkStructure("Portal Updated", [2]string{guidesPath, "Guides"}, [2]string{"/reference-" + randomId, "Reference"}),
			},
		},
	})
}
