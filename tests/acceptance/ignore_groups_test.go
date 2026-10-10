package acceptance_test

import (
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Releases whose Automation API has no ignoreGroups.
var withoutIgnoreGroups = []utils.ApimVersion{utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12}

// Verifies that ignore_groups keeps the API's groups when the configuration stops declaring them.
func TestAPIV4Resource_ignoreGroups(t *testing.T) {
	utils.SkipFor(t, withoutIgnoreGroups...)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	variables := func(ignoreGroups, declareGroup bool) config.Variables {
		return config.Variables{
			"environment_id":  config.StringVariable("DEFAULT"),
			"organization_id": config.StringVariable("DEFAULT"),
			"hrid":            config.StringVariable(randomId),
			"ignore_groups":   config.BoolVariable(ignoreGroups),
			"declare_group":   config.BoolVariable(declareGroup),
		}
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(false, true),
				Check:                    resource.TestCheckResourceAttr("apim_apiv4.test", "groups.#", "1"),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(true, false),
				Check:                    resource.TestCheckResourceAttr("apim_apiv4.test", "groups.#", "1"),
			},
		},
	})
}

// Verifies that ignore_groups keeps the application's groups when the configuration stops declaring them.
func TestApplicationResource_ignoreGroups(t *testing.T) {
	utils.SkipFor(t, withoutIgnoreGroups...)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	variables := func(ignoreGroups, declareGroup bool) config.Variables {
		return config.Variables{
			"hrid":          config.StringVariable(randomId),
			"ignore_groups": config.BoolVariable(ignoreGroups),
			"declare_group": config.BoolVariable(declareGroup),
		}
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(false, true),
				Check:                    resource.TestCheckResourceAttr("apim_application.test", "groups.#", "1"),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(true, false),
				Check:                    resource.TestCheckResourceAttr("apim_application.test", "groups.#", "1"),
			},
		},
	})
}
