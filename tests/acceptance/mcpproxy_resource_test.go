package acceptance_test

import (
	"regexp"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

var emptyPlan = resource.ConfigPlanChecks{
	PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
}

func TestMcpProxyResource_minimal(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_mcp_proxy.test"
	variables := config.Variables{"hrid": config.StringVariable(randomId)}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// The platform does not contact the upstream server when a proxy is applied.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceAddress, "id"),
					resource.TestCheckResourceAttr(resourceAddress, "mode", "PROXY"),
					resource.TestCheckResourceAttr(resourceAddress, "state", "STARTED"),
					resource.TestCheckResourceAttr(resourceAddress, "proxy.server_url", "https://mcp.example.com/mcp"),
				),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables,
				ConfigPlanChecks:         emptyPlan,
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables,
				ResourceName:             resourceAddress,
				ImportState:              true,
				ImportStateIdFunc:        importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id"}, nil),
				ImportStateVerify:        true,
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

func TestMcpProxyResource_update(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_mcp_proxy.test"

	created := config.Variables{
		"hrid":         config.StringVariable(randomId),
		"name":         config.StringVariable("Acceptance"),
		"context_path": config.StringVariable("/mcp/" + randomId),
		"server_url":   config.StringVariable("https://mcp.example.com/mcp"),
		"state":        config.StringVariable("STARTED"),
		"plans":        config.ListVariable(config.StringVariable("bronze")),
	}
	// Everything but the mode and the entity id converges on a later apply.
	updated := config.Variables{
		"hrid":         config.StringVariable(randomId),
		"name":         config.StringVariable("Acceptance, renamed"),
		"description":  config.StringVariable("Described on a later apply"),
		"context_path": config.StringVariable("/mcp/" + randomId + "-moved"),
		"server_url":   config.StringVariable("https://mcp.example.com/other"),
		"state":        config.StringVariable("STOPPED"),
		"plans":        config.ListVariable(config.StringVariable("bronze"), config.StringVariable("gold")),
	}
	// A plan the list no longer carries is closed.
	closed := config.Variables{
		"hrid":         config.StringVariable(randomId),
		"name":         config.StringVariable("Acceptance, renamed"),
		"description":  config.StringVariable("Described on a later apply"),
		"context_path": config.StringVariable("/mcp/" + randomId + "-moved"),
		"server_url":   config.StringVariable("https://mcp.example.com/other"),
		"state":        config.StringVariable("STOPPED"),
		"plans":        config.ListVariable(config.StringVariable("gold")),
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          created,
				Check:                    resource.TestCheckResourceAttr(resourceAddress, "plans.#", "1"),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          updated,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "name", "Acceptance, renamed"),
					resource.TestCheckResourceAttr(resourceAddress, "description", "Described on a later apply"),
					resource.TestCheckResourceAttr(resourceAddress, "context_path", "/mcp/"+randomId+"-moved"),
					resource.TestCheckResourceAttr(resourceAddress, "proxy.server_url", "https://mcp.example.com/other"),
					resource.TestCheckResourceAttr(resourceAddress, "state", "STOPPED"),
					resource.TestCheckResourceAttr(resourceAddress, "plans.#", "2"),
				),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          updated,
				ConfigPlanChecks:         emptyPlan,
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          closed,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "plans.#", "1"),
					resource.TestCheckResourceAttr(resourceAddress, "plans.0.name", "gold"),
				),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          closed,
				ConfigPlanChecks:         emptyPlan,
			},
		},
	})
}

func TestMcpProxyResource_flows_apikey(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_mcp_proxy.test"
	variables := func(limit int) config.Variables {
		return config.Variables{
			"hrid":  config.StringVariable(randomId),
			"limit": config.IntegerVariable(limit),
		}
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(10),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "flows.#", "1"),
					resource.TestCheckResourceAttr(resourceAddress, "flows.0.selectors.0.mcp.methods.0", "tools/call"),
					resource.TestCheckResourceAttr(resourceAddress, "flows.0.request.0.policy", "rate-limit"),
					resource.TestCheckResourceAttr(resourceAddress, "plans.0.security.api_key.source", "HEADER"),
					resource.TestCheckResourceAttr(resourceAddress, "plans.0.flows.#", "1"),
				),
			},
			// The step configuration carries the policy's defaults, as the platform reports it.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(10),
				ConfigPlanChecks:         emptyPlan,
			},
			// The flows of a proxy and of a plan change in place.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(20),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(20),
				ConfigPlanChecks:         emptyPlan,
			},
		},
	})
}

func TestMcpProxyResource_upstream_auth(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_mcp_proxy.test"
	variables := config.Variables{"hrid": config.StringVariable(randomId)}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// 1: a bearer token. The platform never returns it, and the provider does not read
			// the upstream authentication back: the state holds what was declared.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigVariables:          variables,
				Check:                    resource.TestCheckResourceAttr(resourceAddress, "proxy.upstream_auth.bearer.token", "acceptance-test"),
			},
			// 2: the state keeps the declared token across a read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigVariables:          variables,
				ConfigPlanChecks:         emptyPlan,
				Check:                    resource.TestCheckResourceAttr(resourceAddress, "proxy.upstream_auth.bearer.token", "acceptance-test"),
			},
			// 3: another variant.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigVariables:          variables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "proxy.upstream_auth.api_key.api_key_header", "X-Api-Key"),
					resource.TestCheckNoResourceAttr(resourceAddress, "proxy.upstream_auth.bearer.token"),
				),
			},
			// 4: no upstream authentication, declared by leaving the block out.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigVariables:          variables,
				Check:                    resource.TestCheckNoResourceAttr(resourceAddress, "proxy.upstream_auth.api_key.api_key_header"),
			},
			// 5: and read back as declared.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigVariables:          variables,
				ConfigPlanChecks:         emptyPlan,
			},
		},
	})
}

func TestMcpProxyResource_identity_providers(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_mcp_proxy.test"
	variables := config.Variables{"hrid": config.StringVariable(randomId)}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// 1: a generic provider an OAuth2 plan names.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigVariables:          variables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "identity_providers.0.oauth2_generic.name", "keycloak"),
					resource.TestCheckResourceAttr(resourceAddress, "plans.0.security.oauth2.provider", "keycloak"),
				),
			},
			// 2: the platform reports the issuer without its trailing slash and never
			// returns the client secret. Neither shows as a change.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigVariables:          variables,
				ConfigPlanChecks:         emptyPlan,
			},
			// 3: the platform keeps a provider the declaration no longer names, so
			// removing one is refused when planning.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigVariables:          variables,
				ExpectError:              regexp.MustCompile(`no longer declares keycloak`),
			},
		},
	})
}

func TestMcpProxyResource_immutable_fields(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	variables := func(entityId, mode string) config.Variables {
		return config.Variables{
			"hrid":      config.StringVariable(randomId),
			"entity_id": config.StringVariable(entityId),
			"mode":      config.StringVariable(mode),
		}
	}
	immutable := regexp.MustCompile("This attribute cannot be changed on an existing resource")

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables("mcp-proxy."+randomId, "PROXY"),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables("mcp-proxy."+randomId+"-renamed", "PROXY"),
				ExpectError:              immutable,
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables("mcp-proxy."+randomId, "STUDIO"),
				ExpectError:              immutable,
			},
		},
	})
}

func TestMcpProxyResource_plan_security_change(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	variables := config.Variables{"hrid": config.StringVariable(randomId)}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// 1: a keyless plan.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigVariables:          variables,
			},
			// 2: the same plan with another security. The platform refuses it, and
			// its finding reaches the error.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigVariables:          variables,
				ExpectError:              regexp.MustCompile(`cannot\s+be\s+changed\s+in\s+place`),
			},
		},
	})
}

func TestMcpProxyResource_unsorted_lists(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// The platform reports plans sorted by name. Another order is refused
			// when planning, with the order to declare.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          config.Variables{"hrid": config.StringVariable(randomId)},
				ExpectError:              regexp.MustCompile(`must\s+be\s+sorted\s+by\s+name\.\s+Declare\s+it\s+in\s+this\s+order:\s+bronze,\s+gold`),
			},
		},
	})
}

func TestMcpProxyResource_studio(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	endpoint := utils.McpServerURL(t)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_mcp_proxy.test"
	variables := func(enableFGA bool) config.Variables {
		return config.Variables{
			"hrid":       config.StringVariable(randomId),
			"endpoint":   config.StringVariable(endpoint),
			"enable_fga": config.BoolVariable(enableFGA),
		}
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// A studio exposes tools picked from a catalog server, named by its hrid.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "mode", "STUDIO"),
					resource.TestCheckResourceAttr(resourceAddress, "studio.tools.#", "2"),
					resource.TestCheckResourceAttr(resourceAddress, "studio.tools.0.server", "server-"+randomId),
					resource.TestCheckResourceAttr(resourceAddress, "studio.enable_fga", "true"),
				),
			},
			// The platform composes each tool's entity id. It is in the state after a read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(true),
				ConfigPlanChecks:         emptyPlan,
				Check:                    resource.TestCheckResourceAttrSet(resourceAddress, "studio.tools.0.entity_id"),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(false),
				Check:                    resource.TestCheckResourceAttr(resourceAddress, "studio.enable_fga", "false"),
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

func TestMcpProxyResource_studio_missing_upstream_auth(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	endpoint := utils.McpServerURL(t)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Every server a tool comes from needs an upstream_auth entry. The
			// platform refuses a studio without one, and its finding reaches the error.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":     config.StringVariable(randomId),
					"endpoint": config.StringVariable(endpoint),
				},
				ExpectError: regexp.MustCompile(`(?i)upstreamAuth`),
			},
		},
	})
}
