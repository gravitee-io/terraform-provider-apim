package objectplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func PortalStructure() planmodifier.Object {
	return portalStructure{}
}

type portalStructure struct{}

func (portalStructure) Description(context.Context) string {
	return "Requires an explicit top navbar for configured structure and preserves computed structure unless legacy navigation changes."
}

func (m portalStructure) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (portalStructure) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	var navigation types.List
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("navigation"), &navigation)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() && !navigation.IsNull() && !navigation.IsUnknown() {
		resp.Diagnostics.AddAttributeError(req.Path, "Conflicting portal navigation", "navigation and structure cannot be configured together. Configure one navigation representation.")
		return
	}
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() {
		var topNavbar types.List
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, req.Path.AtName("top_navbar"), &topNavbar)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if topNavbar.IsNull() {
			resp.Diagnostics.AddAttributeError(req.Path.AtName("top_navbar"), "Missing portal top navbar", "Set structure.top_navbar to an explicit list. Use [] to clear navigation, or omit structure to preserve navigation.")
			return
		}
	}
	if !req.PlanValue.IsUnknown() || req.ConfigValue.IsUnknown() || req.State.Raw.IsNull() {
		return
	}
	if !navigation.IsNull() {
		var plannedNavigation, stateNavigation types.List
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("navigation"), &plannedNavigation)...)
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("navigation"), &stateNavigation)...)
		if resp.Diagnostics.HasError() || navigation.IsUnknown() || !plannedNavigation.Equal(stateNavigation) {
			return
		}
	}
	resp.PlanValue = req.StateValue
}
