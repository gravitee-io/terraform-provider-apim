package stringplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// deployed is the status of an authorization policy after any apply.
const deployed = "DEPLOYED"

// DeployedOnApply returns a plan modifier for the status of an authorization
// policy: an apply deploys the policy, so a planned change to an existing policy
// plans `DEPLOYED` instead of an unknown status.
//
// Without it, a change the other modifiers take back, such as targets declared in
// another order, still plans an update: the framework marks the status unknown
// before the modifiers run. A policy disabled in the Console plans `DISABLED` ->
// `DEPLOYED` on its next change. On create the status stays unknown.
func DeployedOnApply() planmodifier.String {
	return deployedOnApply{}
}

type deployedOnApply struct{}

func (m deployedOnApply) Description(_ context.Context) string {
	return "An apply deploys the policy."
}

func (m deployedOnApply) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m deployedOnApply) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() || !req.PlanValue.IsUnknown() {
		return
	}
	resp.PlanValue = types.StringValue(deployed)
}
