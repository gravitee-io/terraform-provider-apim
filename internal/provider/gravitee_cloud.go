package provider

import (
	"errors"
	"fmt"
	"slices"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gravitee-io/terraform-provider-apim/internal/sdk/models/shared"
	tfp "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const cloudGateUrlTemplate = "https://%s.cloudgate.gravitee.io/apim/automation"

var defaultCloudUrl = fmt.Sprintf(cloudGateUrlTemplate, "eu")

// CloudInitializer setup url/env/org if required. Keeps non default value intact
func CloudInitializer(auth shared.Security, serverUrl string, data *ApimProviderModel, resp *tfp.ConfigureResponse) string {
	if auth.CloudAuth == nil {
		return serverUrl
	}

	jwtData, err := extractCloudTokenData(*auth.CloudAuth)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Token invalid", err.Error())
		return serverUrl
	}

	if configuredEnvID := data.EnvironmentID.ValueString(); configuredEnvID == "DEFAULT" {
		// Set if unset
		// Check only one is present
		if len(jwtData.Envs) > 1 {
			resp.Diagnostics.AddError(
				"Cloud Token incompatible",
				fmt.Sprintf("cloud token contains more than one environment (%d): environment_id is required in that case", len(jwtData.Envs)))
			return serverUrl
		}
		data.EnvironmentID = basetypes.NewStringValue(jwtData.Envs[0])
	} else if !slices.Contains(jwtData.Envs, configuredEnvID) {
		// if set, check env in token
		resp.Diagnostics.AddError("Cloud Token misconfiguration",
			fmt.Sprintf("cloud token does not contain environment [%s], it must be one of: %s", configuredEnvID, jwtData.Envs))
		return serverUrl
	}

	// Set if unset
	if configuredOrgID := data.OrganizationID.ValueString(); configuredOrgID == "DEFAULT" {
		data.OrganizationID = basetypes.NewStringValue(jwtData.Org)
	} else if configuredOrgID != jwtData.Org {
		// if set, check env in token
		resp.Diagnostics.AddError("Cloud Token misconfiguration",
			fmt.Sprintf("cloud token specifies organization [%s] you cannot use [%s] for organization_id value", jwtData.Org, configuredOrgID))
		return serverUrl
	}

	// Return user defined URL
	if serverUrl != defaultCloudUrl {
		return serverUrl
	}

	// returned computed URL
	return jwtData.baseUrl()
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
