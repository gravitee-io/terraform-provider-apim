package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tfTypes "github.com/gravitee-io/terraform-provider-apim/internal/provider/types"
)

func dynamicDictionary(encryption *tfTypes.DictionaryEncryptionPolicy) DictionaryResourceModel {
	return DictionaryResourceModel{
		Hrid:     types.StringValue("fetched"),
		Name:     types.StringValue("fetched"),
		Type:     types.StringValue("DYNAMIC"),
		Deployed: types.BoolValue(false),
		Dynamic: &tfTypes.DynamicDictionarySpec{
			Provider: &tfTypes.DictionaryProvider{HTTP: &tfTypes.HTTPDictionaryProvider{
				Type:          types.StringValue("HTTP"),
				URL:           types.StringValue("https://api.gravitee.io/echo"),
				Method:        types.StringValue("GET"),
				Specification: types.StringValue("[]"),
			}},
			Trigger:    &tfTypes.DictionaryTrigger{Rate: types.Int64Value(5), Unit: types.StringValue("SECONDS")},
			Encryption: encryption,
		},
	}
}

func TestDictionaryRequestTurnsEncryptionOffWhenThePolicyIsOmitted(t *testing.T) {
	data := dynamicDictionary(nil)

	request, diags := data.toDictionaryRequest(context.Background())

	require.False(t, diags.HasError(), diags)
	require.NotNil(t, request.DictionarySpec.Dynamic.Encryption)
	require.NotNil(t, request.DictionarySpec.Dynamic.Encryption.EncryptOnFetch)
	assert.False(t, *request.DictionarySpec.Dynamic.Encryption.EncryptOnFetch)
}

func TestDictionaryRequestKeepsADeclaredPolicy(t *testing.T) {
	data := dynamicDictionary(&tfTypes.DictionaryEncryptionPolicy{EncryptOnFetch: types.BoolValue(true)})

	request, diags := data.toDictionaryRequest(context.Background())

	require.False(t, diags.HasError(), diags)
	require.NotNil(t, request.DictionarySpec.Dynamic.Encryption.EncryptOnFetch)
	assert.True(t, *request.DictionarySpec.Dynamic.Encryption.EncryptOnFetch)
}

func TestDictionaryRequestSendsNoPolicyForAManualDictionary(t *testing.T) {
	data := dictionaryWithSecret("s3cr3t")

	request, diags := data.toDictionaryRequest(context.Background())

	require.False(t, diags.HasError(), diags)
	assert.Nil(t, request.DictionarySpec.Dynamic)
}
