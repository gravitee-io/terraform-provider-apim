package acceptance_test

import (
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPortalListingResource_minimal(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()

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
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	portalHrid := "portal-" + randomId

	apisLocation := "/apis-" + portalHrid
	apisFeaturedLocation := apisLocation + "/featured-" + portalHrid

	apisInitial := config.ListVariable(
		config.ObjectVariable(config.Variables{
			"api":      config.StringVariable("first"),
			"location": config.StringVariable(apisLocation),
			"order":    config.IntegerVariable(1),
		}),
		config.ObjectVariable(config.Variables{
			"api":      config.StringVariable("second"),
			"location": config.StringVariable(apisLocation),
			"order":    config.IntegerVariable(2),
		}),
	)
	apisUpdated := config.ListVariable(
		config.ObjectVariable(config.Variables{
			"api":      config.StringVariable("second"),
			"location": config.StringVariable(apisLocation),
			"order":    config.IntegerVariable(1),
		}),
		config.ObjectVariable(config.Variables{
			"api":      config.StringVariable("first"),
			"location": config.StringVariable(apisFeaturedLocation),
			"order":    config.IntegerVariable(2),
		}),
	)
	apisReduced := config.ListVariable(
		config.ObjectVariable(config.Variables{
			"api":      config.StringVariable("first"),
			"location": config.StringVariable(apisFeaturedLocation),
			"order":    config.IntegerVariable(1),
		}),
	)

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
			},
		},
	})
}
