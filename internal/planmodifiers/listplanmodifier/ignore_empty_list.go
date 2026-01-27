package listplanmodifier

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// IgnoreEmptyList returns a plan modifier use the config value if null or unknown when the state contains and empty list
func IgnoreEmptyList() planmodifier.List {
	return ignoreEmptyList{}
}

// suppressDiff implements the plan modifier.
type ignoreEmptyList struct {
	strategy int
}

// Description returns a human-readable description of the plan modifier.
func (m ignoreEmptyList) Description(_ context.Context) string {
	return "Will ignore empty lists returned"
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m ignoreEmptyList) MarkdownDescription(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

// PlanModifyList implements the plan modification logic.
func (m ignoreEmptyList) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {

	// nothing set and plan is [], then use unknown
	if req.ConfigValue.IsNull() && !req.StateValue.IsNull() && len(req.StateValue.Elements()) == 0 {
		resp.PlanValue = req.StateValue
	}
}
