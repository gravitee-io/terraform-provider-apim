package provider

import (
	"context"

	"github.com/gravitee-io/terraform-provider-apim/internal/sdk/models/operations"
	"github.com/gravitee-io/terraform-provider-apim/internal/sdk/models/shared"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// toDictionaryRequest sends encryptOnFetch false for a dynamic dictionary without an encryption
// policy, so APIM turns encryption off instead of keeping the policy it stores. APIM 4.13.1 and
// earlier ignore the field.
func (data *DictionaryResourceModel) toDictionaryRequest(ctx context.Context) (*operations.CreateOrUpdateDictionariesRequest, diag.Diagnostics) {
	request, diags := data.ToOperationsCreateOrUpdateDictionariesRequest(ctx)
	if request == nil || diags.HasError() {
		return request, diags
	}
	if dynamic := request.DictionarySpec.Dynamic; dynamic != nil && dynamic.Encryption == nil {
		encryptOnFetch := false
		dynamic.Encryption = &shared.DictionaryEncryptionPolicy{EncryptOnFetch: &encryptOnFetch}
	}
	return request, diags
}
