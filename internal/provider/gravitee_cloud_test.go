package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/internal/sdk/models/shared"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfp "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestCloudInitializerNoCloudAuth(t *testing.T) {
	auth := shared.Security{}
	serverUrl := "fake url"
	var data *ApimProviderModel
	var resp *tfp.ConfigureResponse

	url := CloudInitializer(auth, serverUrl, data, resp)
	assert.Equal(t, url, serverUrl)

}

func TestCloudInitializerValidTokenAndConfig(t *testing.T) {
	file, err := os.ReadFile("testdata/cloud-token.jwt")
	assert.NoError(t, err)

	cloudAuth := string(file)
	auth := shared.Security{
		CloudAuth: &cloudAuth,
	}

	data := &ApimProviderModel{
		EnvironmentID:  types.StringValue("DEFAULT"),
		OrganizationID: types.StringValue("DEFAULT"),
	}
	resp := &tfp.ConfigureResponse{
		Diagnostics: make(diag.Diagnostics, 0),
	}

	url := CloudInitializer(auth, defaultCloudUrl, data, resp)
	assert.Equal(t, url, defaultCloudUrl)
	assert.Len(t, resp.Diagnostics, 0)
	assert.NotEqual(t, "DEFAULT", data.OrganizationID.ValueString())
	assert.NotEqual(t, "DEFAULT", data.EnvironmentID.ValueString())

}

func TestCloudInitializerValidTokenAndNonDefaultURL(t *testing.T) {
	file, err := os.ReadFile("testdata/cloud-token.jwt")
	assert.NoError(t, err)

	cloudAuth := string(file)
	auth := shared.Security{
		CloudAuth: &cloudAuth,
	}

	data := &ApimProviderModel{
		EnvironmentID:  types.StringValue("DEFAULT"),
		OrganizationID: types.StringValue("DEFAULT"),
	}
	resp := &tfp.ConfigureResponse{
		Diagnostics: make(diag.Diagnostics, 0),
	}

	serverUrl := "http://localhost"
	url := CloudInitializer(auth, serverUrl, data, resp)
	assert.Equal(t, url, serverUrl)
	assert.Len(t, resp.Diagnostics, 0)
	assert.NotEqual(t, "DEFAULT", data.OrganizationID.ValueString())
	assert.NotEqual(t, "DEFAULT", data.EnvironmentID.ValueString())

}

func TestCloudInitializerValidTokenNonDefaultGeography(t *testing.T) {
	file, err := os.ReadFile("testdata/cloud-token-us.jwt")
	assert.NoError(t, err)

	cloudAuth := string(file)
	auth := shared.Security{
		CloudAuth: &cloudAuth,
	}

	data := &ApimProviderModel{
		EnvironmentID:  types.StringValue("DEFAULT"),
		OrganizationID: types.StringValue("DEFAULT"),
	}
	resp := &tfp.ConfigureResponse{
		Diagnostics: make(diag.Diagnostics, 0),
	}

	url := CloudInitializer(auth, defaultCloudUrl, data, resp)
	assert.Equal(t, fmt.Sprintf(cloudGateUrlTemplate, "us"), url)
	assert.Len(t, resp.Diagnostics, 0)
	assert.NotEqual(t, "DEFAULT", data.OrganizationID.ValueString())
	assert.NotEqual(t, "DEFAULT", data.EnvironmentID.ValueString())

}

func TestCloudInitializerValidToken2Envs(t *testing.T) {
	file, err := os.ReadFile("testdata/cloud-token-2envs.jwt")
	assert.NoError(t, err)

	cloudAuth := string(file)
	auth := shared.Security{
		CloudAuth: &cloudAuth,
	}

	data := &ApimProviderModel{
		EnvironmentID:  types.StringValue("DEFAULT"),
		OrganizationID: types.StringValue("DEFAULT"),
	}
	resp := &tfp.ConfigureResponse{
		Diagnostics: make(diag.Diagnostics, 0),
	}

	url := CloudInitializer(auth, defaultCloudUrl, data, resp)
	assert.Equal(t, url, defaultCloudUrl)
	assert.Len(t, resp.Diagnostics, 1)
	assert.Equal(t, resp.Diagnostics[0].Severity(), diag.SeverityError)
	assert.Contains(t, resp.Diagnostics[0].Detail(), "environment_id is required")

}

func TestCloudInitializerValidTokenWrongEnv(t *testing.T) {
	file, err := os.ReadFile("testdata/cloud-token.jwt")
	assert.NoError(t, err)

	cloudAuth := string(file)
	auth := shared.Security{
		CloudAuth: &cloudAuth,
	}

	serverUrl := defaultCloudUrl
	data := &ApimProviderModel{
		EnvironmentID:  types.StringValue("foo"),
		OrganizationID: types.StringValue("DEFAULT"),
	}
	resp := &tfp.ConfigureResponse{
		Diagnostics: make(diag.Diagnostics, 0),
	}

	_ = CloudInitializer(auth, serverUrl, data, resp)
	assert.Len(t, resp.Diagnostics, 1)
	assert.Equal(t, resp.Diagnostics[0].Severity(), diag.SeverityError)
	assert.Contains(t, resp.Diagnostics[0].Detail(), "[foo]")
}

func TestCloudInitializerValidTokenWrongOrg(t *testing.T) {
	file, err := os.ReadFile("testdata/cloud-token.jwt")
	assert.NoError(t, err)

	cloudAuth := string(file)
	auth := shared.Security{
		CloudAuth: &cloudAuth,
	}

	serverUrl := defaultCloudUrl
	data := &ApimProviderModel{
		EnvironmentID:  types.StringValue("DEFAULT"),
		OrganizationID: types.StringValue("foo"),
	}
	resp := &tfp.ConfigureResponse{
		Diagnostics: make(diag.Diagnostics, 0),
	}

	_ = CloudInitializer(auth, serverUrl, data, resp)
	assert.Len(t, resp.Diagnostics, 1)
	assert.Equal(t, resp.Diagnostics[0].Severity(), diag.SeverityError)
	assert.Contains(t, resp.Diagnostics[0].Detail(), "[foo]")
}

func TestCloudInitializerValidInvalidToken(t *testing.T) {
	file, err := os.ReadFile("testdata/invalid.jwt")
	assert.NoError(t, err)

	cloudAuth := string(file)
	auth := shared.Security{
		CloudAuth: &cloudAuth,
	}

	serverUrl := defaultCloudUrl
	data := &ApimProviderModel{
		EnvironmentID:  types.StringValue("DEFAULT"),
		OrganizationID: types.StringValue("DEFAULT"),
	}
	resp := &tfp.ConfigureResponse{
		Diagnostics: make(diag.Diagnostics, 0),
	}

	_ = CloudInitializer(auth, serverUrl, data, resp)
	assert.Len(t, resp.Diagnostics, 1)
	assert.Equal(t, resp.Diagnostics[0].Severity(), diag.SeverityError)
	assert.Contains(t, resp.Diagnostics[0].Summary(), "invalid")
	assert.Contains(t, resp.Diagnostics[0].Detail(), "malformed")
}

func TestCloudInitializerValidNotCloudToken(t *testing.T) {
	file, err := os.ReadFile("testdata/not-cloud.jwt")
	assert.NoError(t, err)

	cloudAuth := string(file)
	auth := shared.Security{
		CloudAuth: &cloudAuth,
	}

	serverUrl := defaultCloudUrl
	data := &ApimProviderModel{
		EnvironmentID:  types.StringValue("DEFAULT"),
		OrganizationID: types.StringValue("DEFAULT"),
	}
	resp := &tfp.ConfigureResponse{
		Diagnostics: make(diag.Diagnostics, 0),
	}

	_ = CloudInitializer(auth, serverUrl, data, resp)
	assert.Len(t, resp.Diagnostics, 1)
	assert.Equal(t, resp.Diagnostics[0].Severity(), diag.SeverityError)
	assert.Contains(t, resp.Diagnostics[0].Summary(), "invalid")
	assert.Contains(t, resp.Diagnostics[0].Detail(), "required claims")
}
