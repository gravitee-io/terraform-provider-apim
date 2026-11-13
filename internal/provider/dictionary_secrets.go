package provider

import (
	"context"

	tfTypes "github.com/gravitee-io/terraform-provider-apim/internal/provider/types"
	"github.com/gravitee-io/terraform-provider-apim/internal/sdk/models/shared"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// refreshKeepingSecrets refreshes the model from APIM's answer, which reads every secret back
// masked. A secret the model already holds keeps its value, so it plans no change. A secret only
// APIM holds keeps the mask, and one APIM no longer holds is dropped, so both plan a change. A key
// APIM holds plain that the model holds as a secret stays out of properties: the plan shows the
// secret changing instead of printing the plain value.
func (data *DictionaryResourceModel) refreshKeepingSecrets(ctx context.Context, resp *shared.DictionaryState) diag.Diagnostics {
	var secrets map[string]tfTypes.EncryptableValue
	if data.Manual != nil {
		secrets = data.Manual.EncryptedProperties
	}
	diags := data.RefreshFromSharedDictionaryState(ctx, resp)
	if data.Manual == nil || secrets == nil {
		return diags
	}
	if data.Manual.EncryptedProperties == nil {
		data.Manual.EncryptedProperties = map[string]tfTypes.EncryptableValue{}
	}
	for key := range data.Manual.EncryptedProperties {
		if secret, ok := secrets[key]; ok {
			data.Manual.EncryptedProperties[key] = secret
		}
	}
	for key := range secrets {
		delete(data.Manual.Properties, key)
	}
	if len(data.Manual.Properties) == 0 {
		data.Manual.Properties = nil
	}
	return diags
}
