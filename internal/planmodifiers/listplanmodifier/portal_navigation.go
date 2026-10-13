package listplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func PortalNavigation() planmodifier.List {
	return portalNavigation{}
}

type portalNavigation struct{}

func (portalNavigation) Description(context.Context) string {
	return "Preserves computed legacy navigation unless configured structure changes."
}

func (m portalNavigation) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (portalNavigation) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if !req.PlanValue.IsUnknown() || req.ConfigValue.IsUnknown() || req.State.Raw.IsNull() {
		return
	}
	var structure types.Object
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("structure"), &structure)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !structure.IsNull() {
		var plannedStructure, stateStructure types.Object
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("structure"), &plannedStructure)...)
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("structure"), &stateStructure)...)
		if resp.Diagnostics.HasError() || structure.IsUnknown() || !plannedStructure.Equal(stateStructure) {
			return
		}
	}
	resp.PlanValue = req.StateValue
}
