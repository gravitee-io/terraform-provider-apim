package acceptance_test

import (
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Releases whose Automation API has no defaultMemberRoles, ignoreMembers or ignoreGroups.
var withoutGroupAutomationM2 = []utils.ApimVersion{utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12}

// Verifies that default member roles are created, read back, imported, and converge when a scope is removed or all are cleared.
func TestGroupResource_defaultMemberRoles(t *testing.T) {
	utils.SkipFor(t, withoutGroupAutomationM2...)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_group.test"
	allScopes := config.ObjectVariable(map[string]config.Variable{
		"api":         config.StringVariable("USER"),
		"application": config.StringVariable("USER"),
		"api_product": config.StringVariable("USER"),
	})
	apiOnly := config.ObjectVariable(map[string]config.Variable{"api": config.StringVariable("OWNER")})
	none := config.ObjectVariable(map[string]config.Variable{})

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          config.Variables{"hrid": config.StringVariable(randomId), "default_member_roles": allScopes},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "default_member_roles.application", "USER"),
					resource.TestCheckResourceAttr(resourceAddress, "default_member_roles.api_product", "USER"),
				),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          config.Variables{"hrid": config.StringVariable(randomId), "default_member_roles": allScopes},
				ResourceName:             resourceAddress,
				ImportState:              true,
				ImportStateIdFunc:        importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id"}, nil),
				ImportStateVerify:        true,
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          config.Variables{"hrid": config.StringVariable(randomId), "default_member_roles": apiOnly},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "default_member_roles.api", "OWNER"),
					resource.TestCheckNoResourceAttr(resourceAddress, "default_member_roles.application"),
					resource.TestCheckNoResourceAttr(resourceAddress, "default_member_roles.api_product"),
				),
			},
			{
				// {} clears all three; APIM answers {} back, so the follow-up plan is empty
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          config.Variables{"hrid": config.StringVariable(randomId), "default_member_roles": none},
				Check:                    resource.TestCheckNoResourceAttr(resourceAddress, "default_member_roles.api"),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          config.Variables{"hrid": config.StringVariable(randomId), "default_member_roles": none},
				ResourceName:             resourceAddress,
				ImportState:              true,
				ImportStateIdFunc:        importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id"}, nil),
				ImportStateVerify:        true,
			},
		},
	})
}

// Verifies that ignore_members keeps the members the group already has when the configuration declares none.
func TestGroupResource_ignoreMembers(t *testing.T) {
	utils.SkipFor(t, withoutGroupAutomationM2...)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_group.test"
	members := config.ListVariable(config.ObjectVariable(config.Variables{
		"source":    config.StringVariable("memory"),
		"source_id": config.StringVariable("api1"),
		"roles": config.MapVariable(map[string]config.Variable{
			"API":         config.StringVariable("USER"),
			"APPLICATION": config.StringVariable("USER"),
			"INTEGRATION": config.StringVariable("USER"),
		}),
	}))

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId), "ignore_members": config.BoolVariable(false), "members": members,
				},
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          config.Variables{"hrid": config.StringVariable(randomId), "ignore_members": config.BoolVariable(true)},
				Check:                    resource.TestCheckResourceAttr(resourceAddress, "members.#", "1"),
			},
		},
	})
}
