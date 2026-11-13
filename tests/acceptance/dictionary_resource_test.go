package acceptance_test

import (
	"regexp"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
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

// Verifies the secret properties of a manual dictionary: APIM returns their
// values masked, so Terraform keeps what it sent and shows no drift.
func TestDictionaryResource_manual_encrypted(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_dictionary.test"
	mask := "••••••••••••"
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
					checkDictionaryEncryptsKey(randomId, "api-key", "s3cr3t"),
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
			// Verifies no drift while APIM reads the secret back masked.
			{
				PreConfig:                requireMaskedRead(t, randomId, "api-key", "s3cr3t"),
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, "s3cr3t", true),
				PlanOnly:                 true,
			},
			// Verifies a secret rotation: a hidden update, and APIM stores the new secret encrypted.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, "r0tated", true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceAddress, plancheck.ResourceActionUpdate),
						plancheck.ExpectSensitiveValue(resourceAddress,
							tfjsonpath.New("manual").AtMapKey("encrypted_properties").AtMapKey("api-key").AtMapKey("value")),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "manual.encrypted_properties.api-key.value", "r0tated"),
					checkDictionaryEncryptsKey(randomId, "api-key", "r0tated"),
				),
			},
			// Verifies an empty secret is refused when planning.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, "", true),
				PlanOnly:                 true,
				ExpectError:              regexp.MustCompile(`UTF-8 character count\s+must be at least 1`),
			},
			// Verifies the mask is accepted: APIM keeps the stored value, which a read cannot show.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, mask, true),
				Check:                    resource.TestCheckResourceAttr(resourceAddress, "manual.encrypted_properties.api-key.value", mask),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, mask, true),
				PlanOnly:                 true,
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
			// Verifies APIM rejects a dictionary left with neither plain nor secret properties.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(nil, "r0tated", false),
				ExpectError:              regexp.MustCompile(`at least one property`),
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

// Verifies the import of a dictionary holding a secret: the state holds no
// secret, so the first plan updates it, and the secret stays encrypted.
func TestDictionaryResource_manual_encrypted_import(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_dictionary.test"
	variables := config.Variables{
		"hrid":   config.StringVariable(randomId),
		"secret": config.StringVariable("s3cr3t"),
	}
	t.Cleanup(func() { deleteDictionary(randomId) })

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies the import plans an update of the secret, hidden.
			{
				PreConfig: func() {
					err := putDictionary(map[string]any{
						"hrid":     randomId,
						"name":     "[Terraform] Imported encrypted dictionary",
						"deployed": false,
						"type":     "MANUAL",
						"manual": map[string]any{
							"properties":          map[string]any{"env": "test"},
							"encryptedProperties": map[string]any{"api-key": map[string]any{"value": "s3cr3t"}},
						},
					})
					if err != nil {
						t.Fatal(err)
					}
				},
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceAddress, plancheck.ResourceActionUpdate),
						plancheck.ExpectSensitiveValue(resourceAddress,
							tfjsonpath.New("manual").AtMapKey("encrypted_properties").AtMapKey("api-key").AtMapKey("value")),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "manual.encrypted_properties.api-key.value", "s3cr3t"),
					checkDictionaryEncryptsKey(randomId, "api-key", "s3cr3t"),
				),
			},
			// Verifies no drift.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables,
				PlanOnly:                 true,
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

// Verifies a dictionary changed outside Terraform. A secret made plain plans an
// update that encrypts it again, without showing the plain value. A secret added
// or deleted plans its removal or its restore. A plain key encrypted outside
// Terraform plans an update that APIM rejects, until the configuration declares
// the key encrypted.
func TestDictionaryResource_manual_encrypted_outside_terraform(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_dictionary.test"
	stringMap := func(values map[string]string) config.Variable {
		variables := make(map[string]config.Variable, len(values))
		for key, value := range values {
			variables[key] = config.StringVariable(value)
		}
		return config.MapVariable(variables)
	}
	variables := func(properties, secrets map[string]string) config.Variables {
		return config.Variables{
			"hrid":       config.StringVariable(randomId),
			"properties": stringMap(properties),
			"secrets":    stringMap(secrets),
		}
	}
	plain := map[string]string{"env": "test"}
	plainWithToken := map[string]string{"env": "test", "token": "plain-token"}
	secret := map[string]string{"api-key": "s3cr3t"}
	secretWithToken := map[string]string{"api-key": "s3cr3t", "token": "console-token"}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, secret),
				Check:                    checkDictionaryEncryptsKey(randomId, "api-key", "s3cr3t"),
			},
			// Verifies a secret made plain outside Terraform plans a change.
			{
				PreConfig:                makePlainOutsideTerraform(t, randomId, "api-key", "now-plain"),
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, secret),
				PlanOnly:                 true,
				ExpectNonEmptyPlan:       true,
			},
			// Verifies the apply encrypts it again, and the plan never shows the plain value.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, secret),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceAddress, plancheck.ResourceActionUpdate),
						expectPlanWithoutValue("now-plain"),
					},
				},
				Check: checkDictionaryEncryptsKey(randomId, "api-key", "s3cr3t"),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, secret),
				PlanOnly:                 true,
			},
			// Verifies a secret added outside Terraform plans its removal.
			{
				PreConfig:                addSecretOutsideTerraform(t, randomId, "added-in-console", "console-value"),
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, secret),
				PlanOnly:                 true,
				ExpectNonEmptyPlan:       true,
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, secret),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceAddress, plancheck.ResourceActionUpdate),
					},
				},
				Check: checkDictionaryLacksKey(randomId, "added-in-console"),
			},
			// Verifies a secret deleted outside Terraform plans its restore.
			{
				PreConfig:                deleteSecretOutsideTerraform(t, randomId, "api-key"),
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, secret),
				PlanOnly:                 true,
				ExpectNonEmptyPlan:       true,
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, secret),
				Check:                    checkDictionaryEncryptsKey(randomId, "api-key", "s3cr3t"),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, secret),
				PlanOnly:                 true,
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plainWithToken, secret),
				Check:                    resource.TestCheckResourceAttr(resourceAddress, "manual.properties.token", "plain-token"),
			},
			// Verifies a plain key encrypted outside Terraform plans a change.
			{
				PreConfig:                encryptOutsideTerraform(t, randomId, "token", "console-token"),
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plainWithToken, secret),
				PlanOnly:                 true,
				ExpectNonEmptyPlan:       true,
			},
			// Verifies APIM rejects the key under properties.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plainWithToken, secret),
				ExpectError:              regexp.MustCompile(`encryptedToPlain`),
			},
			// Verifies declaring the key encrypted converges.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, secretWithToken),
				Check:                    checkDictionaryEncryptsKey(randomId, "token", "console-token"),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          variables(plain, secretWithToken),
				PlanOnly:                 true,
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

// Verifies the encryption policy of a dynamic dictionary: omitting it turns
// encryption off.
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
			// Verifies that omitting the policy turns encryption off.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceAddress, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(resourceAddress,
							tfjsonpath.New("dynamic").AtMapKey("encryption").AtMapKey("encrypt_on_fetch"), knownvalue.Bool(false)),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "dynamic.encryption.encrypt_on_fetch", "false"),
					checkDictionaryEncryptsOnFetch(randomId, false),
				),
			},
			// Verifies no drift while the policy stays omitted.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
				},
				PlanOnly: true,
			},
			// Verifies that encryption on fetch can be turned on again, then off explicitly.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":             config.StringVariable(randomId),
					"encrypt_on_fetch": config.BoolVariable(true),
				},
				Check: checkDictionaryEncryptsOnFetch(randomId, true),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":             config.StringVariable(randomId),
					"encrypt_on_fetch": config.BoolVariable(false),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "dynamic.encryption.encrypt_on_fetch", "false"),
					checkDictionaryEncryptsOnFetch(randomId, false),
				),
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

// Verifies the first apply of a dynamic dictionary whose HRID APIM already holds
// with encryption on: omitting the policy turns it off at once.
func TestDictionaryResource_dynamic_encryption_existing_hrid(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10, utils.ApimV4_11, utils.ApimV4_12)
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_dictionary.test"
	fixture := config.StaticDirectory("testdata/TestDictionaryResource_dynamic_encryption")
	variables := config.Variables{"hrid": config.StringVariable(randomId)}
	t.Cleanup(func() { deleteDictionary(randomId) })

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					err := putDictionary(map[string]any{
						"hrid":     randomId,
						"name":     "[Terraform] Encrypted dynamic dictionary",
						"deployed": false,
						"type":     "DYNAMIC",
						"dynamic": map[string]any{
							"provider": map[string]any{
								"type":          "HTTP",
								"url":           "https://api.gravitee.io/echo",
								"method":        "GET",
								"specification": `[{"operation": "shift", "spec": {"headers": {"*": {"$": "[#2].key", "@": "[#2].value"}}}}]`,
							},
							"trigger":    map[string]any{"rate": 5, "unit": "SECONDS"},
							"encryption": map[string]any{"encryptOnFetch": true},
						},
					})
					if err != nil {
						t.Fatal(err)
					}
				},
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          fixture,
				ConfigVariables:          variables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAddress, "dynamic.encryption.encrypt_on_fetch", "false"),
					checkDictionaryEncryptsOnFetch(randomId, false),
				),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          fixture,
				ConfigVariables:          variables,
				PlanOnly:                 true,
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}
