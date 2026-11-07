package listplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// KeepStateOrder returns a plan modifier that plans the state's order when the
// planned list holds the same elements as the state, in another order. It is for
// the lists the platform reports sorted, such as the targets of an authorization
// entity: the provider keeps the declared order after an apply and takes the
// platform's on the next read, so a list declared in another order would plan a
// change on every run.
//
// The attribute has to be computed: Terraform refuses a planned value that
// differs from the configuration of an attribute that is not. Computed, a list
// left out of the configuration would keep the state's elements with no change
// planned, so it plans an empty list instead: the platform replaces these lists
// whole on every apply.
func KeepStateOrder() planmodifier.List {
	return keepStateOrder{}
}

type keepStateOrder struct{}

func (m keepStateOrder) Description(_ context.Context) string {
	return "The order of the elements makes no difference."
}

func (m keepStateOrder) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m keepStateOrder) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.ConfigValue.IsNull() {
		empty, diags := types.ListValue(req.PlanValue.ElementType(ctx), []attr.Value{})
		resp.Diagnostics.Append(diags...)
		resp.PlanValue = empty
		return
	}
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() || req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}
	planned, err := req.PlanValue.ToTerraformValue(ctx)
	if err != nil || !planned.IsFullyKnown() {
		return
	}
	if sameElements(req.PlanValue.Elements(), req.StateValue.Elements()) {
		resp.PlanValue = req.StateValue
	}
}

// sameElements tells whether both lists hold the same elements as many times
// each, in any order.
func sameElements(planned, held []attr.Value) bool {
	if len(planned) != len(held) {
		return false
	}
	matched := make([]bool, len(held))
	for _, element := range planned {
		found := false
		for i, candidate := range held {
			if !matched[i] && element.Equal(candidate) {
				matched[i] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
