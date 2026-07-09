package acceptance_test

import (
	"encoding/base64"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cloudTestOrg  = "2806dece-d045-463d-86de-ced045963d84"
	cloudTestEnv  = "9b5e326f-f68a-45f9-9e32-6ff68ae5f92c"
	cloudTestEnv2 = "357c35e3-6635-443a-9812-312e2f548606"
	cloudWrongOrg = "00000000-0000-0000-0000-000000000001"
	cloudWrongEnv = "11111111-1111-1111-1111-111111111111"
)

// Verifies cloud token scope validation rejects misconfiguration before API calls.
func TestCloudAuthScope_rejectsUnauthorizedScope(t *testing.T) {
	
	randomID := "test-" + acctest.RandString(10)
	validToken := buildCloudToken(t, validCloudTokenClaims())
	multiEnvToken := buildCloudToken(t, validCloudTokenClaims(cloudTestEnv, cloudTestEnv2))
	nonCloudToken := buildCloudToken(t, map[string]any{
		"sub": "1234567890",
	})

	baseVars := config.Variables{
		"hrid": config.StringVariable(randomID),
	}

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: mergeVariables(baseVars, config.Variables{
					"cloud_token": config.StringVariable("not-a-jwt"),
				}),
				ExpectError: regexp.MustCompile(`Cloud Token invalid`),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: mergeVariables(baseVars, config.Variables{
					"cloud_token": config.StringVariable(nonCloudToken),
				}),
				ExpectError: regexp.MustCompile(`required claims`),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: mergeVariables(baseVars, config.Variables{
					"cloud_token": config.StringVariable(multiEnvToken),
				}),
				ExpectError: regexp.MustCompile(`Cloud Token incompatible`),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: mergeVariables(baseVars, config.Variables{
					"cloud_token":             config.StringVariable(validToken),
					"provider_environment_id": config.StringVariable(cloudWrongEnv),
				}),
				ExpectError: regexp.MustCompile(`Cloud Token misconfiguration`),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: mergeVariables(baseVars, config.Variables{
					"cloud_token":              config.StringVariable(validToken),
					"provider_organization_id": config.StringVariable(cloudWrongOrg),
				}),
				ExpectError: regexp.MustCompile(`Cloud Token misconfiguration`),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: mergeVariables(baseVars, config.Variables{
					"cloud_token":             config.StringVariable(validToken),
					"resource_environment_id": config.StringVariable(cloudWrongEnv),
				}),
				ExpectError: regexp.MustCompile(`cloud token scope violation`),
			},
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: mergeVariables(baseVars, config.Variables{
					"cloud_token":              config.StringVariable(validToken),
					"resource_organization_id": config.StringVariable(cloudWrongOrg),
				}),
				ExpectError: regexp.MustCompile(`cloud token scope violation`),
			},
		},
	})
}

func mergeVariables(base config.Variables, overrides config.Variables) config.Variables {
	merged := make(config.Variables, len(base)+len(overrides))
	for key, value := range base {
		merged[key] = value
	}
	for key, value := range overrides {
		merged[key] = value
	}
	return merged
}

func buildCloudToken(t *testing.T, claims map[string]any) string {
	t.Helper()

	header := base64.RawURLEncoding.EncodeToString([]byte(`{"typ":"JWT","alg":"none"}`))
	payload, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("failed to marshal cloud token claims: %v", err)
	}

	return header + "." + base64.RawURLEncoding.EncodeToString(payload) + "."
}

func validCloudTokenClaims(envs ...string) map[string]any {
	if len(envs) == 0 {
		envs = []string{cloudTestEnv}
	}

	return map[string]any{
		"org":    cloudTestOrg,
		"envs":   envs,
		"cpg":    "eu",
		"iss":    "GraviteeCloud",
		"aud":    "CloudGate",
		"target": "apim",
		"scopes": []string{"automation", "rest"},
	}
}
