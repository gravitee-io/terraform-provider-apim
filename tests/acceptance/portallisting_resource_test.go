package acceptance_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPortalListingResource_minimal(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()
	lockPortalAcceptanceTest(t)

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_portal_listing.test"

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"organization_id": config.StringVariable(organizationId),
				},
			},
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
				ImportStateIdFunc: importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id", "portal_hrid"}, nil),
				ImportStateVerify: true,
			},
		},
	})
}

func TestPortalListingResource_update(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()
	lockPortalAcceptanceTest(t)

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	portalHrid := "portal-" + randomId
	resourceAddress := "apim_portal_listing.test"
	rootLocation := "/apis-" + portalHrid
	featuredLocation := rootLocation + "/featured-" + portalHrid

	apisInitial := config.ListVariable(
		config.ObjectVariable(config.Variables{
			"api":        config.StringVariable("first"),
			"location":   config.StringVariable("apis"),
			"order":      config.IntegerVariable(1),
			"visibility": config.StringVariable("PRIVATE"),
		}),
		config.ObjectVariable(config.Variables{
			"api":      config.StringVariable("second"),
			"location": config.StringVariable("apis"),
			"order":    config.IntegerVariable(2),
		}),
	)
	apisUpdated := config.ListVariable(
		config.ObjectVariable(config.Variables{
			"api":        config.StringVariable("second"),
			"location":   config.StringVariable("apis"),
			"order":      config.IntegerVariable(1),
			"visibility": config.StringVariable("PRIVATE"),
		}),
		config.ObjectVariable(config.Variables{
			"api":      config.StringVariable("first"),
			"location": config.StringVariable("featured"),
			"order":    config.IntegerVariable(2),
		}),
	)
	apisReduced := config.ListVariable(
		config.ObjectVariable(config.Variables{
			"api":      config.StringVariable("first"),
			"location": config.StringVariable("featured"),
			"order":    config.IntegerVariable(1),
		}),
	)

	type expectedEntry struct {
		apiHrid    string
		location   string
		order      int
		visibility string
	}
	checkEntries := func(entries ...expectedEntry) resource.TestCheckFunc {
		checks := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceAddress, "apis.#", strconv.Itoa(len(entries))),
			resource.TestCheckResourceAttr("apim_portal.primary", "structure.top_navbar.0.path", rootLocation),
			resource.TestCheckResourceAttr("apim_portal.primary", "structure.top_navbar.1.path", featuredLocation),
		}
		for i, entry := range entries {
			prefix := fmt.Sprintf("apis.%d.", i)
			checks = append(checks,
				resource.TestCheckResourceAttr(resourceAddress, prefix+"api_hrid", entry.apiHrid),
				resource.TestCheckResourceAttr(resourceAddress, prefix+"location", entry.location),
				resource.TestCheckResourceAttr(resourceAddress, prefix+"order", strconv.Itoa(entry.order)),
			)
			if entry.visibility == "" {
				checks = append(checks, resource.TestCheckNoResourceAttr(resourceAddress, prefix+"visibility"))
			} else {
				checks = append(checks, resource.TestCheckResourceAttr(resourceAddress, prefix+"visibility", entry.visibility))
			}
			name := "terraform_example_a-" + randomId
			if entry.apiHrid == "api-b-"+randomId {
				name = "terraform_example_b-" + randomId
			}
			effective := entry.visibility
			if effective == "" {
				effective = "PUBLIC"
				if entry.location == featuredLocation {
					effective = "PRIVATE"
				}
			}
			checks = append(checks, checkNavigationVisibility("apim_portal.primary", false, name, effective))
		}
		return resource.ComposeTestCheckFunc(checks...)
	}
	initialCheck := checkEntries(
		expectedEntry{"api-a-" + randomId, rootLocation, 1, "PRIVATE"},
		expectedEntry{"api-b-" + randomId, rootLocation, 2, ""},
	)
	updatedCheck := checkEntries(
		expectedEntry{"api-b-" + randomId, rootLocation, 1, "PRIVATE"},
		expectedEntry{"api-a-" + randomId, featuredLocation, 2, ""},
	)
	reducedCheck := checkEntries(expectedEntry{"api-a-" + randomId, featuredLocation, 1, ""})

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"organization_id": config.StringVariable(organizationId),
					"portal_hrid":     config.StringVariable(portalHrid),
					"apis":            apisInitial,
				},
				Check: initialCheck,
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"organization_id": config.StringVariable(organizationId),
					"portal_hrid":     config.StringVariable(portalHrid),
					"apis":            apisUpdated,
				},
				Check: updatedCheck,
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"organization_id": config.StringVariable(organizationId),
					"portal_hrid":     config.StringVariable(portalHrid),
					"apis":            apisReduced,
				},
				Check: reducedCheck,
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"organization_id": config.StringVariable(organizationId),
					"portal_hrid":     config.StringVariable(portalHrid),
					"apis":            apisUpdated,
				},
				Check: updatedCheck,
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"organization_id": config.StringVariable(organizationId),
					"portal_hrid":     config.StringVariable(portalHrid),
					"apis":            apisUpdated,
				},
				PlanOnly: true,
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"organization_id": config.StringVariable(organizationId),
					"portal_hrid":     config.StringVariable(portalHrid),
					"apis":            apisUpdated,
				},
				ResourceName:      resourceAddress,
				ImportState:       true,
				ImportStateIdFunc: importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id", "portal_hrid"}, nil),
				ImportStateVerify: true,
			},
		},
	})
}
