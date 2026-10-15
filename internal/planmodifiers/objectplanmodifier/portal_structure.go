package objectplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func PortalStructure() planmodifier.Object {
	return portalStructure{}
}

type portalStructure struct{}

func (portalStructure) Description(context.Context) string {
	return "Requires an explicit top navbar for configured structure and preserves computed structure when omitted."
}

func (m portalStructure) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (portalStructure) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
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
	resp.PlanValue = req.StateValue
}
