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

// APIM returns both aliases in state, but accepts only one in a PUT request.
func (data *PortalResourceModel) toConfiguredPortalRequest(ctx context.Context, config tfsdk.Config, plan tfsdk.Plan) (*operations.CreateOrUpdatePortalRequest, diag.Diagnostics) {
	var diags diag.Diagnostics
	var configNavigation, planNavigation types.List
	var configStructure, planStructure types.Object
	diags.Append(config.GetAttribute(ctx, path.Root("navigation"), &configNavigation)...)
	diags.Append(config.GetAttribute(ctx, path.Root("structure"), &configStructure)...)
	diags.Append(plan.GetAttribute(ctx, path.Root("navigation"), &planNavigation)...)
	diags.Append(plan.GetAttribute(ctx, path.Root("structure"), &planStructure)...)
	if diags.HasError() {
		return nil, diags
	}

	navigationConfigured := !configNavigation.IsNull() && !(configNavigation.IsUnknown() && planNavigation.IsNull())
	structureConfigured := !configStructure.IsNull() && !(configStructure.IsUnknown() && planStructure.IsNull())
	if navigationConfigured && structureConfigured {
		diags.AddError("Conflicting portal navigation", "navigation and structure cannot be configured together. Configure one navigation representation.")
		return nil, diags
	}

	requestData := *data
	switch {
	case navigationConfigured:
		diags.Append(knownPortalNavigationInput(ctx, "navigation", planNavigation)...)
		if !diags.HasError() {
			diags.Append(plan.GetAttribute(ctx, path.Root("navigation"), &requestData.Navigation)...)
		}
		requestData.Structure = nil
	case structureConfigured:
		diags.Append(knownPortalNavigationInput(ctx, "structure", planStructure)...)
		if !diags.HasError() && planStructure.Attributes()["top_navbar"].IsNull() {
			diags.AddAttributeError(path.Root("structure").AtName("top_navbar"), "Missing portal top navbar", "Set structure.top_navbar to an explicit list. Use [] to clear navigation, or omit structure to preserve navigation.")
		}
		if !diags.HasError() {
			diags.Append(plan.GetAttribute(ctx, path.Root("structure"), &requestData.Structure)...)
		}
		requestData.Navigation = nil
	default:
		// PUT replaces navigation. Preserve refreshed state during unrelated updates.
		if requestData.Structure != nil {
			requestData.Navigation = nil
		}
	}
	if diags.HasError() {
		return nil, diags
	}
	request, requestDiags := requestData.ToOperationsCreateOrUpdatePortalRequest(ctx)
	diags.Append(requestDiags...)
	return request, diags
}

func knownPortalNavigationInput(ctx context.Context, name string, value attr.Value) diag.Diagnostics {
	var diags diag.Diagnostics
	raw, err := value.ToTerraformValue(ctx)
	if err != nil {
		diags.AddAttributeError(path.Root(name), "Invalid portal navigation", err.Error())
	} else if raw.IsNull() || !raw.IsFullyKnown() {
		diags.AddAttributeError(path.Root(name), "Unresolved portal navigation", name+" must have a known value when applying the portal.")
	}
	return diags
}
