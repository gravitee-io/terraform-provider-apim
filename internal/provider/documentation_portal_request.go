package provider

import (
	"context"

	"github.com/gravitee-io/terraform-provider-apim/internal/sdk/models/operations"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func (data *DocumentationPortalResourceModel) toConfiguredDocumentationPortalRequest(ctx context.Context, config tfsdk.Config, plan tfsdk.Plan) (*operations.CreateOrUpdatePortalDocumentationRequest, diag.Diagnostics) {
	visibility, diags := configuredPortalVisibility(ctx, config)
	if diags.HasError() {
		return nil, diags
	}
	location, order, placementDiags := configuredDocumentationPlacement(ctx, config)
	diags.Append(placementDiags...)
	if diags.HasError() {
		return nil, diags
	}
	requestData := *data
	requestData.Visibility = visibility
	requestData.Location = location
	requestData.Order = order
	request, requestDiags := requestData.ToOperationsCreateOrUpdatePortalDocumentationRequest(ctx)
	diags.Append(requestDiags...)
	return request, diags
}
