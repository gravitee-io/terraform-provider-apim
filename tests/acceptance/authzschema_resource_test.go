package acceptance_test

import (
	"regexp"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAuthzSchemaResource_minimal(t *testing.T) {
	utils.SkipUnlessAuthzServes(t, utils.AuthzSchemas)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_authz_schema.test"
	variables := config.Variables{"hrid": config.StringVariable(randomId)}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Without a target the schema is stored and published nowhere.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceAddress, "id"),
					resource.TestCheckNoResourceAttr(resourceAddress, "target"),
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

func TestAuthzSchemaResource_update(t *testing.T) {
	utils.SkipUnlessAuthzServes(t, utils.AuthzSchemas)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_authz_schema.test"
	variables := func(name, schema string) config.Variables {
		return config.Variables{
			"hrid":   config.StringVariable(randomId),
			"name":   config.StringVariable(name),
			"schema": config.StringVariable(schema),
		}
	}
	created := variables("Acceptance", "entity User in [Group];\nentity Group;\n")
	updated := variables("Acceptance, renamed", "entity User in [Group];\nentity Group { department?: String };\nentity Hotel;\n")

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          created,
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          updated,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "name", "Acceptance, renamed"),
					resource.TestMatchResourceAttr(resourceAddress, "schema", regexp.MustCompile("entity Hotel;")),
				),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          updated,
				ConfigPlanChecks:         emptyPlan,
			},
		},
	})
}

func TestAuthzSchemaResource_target(t *testing.T) {
	utils.SkipUnlessAuthzServes(t, utils.AuthzSchemas, utils.AuthzPdps)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_authz_schema.test"
	variables := func(target int) config.Variables {
		return config.Variables{
			"hrid":   config.StringVariable(randomId),
			"target": config.IntegerVariable(target),
		}
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(0),
				Check:                    resource.TestCheckResourceAttr(resourceAddress, "target", randomId+"-first"),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(0),
				ConfigPlanChecks:         emptyPlan,
			},
			// Moved to another PDP gateway, the schema is withdrawn from the first.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(1),
				Check:                    resource.TestCheckResourceAttr(resourceAddress, "target", randomId+"-second"),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(1),
				ConfigPlanChecks:         emptyPlan,
			},
		},
	})
}
