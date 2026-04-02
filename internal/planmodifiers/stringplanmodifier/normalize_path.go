package stringplanmodifier

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// NormalizePath returns a plan modifier that suppresses diffs caused by APIM
// normalizing paths by adding a trailing slash (e.g. "/my-api" → "/my-api/").
func NormalizePath() planmodifier.String {
	return normalizePath{}
}

type normalizePath struct{}

func (m normalizePath) Description(_ context.Context) string {
	return "Suppresses diff when the only difference is a trailing slash added by APIM."
}

func (m normalizePath) MarkdownDescription(_ context.Context) string {
	return "Suppresses diff when the only difference is a trailing slash added by APIM."
}

func (m normalizePath) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	planVal := req.PlanValue.ValueString()
	stateVal := req.StateValue.ValueString()

	if !strings.HasSuffix(planVal, "/") && stateVal == planVal+"/" {
		resp.PlanValue = req.StateValue
	}
}
