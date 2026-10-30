package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func configuredDocumentationPlacement(ctx context.Context, config tfsdk.Config) (types.String, types.Int64, diag.Diagnostics) {
	var location types.String
	var order types.Int64
	diags := config.GetAttribute(ctx, path.Root("location"), &location)
	diags.Append(config.GetAttribute(ctx, path.Root("order"), &order)...)
	return location, order, diags
}
