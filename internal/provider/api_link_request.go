package provider

import (
	"context"

	"github.com/gravitee-io/terraform-provider-apim/internal/sdk/models/operations"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func (data *APILinkResourceModel) toConfiguredAPILinkRequest(ctx context.Context, config tfsdk.Config, plan tfsdk.Plan) (*operations.CreateOrUpdateAPILinkRequest, diag.Diagnostics) {
	visibility, diags := configuredPortalVisibility(ctx, config, plan)
	if diags.HasError() {
		return nil, diags
	}
	requestData := *data
	requestData.Visibility = visibility
	request, requestDiags := requestData.ToOperationsCreateOrUpdateAPILinkRequest(ctx)
	diags.Append(requestDiags...)
	return request, diags
}
