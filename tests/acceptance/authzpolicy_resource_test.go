package acceptance_test

import (
	"regexp"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAuthzPolicyResource_minimal(t *testing.T) {
	utils.SkipUnlessAuthzServes(t, utils.AuthzPolicies)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_authz_policy.test"
	variables := config.Variables{"hrid": config.StringVariable(randomId)}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceAddress, "id"),
					resource.TestCheckResourceAttr(resourceAddress, "policy_status", "DEPLOYED"),
					// No target: the default engine on every gateway.
					resource.TestCheckResourceAttr(resourceAddress, "targets.#", "0"),
				),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables,
				ConfigPlanChecks:         emptyPlan,
			},
			authzImportStep(resourceAddress, variables),
		},
	})
}

func TestAuthzPolicyResource_update(t *testing.T) {
	utils.SkipUnlessAuthzServes(t, utils.AuthzPolicies, utils.AuthzPdps)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_authz_policy.test"
	permit := `permit(principal in Group::"engineering", action == Action::"read", resource);`
	forbid := `forbid(principal in Group::"contractors", action == Action::"read", resource);`
	variables := func(name, policyText string, reversed bool) config.Variables {
		return config.Variables{
			"hrid":        config.StringVariable(randomId),
			"name":        config.StringVariable(name),
			"policy_text": config.StringVariable(policyText),
			"reversed":    config.BoolVariable(reversed),
		}
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables("Acceptance", permit, false),
				Check:                    resource.TestCheckResourceAttr(resourceAddress, "targets.#", "2"),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables("Acceptance, renamed", forbid, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "name", "Acceptance, renamed"),
					resource.TestCheckResourceAttr(resourceAddress, "policy_text", forbid),
					resource.TestCheckResourceAttr(resourceAddress, "policy_status", "DEPLOYED"),
				),
			},
			// The platform reports targets sorted; another declared order plans nothing.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables("Acceptance, renamed", forbid, true),
				ConfigPlanChecks:         emptyPlan,
			},
		},
	})
}

func TestAuthzPolicyResource_immutable_fields(t *testing.T) {
	utils.SkipUnlessAuthzServes(t, utils.AuthzPolicies)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	variables := func(entityId, policyType string) config.Variables {
		return config.Variables{
			"hrid":        config.StringVariable(randomId),
			"entity_id":   config.StringVariable(entityId),
			"policy_type": config.StringVariable(policyType),
		}
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables("policy."+randomId, "MCP"),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables("policy."+randomId+".other", "MCP"),
				ExpectError:              immutableAttribute,
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables("policy."+randomId, "CUSTOM"),
				ExpectError:              immutableAttribute,
			},
		},
	})
}

func TestAuthzPolicyResource_refused(t *testing.T) {
	utils.SkipUnlessAuthzServes(t, utils.AuthzPolicies)
	t.Parallel()

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// APIM's finding reaches the diagnostic.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          config.Variables{"hrid": config.StringVariable("test-" + acctest.RandString(10))},
				ExpectError:              regexp.MustCompile(`Policy templates aren't supported`),
			},
		},
	})
}
