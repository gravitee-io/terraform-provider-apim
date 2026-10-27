package provider

import (
	"context"

	"github.com/gravitee-io/terraform-provider-apim/internal/sdk/models/operations"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
		// Omitted computed fields may be unknown in the plan, but must not be sent from prior state.
		structure := configuredPortalStructure(ctx, configStructure, planStructure)
		diags.Append(knownPortalStructureInput(ctx, structure)...)
		if !diags.HasError() && structure.Attributes()["top_navbar"].IsNull() {
			diags.AddAttributeError(path.Root("structure").AtName("top_navbar"), "Missing portal top navbar", "Set structure.top_navbar to an explicit list. Use [] to clear navigation, or omit structure to preserve navigation.")
		}
		if !diags.HasError() {
			diags.Append(structure.As(ctx, &requestData.Structure, basetypes.ObjectAsOptions{})...)
		}
	}
	if diags.HasError() {
		return nil, diags
	}

	request, requestDiags := requestData.ToOperationsCreateOrUpdatePortalRequest(ctx)
	diags.Append(requestDiags...)
	return request, diags
}

func configuredPortalStructure(ctx context.Context, config, plan types.Object) types.Object {
	if config.IsUnknown() || plan.IsNull() || plan.IsUnknown() {
		return plan
	}
	configNavbar, configOK := config.Attributes()["top_navbar"].(types.List)
	planNavbar, planOK := plan.Attributes()["top_navbar"].(types.List)
	if !configOK || !planOK || configNavbar.IsNull() || configNavbar.IsUnknown() || planNavbar.IsNull() || planNavbar.IsUnknown() {
		return plan
	}
	configured := configNavbar.Elements()
	entries := planNavbar.Elements()
	if len(configured) != len(entries) {
		return plan
	}
	for i, entry := range entries {
		planned, planOK := entry.(types.Object)
		input, configOK := configured[i].(types.Object)
		if !planOK || !configOK || input.IsNull() || input.IsUnknown() || planned.IsNull() || planned.IsUnknown() {
			continue
		}
		if input.Attributes()["visibility"].IsNull() {
			attributes := planned.Attributes()
			attributes["visibility"] = types.StringNull()
			entries[i] = types.ObjectValueMust(planned.AttributeTypes(ctx), attributes)
		}
	}
	attributes := plan.Attributes()
	attributes["top_navbar"] = types.ListValueMust(planNavbar.ElementType(ctx), entries)
	return types.ObjectValueMust(plan.AttributeTypes(ctx), attributes)
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
