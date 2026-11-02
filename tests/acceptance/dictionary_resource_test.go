package acceptance_test

import (
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Verifies the create, read, import, and delete lifecycle of the
// `apim_dictionary` resource.
func TestDictionaryResource_manual_minimal(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_dictionary.test"

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
				},
			},
			// Verifies resource import.
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
				ImportStateIdFunc: importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id"}, nil),
				ImportStateVerify: true,
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

// Verifies create and update of the `apim_dictionary` resource.
// Name and deployed attributes are updated.
func TestDictionaryResource_manual_update(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":     config.StringVariable(randomId),
					"name":     config.StringVariable("My Dictionary"),
					"deployed": config.BoolVariable(false),
				},
			},
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":     config.StringVariable(randomId),
					"name":     config.StringVariable("My Dictionary - updated"),
					"deployed": config.BoolVariable(false),
				},
			},
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":     config.StringVariable(randomId),
					"name":     config.StringVariable("My Dictionary - updated"),
					"deployed": config.BoolVariable(true),
				},
			},

			// Testing framework implicitly verifies resource delete.
		},
	})
}

// Verifies the create, read, import, and delete lifecycle of the
// `apim_dictionary` resource.
func TestDictionaryResource_dynamic_minimal(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_dictionary.test"

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
				},
			},
			// Verifies resource import.
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
				ImportStateIdFunc: importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id"}, nil),
				ImportStateVerify: true,
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

// Verifies create and update of the `apim_dictionary` resource.
// Name and deployed attributes are updated.
func TestDictionaryResource_dynamic_update(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":     config.StringVariable(randomId),
					"method":   config.StringVariable("GET"),
					"deployed": config.BoolVariable(false),
					"headers": config.ListVariable(config.ObjectVariable(config.Variables{
						"name":  config.StringVariable("X-Test"),
						"value": config.StringVariable("OK"),
					})),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectSensitiveValue("apim_dictionary.test",
							tfjsonpath.New("dynamic").AtMapKey("provider").AtMapKey("http").AtMapKey("url")),
						plancheck.ExpectSensitiveValue("apim_dictionary.test",
							tfjsonpath.New("dynamic").AtMapKey("provider").AtMapKey("http").
								AtMapKey("headers").AtSliceIndex(0).AtMapKey("value")),
					},
				},
			},
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":     config.StringVariable(randomId),
					"method":   config.StringVariable("POST"),
					"deployed": config.BoolVariable(true),
					"headers": config.ListVariable(config.ObjectVariable(config.Variables{
						"name":  config.StringVariable("X-Test-Plus"),
						"value": config.StringVariable("Rock On"),
					})),
				},
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

// Verifies the secret properties of a manual dictionary: APIM never returns
// their values, so Terraform keeps what it sent and shows no drift.
func TestDictionaryResource_manual_encrypted(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_dictionary.test"
	plain := config.MapVariable(map[string]config.Variable{"env": config.StringVariable("test")})
	madePlain := config.MapVariable(map[string]config.Variable{
		"env":     config.StringVariable("test"),
		"api-key": config.StringVariable("not-secret"),
	})
	variables := func(properties config.Variable, secret string, withSecret bool) config.Variables {
		vars := config.Variables{
			"hrid":        config.StringVariable(randomId),
			"secret":      config.StringVariable(secret),
			"with_secret": config.BoolVariable(withSecret),
		}
		if properties != nil {
			vars["properties"] = properties
		}
		return vars
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create with a plain and a secret property.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, "s3cr3t", true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectSensitiveValue(resourceAddress,
							tfjsonpath.New("manual").AtMapKey("encrypted_properties").AtMapKey("api-key").AtMapKey("value")),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "manual.properties.env", "test"),
					resource.TestCheckResourceAttr(resourceAddress, "manual.encrypted_properties.api-key.value", "s3cr3t"),
				),
			},
			// Verifies resource import: secret values are not read back.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"organization_id": config.StringVariable(organizationId),
					"properties":      plain,
					"secret":          config.StringVariable("s3cr3t"),
					"with_secret":     config.BoolVariable(true),
				},
				ResourceName:            resourceAddress,
				ImportState:             true,
				ImportStateIdFunc:       importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id"}, nil),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"manual.encrypted_properties"},
			},
			// Verifies no drift.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, "s3cr3t", true),
				PlanOnly:                 true,
			},
			// Verifies a secret rotation.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, "r0tated", true),
				Check:                    resource.TestCheckResourceAttr(resourceAddress, "manual.encrypted_properties.api-key.value", "r0tated"),
			},
			// Verifies the secret removal.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, "r0tated", false),
				Check:                    resource.TestCheckResourceAttr(resourceAddress, "manual.encrypted_properties.%", "0"),
			},
			// Verifies a dictionary holding only a secret: APIM returns empty plain properties.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(nil, "r0tated", true),
				Check:                    resource.TestCheckResourceAttr(resourceAddress, "manual.properties.%", "0"),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(nil, "r0tated", true),
				PlanOnly:                 true,
			},
			// Verifies a secret made plain again, in two applies: remove it, then add it as plain.
			// APIM rejects a manual dictionary without any property, so a plain one stays.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, "r0tated", false),
				Check:                    resource.TestCheckResourceAttr(resourceAddress, "manual.encrypted_properties.%", "0"),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(madePlain, "r0tated", false),
				Check:                    resource.TestCheckResourceAttr(resourceAddress, "manual.properties.api-key", "not-secret"),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(madePlain, "r0tated", false),
				PlanOnly:                 true,
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

// Verifies the encryption policy of a dynamic dictionary: omitting it keeps
// the current policy.
func TestDictionaryResource_dynamic_encryption(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_dictionary.test"

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create with encryption on fetch.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":             config.StringVariable(randomId),
					"encrypt_on_fetch": config.BoolVariable(true),
				},
				Check: resource.TestCheckResourceAttr(resourceAddress, "dynamic.encryption.encrypt_on_fetch", "true"),
			},
			// Verifies resource import.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":   config.StringVariable(environmentId),
					"hrid":             config.StringVariable(randomId),
					"organization_id":  config.StringVariable(organizationId),
					"encrypt_on_fetch": config.BoolVariable(true),
				},
				ResourceName:      resourceAddress,
				ImportState:       true,
				ImportStateIdFunc: importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id"}, nil),
				ImportStateVerify: true,
			},
			// Verifies no drift.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":             config.StringVariable(randomId),
					"encrypt_on_fetch": config.BoolVariable(true),
				},
				PlanOnly: true,
			},
			// Verifies that omitting the policy keeps it.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.TestCheckResourceAttr(resourceAddress, "dynamic.encryption.encrypt_on_fetch", "true"),
			},
			// Verifies that encryption on fetch can be turned off.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":             config.StringVariable(randomId),
					"encrypt_on_fetch": config.BoolVariable(false),
				},
				Check: resource.TestCheckResourceAttr(resourceAddress, "dynamic.encryption.encrypt_on_fetch", "false"),
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}
