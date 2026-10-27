package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func configuredDocumentationPlacement(ctx context.Context, config tfsdk.Config, plan tfsdk.Plan) (types.String, types.Int64, diag.Diagnostics) {
	var location types.String
	var order types.Int64
	diags := config.GetAttribute(ctx, path.Root("location"), &location)
	diags.Append(config.GetAttribute(ctx, path.Root("order"), &order)...)
	if !diags.HasError() && location.IsUnknown() {
		diags.Append(plan.GetAttribute(ctx, path.Root("location"), &location)...)
	}
	if !diags.HasError() && order.IsUnknown() {
		diags.Append(plan.GetAttribute(ctx, path.Root("order"), &order)...)
	}
	if !diags.HasError() && location.IsUnknown() {
		diags.AddAttributeError(path.Root("location"), "Unresolved documentation location", "Configured location must be known when applying the resource.")
	}
	if !diags.HasError() && order.IsUnknown() {
		diags.AddAttributeError(path.Root("order"), "Unresolved documentation order", "Configured order must be known when applying the resource.")
	}
	return location, order, diags
}
