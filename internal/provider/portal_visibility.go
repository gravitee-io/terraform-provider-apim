package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Visibility omission requests inheritance, not the effective value from a previous parent.
func configuredPortalVisibility(ctx context.Context, config tfsdk.Config, plan tfsdk.Plan) (types.String, diag.Diagnostics) {
	var visibility types.String
	diags := config.GetAttribute(ctx, path.Root("visibility"), &visibility)
	if !diags.HasError() && visibility.IsUnknown() {
		diags.Append(plan.GetAttribute(ctx, path.Root("visibility"), &visibility)...)
	}
	if !diags.HasError() && visibility.IsUnknown() {
		diags.AddAttributeError(path.Root("visibility"), "Unresolved portal visibility", "Configured visibility must be known when applying the resource.")
	}
	return visibility, diags
}
