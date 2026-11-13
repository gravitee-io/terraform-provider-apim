package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tfTypes "github.com/gravitee-io/terraform-provider-apim/internal/provider/types"
	"github.com/gravitee-io/terraform-provider-apim/internal/sdk"
	"github.com/gravitee-io/terraform-provider-apim/internal/sdk/models/shared"
)

// From APIM 4.13.3, the Automation API answers every encrypted dictionary value with a mask, on a
// read and on a write alike. The answer goes through the SDK and refreshKeepingSecrets here, so a
// mask that reaches the state fails in this test instead of as a change in every plan.
const maskedDictionaryAnswer = `{
	"hrid": "secrets",
	"name": "secrets",
	"type": "MANUAL",
	"deployed": false,
	"organizationId": "DEFAULT",
	"environmentId": "DEFAULT",
	"manual": {
		"properties": {"env": "test"},
		"encryptedProperties": {"api-key": {"value": "••••••••••••"}}
	}
}`

func maskedDictionaryClient(t *testing.T) *sdk.GraviteeApim {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(maskedDictionaryAnswer))
	}))
	t.Cleanup(server.Close)
	return sdk.New(
		sdk.WithServerURL(server.URL),
		sdk.WithSecurity(shared.Security{BasicAuth: &shared.SchemeBasicAuth{Username: "admin", Password: "admin"}}),
	)
}

func dictionaryWithSecret(secret string) DictionaryResourceModel {
	return DictionaryResourceModel{
		Hrid:           types.StringValue("secrets"),
		Name:           types.StringValue("secrets"),
		Type:           types.StringValue("MANUAL"),
		Deployed:       types.BoolValue(false),
		OrganizationID: types.StringValue("DEFAULT"),
		EnvironmentID:  types.StringValue("DEFAULT"),
		Manual: &tfTypes.ManualDictionarySpec{
			Properties: map[string]types.String{"env": types.StringValue("test")},
			EncryptedProperties: map[string]tfTypes.EncryptableValue{
				"api-key": {Value: types.StringValue(secret)},
			},
		},
	}
}

func readDictionary(t *testing.T, data *DictionaryResourceModel) {
	t.Helper()
	ctx := context.Background()
	request, diags := data.ToOperationsGetDictionaryRequest(ctx)
	require.False(t, diags.HasError(), diags)
	res, err := maskedDictionaryClient(t).Dictionaries.Get(ctx, *request)
	require.NoError(t, err)
	require.NotNil(t, res.DictionaryState)
	diags = data.refreshKeepingSecrets(ctx, res.DictionaryState)
	require.False(t, diags.HasError(), diags)
}

func TestDictionaryReadKeepsTheSecretOverTheMask(t *testing.T) {
	data := dictionaryWithSecret("s3cr3t")

	readDictionary(t, &data)

	require.NotNil(t, data.Manual)
	assert.Equal(t, "s3cr3t", data.Manual.EncryptedProperties["api-key"].Value.ValueString())
	assert.Equal(t, map[string]types.String{"env": types.StringValue("test")}, data.Manual.Properties)
}

func TestDictionaryWriteAnswerKeepsTheSecretOverTheMask(t *testing.T) {
	ctx := context.Background()
	data := dictionaryWithSecret("r0tated")
	request, diags := data.ToOperationsCreateOrUpdateDictionariesRequest(ctx)
	require.False(t, diags.HasError(), diags)

	res, err := maskedDictionaryClient(t).Dictionaries.CreateOrUpdate(ctx, *request)
	require.NoError(t, err)
	require.NotNil(t, res.DictionaryState)
	diags = data.refreshKeepingSecrets(ctx, res.DictionaryState)
	require.False(t, diags.HasError(), diags)

	require.NotNil(t, data.Manual)
	assert.Equal(t, "r0tated", data.Manual.EncryptedProperties["api-key"].Value.ValueString())
}

// After an import there is no prior state to keep the secret from: the key holds the mask, and the
// configured value then plans as an update.
func TestDictionaryImportHoldsTheMask(t *testing.T) {
	data := DictionaryResourceModel{
		Hrid:           types.StringValue("secrets"),
		OrganizationID: types.StringValue("DEFAULT"),
		EnvironmentID:  types.StringValue("DEFAULT"),
	}

	readDictionary(t, &data)

	require.NotNil(t, data.Manual)
	assert.Equal(t, "••••••••••••", data.Manual.EncryptedProperties["api-key"].Value.ValueString())
	assert.Equal(t, map[string]types.String{"env": types.StringValue("test")}, data.Manual.Properties)
}
