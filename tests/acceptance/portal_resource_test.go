package acceptance_test

import (
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
