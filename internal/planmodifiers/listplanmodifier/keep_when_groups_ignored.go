package listplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// KeepWhenGroupsIgnored returns a plan modifier for a resource's groups: while the resource sets
// ignore_groups = true, APIM owns the groups and an apply leaves them as they are, so the plan keeps the
// state's value instead of showing a change the apply will not make.
func KeepWhenGroupsIgnored() planmodifier.List {
	return keepWhenGroupsIgnored{}
}

type keepWhenGroupsIgnored struct{}

func (m keepWhenGroupsIgnored) Description(_ context.Context) string {
	return "Keeps the groups of the state while ignore_groups is true."
}

func (m keepWhenGroupsIgnored) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m keepWhenGroupsIgnored) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	var ignoreGroups types.Bool
	if diags := req.Plan.GetAttribute(ctx, path.Root("ignore_groups"), &ignoreGroups); diags.HasError() {
		return
	}
	if ignoreGroups.IsNull() || ignoreGroups.IsUnknown() || !ignoreGroups.ValueBool() {
		return
	}
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() && len(req.ConfigValue.Elements()) > 0 {
		resp.Diagnostics.AddAttributeWarning(req.Path, "groups is not applied while ignore_groups is true",
			"APIM keeps the groups it has. Remove groups from the configuration, or set ignore_groups = false to manage them here.")
	}
	resp.PlanValue = req.StateValue
}
