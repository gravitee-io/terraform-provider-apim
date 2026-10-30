package provider

import (
	"context"

	"github.com/gravitee-io/terraform-provider-apim/internal/sdk/models/operations"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (data *PortalResourceModel) toConfiguredPortalRequest(ctx context.Context, config tfsdk.Config, plan tfsdk.Plan) (*operations.CreateOrUpdatePortalRequest, diag.Diagnostics) {
	var diags diag.Diagnostics
	var configStructure types.Object
	diags.Append(config.GetAttribute(ctx, path.Root("structure"), &configStructure)...)
	if diags.HasError() {
		return nil, diags
	}

	requestData := *data
	if configStructure.IsNull() {
		diags.AddAttributeError(path.Root("structure"), "Missing portal structure", "Configure the required structure object and its complete top_navbar list. Set top_navbar to [] to clear navigation.")
	} else {
		// The structure is the complete desired navigation; reading it from
		// configuration prevents prior state from restoring removed entries.
		configNavbar, configNavbarOK := configStructure.Attributes()["top_navbar"].(types.List)
		if !diags.HasError() && (!configNavbarOK || configNavbar.IsNull()) {
			diags.AddAttributeError(path.Root("structure").AtName("top_navbar"), "Missing portal top navbar", "Configure the complete desired navigation list. Set top_navbar to [] to clear all entries.")
		}
		if !diags.HasError() {
			diags.Append(config.GetAttribute(ctx, path.Root("structure"), &requestData.Structure)...)
		}
	}
	if diags.HasError() {
		return nil, diags
	}

	request, requestDiags := requestData.ToOperationsCreateOrUpdatePortalRequest(ctx)
	diags.Append(requestDiags...)
	return request, diags
}
