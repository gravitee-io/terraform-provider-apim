package stringplanmodifier

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// NormalizePath returns a plan modifier that ensures the path value
// always ends with a trailing slash. This prevents drift caused by
// APIM normalizing paths server-side (e.g. "/my-api" → "/my-api/").
func NormalizePath() planmodifier.String {
	return normalizePath{}
}

type normalizePath struct{}

func (m normalizePath) Description(_ context.Context) string {
	return "Normalizes the path value by ensuring it ends with a trailing slash."
}

func (m normalizePath) MarkdownDescription(_ context.Context) string {
	return "Normalizes the path value by ensuring it ends with a trailing slash."
}

func (m normalizePath) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	v := req.PlanValue.ValueString()
	if !strings.HasSuffix(v, "/") {
		resp.PlanValue = req.PlanValue
		// Use the state value if it's the same path with a trailing slash,
		// meaning APIM normalized it.
		if !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
			stateVal := req.StateValue.ValueString()
			if stateVal == v+"/" {
				resp.PlanValue = req.StateValue
			}
		}
	}
}
