package acceptance_test

import (
	"regexp"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var immutableAttribute = regexp.MustCompile("This attribute cannot be changed on an existing resource")

func authzImportStep(resourceAddress string, variables config.Variables) resource.TestStep {
	return resource.TestStep{
		ProtoV6ProviderFactories: testProviders(),
		ConfigDirectory:          config.TestNameDirectory(),
		ConfigVariables:          variables,
		ResourceName:             resourceAddress,
		ImportState:              true,
		ImportStateIdFunc:        importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id"}, nil),
		ImportStateVerify:        true,
	}
}

func TestAuthzPdpResource_minimal(t *testing.T) {
	utils.SkipUnlessAuthzServes(t, utils.AuthzPdps)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_authz_pdp.test"
	variables := config.Variables{"hrid": config.StringVariable(randomId)}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceAddress, "id"),
					// An untagged PDP gateway is routed to by its targetPdpId alone.
					resource.TestCheckResourceAttr(resourceAddress, "routing_scope", randomId),
					resource.TestCheckResourceAttr(resourceAddress, "strict", "false"),
					resource.TestCheckNoResourceAttr(resourceAddress, "tag"),
					resource.TestCheckNoResourceAttr(resourceAddress, "endpoint"),
				),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables,
				ConfigPlanChecks:         emptyPlan,
			},
			authzImportStep(resourceAddress, variables),
			// Testing framework implicitly verifies resource delete.
		},
	})
}

func TestAuthzPdpResource_tagged_endpoint(t *testing.T) {
	utils.SkipUnlessAuthzServes(t, utils.AuthzPdps)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_authz_pdp.test"
	withEndpoint := config.Variables{
		"hrid":          config.StringVariable(randomId),
		"endpoint_path": config.StringVariable("/authz/" + randomId),
	}
	withoutEndpoint := config.Variables{"hrid": config.StringVariable(randomId)}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          withEndpoint,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "tag", "eu"),
					resource.TestCheckResourceAttr(resourceAddress, "routing_scope", randomId+"@eu"),
					resource.TestCheckResourceAttr(resourceAddress, "endpoint.path", "/authz/"+randomId),
					resource.TestCheckResourceAttrSet(resourceAddress, "well_known_path"),
				),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          withEndpoint,
				ConfigPlanChecks:         emptyPlan,
			},
			// The endpoint is removed in place.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          withoutEndpoint,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceAddress, "endpoint"),
					resource.TestCheckNoResourceAttr(resourceAddress, "well_known_path"),
					resource.TestCheckResourceAttr(resourceAddress, "routing_scope", randomId+"@eu"),
				),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          withoutEndpoint,
				ConfigPlanChecks:         emptyPlan,
			},
		},
	})
}

func TestAuthzPdpResource_immutable_fields(t *testing.T) {
	utils.SkipUnlessAuthzServes(t, utils.AuthzPdps)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	variables := func(targetPdpId, tag string) config.Variables {
		return config.Variables{
			"hrid":          config.StringVariable(randomId),
			"target_pdp_id": config.StringVariable(targetPdpId),
			"tag":           config.StringVariable(tag),
		}
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(randomId, "eu"),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(randomId+"-other", "eu"),
				ExpectError:              immutableAttribute,
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(randomId, "us"),
				ExpectError:              immutableAttribute,
			},
		},
	})
}
