package provider

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gravitee-io/terraform-provider-apim/internal/sdk/models/shared"
	tfp "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const cloudGateUrlTemplate = "https://%s.cloudgate.gravitee.io/apim/automation"

var defaultCloudUrl = fmt.Sprintf(cloudGateUrlTemplate, "eu")

// CloudInitResult holds provider configuration derived from a Gravitee Cloud token.
type CloudInitResult struct {
	ServerURL string
	Claims    *CloudTokenClaimsData
}

type orgEnv struct {
	Org string
	Env string
}

func fromData(data *ApimProviderModel) orgEnv {
	return orgEnv{
		Org: data.OrganizationID.ValueString(),
		Env: data.EnvironmentID.ValueString(),
	}
}

// CloudInitializer setup url/env/org if required. Keeps non default value intact.
func CloudInitializer(ctx context.Context, auth shared.Security, serverUrl string, data *ApimProviderModel, resp *tfp.ConfigureResponse) CloudInitResult {
	if auth.CloudAuth == nil {
		return CloudInitResult{ServerURL: serverUrl}
	}

	jwtData, err := extractCloudTokenData(*auth.CloudAuth)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Token invalid", err.Error())
		return CloudInitResult{ServerURL: serverUrl}
	}

	if configuredEnvID := data.EnvironmentID.ValueString(); configuredEnvID == "DEFAULT" {
		// Set if unset
		// Check only one is present
		if len(jwtData.Envs) > 1 {
			resp.Diagnostics.AddError(
				"Cloud Token incompatible",
				fmt.Sprintf("cloud token contains more than one environment (%d): environment_id is required in that case", len(jwtData.Envs)))
			return CloudInitResult{ServerURL: serverUrl}
		}
		data.EnvironmentID = basetypes.NewStringValue(jwtData.Envs[0])
	} else if err := validateCloudScope(jwtData, fromData(data)); err != nil {
		tflog.Error(ctx, "cloud token scope violation", map[string]interface{}{
			"error": err,
		})
		resp.Diagnostics.AddError("Cloud Token misconfiguration", err.Error())
		return CloudInitResult{ServerURL: serverUrl}
	}

	// Set if unset
	if configuredOrgID := data.OrganizationID.ValueString(); configuredOrgID == "DEFAULT" {
		data.OrganizationID = basetypes.NewStringValue(jwtData.Org)
	} else if err := validateCloudScope(jwtData, fromData(data)); err != nil {
		tflog.Error(ctx, "cloud token scope violation", map[string]interface{}{
			"error": err,
		})
		resp.Diagnostics.AddError("Cloud Token misconfiguration", err.Error())
		return CloudInitResult{ServerURL: serverUrl}
	}

	claims := jwtData
	result := CloudInitResult{
		ServerURL: serverUrl,
		Claims:    &claims,
	}

	// Return user defined URL
	if serverUrl != defaultCloudUrl {
		return result
	}

	// returned computed URL
	result.ServerURL = jwtData.baseUrl()
	return result
}

func validateCloudScope(claims CloudTokenClaimsData, orgEnv orgEnv) error {
	if orgEnv.Env != "" && orgEnv.Env != "DEFAULT" && !slices.Contains(claims.Envs, orgEnv.Env) {
		return fmt.Errorf("cloud token does not contain environment [%s], it must be one of: %s", orgEnv.Env, claims.Envs)
	}
	if orgEnv.Org != "" && orgEnv.Org != "DEFAULT" && orgEnv.Org != claims.Org {
		return fmt.Errorf("cloud token specifies organization [%s], you cannot use [%s] for organization_id value", claims.Org, orgEnv.Org)
	}
	return nil
}

func extractCloudTokenData(jwtToken string) (CloudTokenClaimsData, error) {
	claims := &CloudTokenClaims{}

	_, _, err := jwt.NewParser().ParseUnverified(jwtToken, claims)

	if err != nil {
		return CloudTokenClaimsData{}, err
	}

	if !claims.isValid() {
		return CloudTokenClaimsData{}, errors.New("cloud token does not contains all required claims")
	}

	return claims.CloudTokenClaimsData, nil
}

type CloudTokenClaims struct {
	jwt.RegisteredClaims
	CloudTokenClaimsData
}

type CloudTokenClaimsData struct {
	Org       string   `json:"org"`
	Envs      []string `json:"envs"`
	Geography string   `json:"cpg"`
}

func (d CloudTokenClaimsData) baseUrl() string {
	return fmt.Sprintf(cloudGateUrlTemplate, d.Geography)
}

func (d CloudTokenClaimsData) isValid() bool {
	return d.Org != "" && d.Envs != nil && len(d.Envs) > 0 && d.Geography != ""
}
