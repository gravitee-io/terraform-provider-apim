package provider

import (
	"context"

	"github.com/gravitee-io/terraform-provider-apim/internal/sdk/models/operations"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (data *PortalResourceModel) toConfiguredPortalRequest(ctx context.Context, config tfsdk.Config, plan tfsdk.Plan) (*operations.CreateOrUpdatePortalRequest, diag.Diagnostics) {
	var diags diag.Diagnostics
	var configStructure, planStructure types.Object
	diags.Append(config.GetAttribute(ctx, path.Root("structure"), &configStructure)...)
	diags.Append(plan.GetAttribute(ctx, path.Root("structure"), &planStructure)...)
	if diags.HasError() {
		return nil, diags
	}

	requestData := *data
	if !configStructure.IsNull() {
		diags.Append(knownPortalStructureInput(ctx, planStructure)...)
		if !diags.HasError() && planStructure.Attributes()["top_navbar"].IsNull() {
			diags.AddAttributeError(path.Root("structure").AtName("top_navbar"), "Missing portal top navbar", "Set structure.top_navbar to an explicit list. Use [] to clear navigation, or omit structure to preserve navigation.")
		}
		if !diags.HasError() {
			diags.Append(plan.GetAttribute(ctx, path.Root("structure"), &requestData.Structure)...)
		}
	}
	if diags.HasError() {
		return nil, diags
	}

	request, requestDiags := requestData.ToOperationsCreateOrUpdatePortalRequest(ctx)
	diags.Append(requestDiags...)
	return request, diags
}

func knownPortalStructureInput(ctx context.Context, value attr.Value) diag.Diagnostics {
	var diags diag.Diagnostics
	raw, err := value.ToTerraformValue(ctx)
	if err != nil {
		diags.AddAttributeError(path.Root("structure"), "Invalid portal structure", err.Error())
	} else if raw.IsNull() || !raw.IsFullyKnown() {
		diags.AddAttributeError(path.Root("structure"), "Unresolved portal structure", "structure must have a known value when applying the portal.")
	}
	return diags
}
