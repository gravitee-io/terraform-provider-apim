package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Visibility omission requests inheritance, not the effective value from a previous parent.
func configuredPortalVisibility(ctx context.Context, config tfsdk.Config) (types.String, diag.Diagnostics) {
	var visibility types.String
	diags := config.GetAttribute(ctx, path.Root("visibility"), &visibility)
	return visibility, diags
}
