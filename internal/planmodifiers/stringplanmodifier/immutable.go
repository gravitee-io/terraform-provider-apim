package stringplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// Immutable returns a plan modifier that prevents changing the value of an attribute on an existing resource.
// On resource creation (state is null), any value is accepted.
func Immutable() planmodifier.String {
	return immutable{}
}

type immutable struct{}

func (m immutable) Description(_ context.Context) string {
	return "Value cannot be changed after resource creation."
}

func (m immutable) MarkdownDescription(_ context.Context) string {
	return "Value cannot be changed after resource creation."
}

func (m immutable) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Allow any value on creation.
	if req.StateValue.IsNull() {
		return
	}

	// If value hasn't changed, nothing to do.
	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Immutable attribute",
		"This attribute cannot be changed on an existing resource.",
	)
}
