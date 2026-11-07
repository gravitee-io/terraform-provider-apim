package acceptance_test

import (
	"regexp"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAuthzEntityResource_minimal(t *testing.T) {
	utils.SkipUnlessAuthzServes(t, utils.AuthzEntities)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_authz_entity.test"
	variables := config.Variables{"hrid": config.StringVariable(randomId)}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceAddress, "id"),
					resource.TestCheckResourceAttr(resourceAddress, "active", "true"),
					resource.TestCheckResourceAttr(resourceAddress, "targets.#", "0"),
					resource.TestCheckResourceAttr(resourceAddress, "parents.#", "0"),
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

func TestAuthzEntityResource_attributes_parents(t *testing.T) {
	utils.SkipUnlessAuthzServes(t, utils.AuthzEntities)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_authz_entity.test"
	variables := func(attributes string, withParent bool) config.Variables {
		return config.Variables{
			"hrid":        config.StringVariable(randomId),
			"attributes":  config.StringVariable(attributes),
			"with_parent": config.BoolVariable(withParent),
		}
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// APIM reports the attributes in its own key order and spacing: the plan
			// after the apply, which the framework checks, is empty.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(`{ "level": 2,  "department": "travel" }`, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "parents.#", "1"),
					resource.TestCheckResourceAttr(resourceAddress, "parents.0.entity_type", "Group"),
					resource.TestCheckResourceAttr(resourceAddress, "parents.0.entity_id", randomId+"-parent"),
				),
			},
			// Attributes and parents are replaced whole.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(`{"department": "sales"}`, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "parents.#", "0"),
					resource.TestMatchResourceAttr(resourceAddress, "attributes", regexp.MustCompile(`^\{\s*"department"\s*:\s*"sales"\s*\}$`)),
				),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(`{"department": "sales"}`, false),
				ConfigPlanChecks:         emptyPlan,
			},
		},
	})
}

// AUTHZ-82: the platform reports targets and parents sorted; another declared
// order plans nothing.
func TestAuthzEntityResource_reordered_lists(t *testing.T) {
	utils.SkipUnlessAuthzServes(t, utils.AuthzEntities, utils.AuthzPdps)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_authz_entity.test"
	variables := func(reversed bool) config.Variables {
		return config.Variables{
			"hrid":     config.StringVariable(randomId),
			"reversed": config.BoolVariable(reversed),
		}
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "targets.#", "2"),
					resource.TestCheckResourceAttr(resourceAddress, "parents.#", "2"),
				),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(true),
				ConfigPlanChecks:         emptyPlan,
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(false),
				ConfigPlanChecks:         emptyPlan,
			},
		},
	})
}

func TestAuthzEntityResource_targets_from_new_pdp(t *testing.T) {
	utils.SkipUnlessAuthzServes(t, utils.AuthzEntities, utils.AuthzPdps)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_authz_entity.test"
	variables := config.Variables{"hrid": config.StringVariable(randomId)}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// The routing scope is unknown until the PDP gateway is created in the same apply.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "targets.#", "1"),
					resource.TestCheckResourceAttr(resourceAddress, "targets.0", randomId+"-pdp"),
				),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables,
				ConfigPlanChecks:         emptyPlan,
			},
		},
	})
}

func TestAuthzEntityResource_immutable_fields(t *testing.T) {
	utils.SkipUnlessAuthzServes(t, utils.AuthzEntities)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	variables := func(entityId, kind, entityType string) config.Variables {
		return config.Variables{
			"hrid":        config.StringVariable(randomId),
			"entity_id":   config.StringVariable(entityId),
			"kind":        config.StringVariable(kind),
			"entity_type": config.StringVariable(entityType),
		}
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(randomId, "PRINCIPAL", "Team"),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(randomId+"-other", "PRINCIPAL", "Team"),
				ExpectError:              immutableAttribute,
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(randomId, "RESOURCE", "Team"),
				ExpectError:              immutableAttribute,
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(randomId, "PRINCIPAL", "Squad"),
				ExpectError:              immutableAttribute,
			},
		},
	})
}

func TestAuthzEntityResource_refused(t *testing.T) {
	utils.SkipUnlessAuthzServes(t, utils.AuthzEntities)
	t.Parallel()

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// APIM's finding reaches the diagnostic.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          config.Variables{"hrid": config.StringVariable("test-" + acctest.RandString(10))},
				ExpectError:              regexp.MustCompile(`not spelled the way the engine spells it`),
			},
		},
	})
}
