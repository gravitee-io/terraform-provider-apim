package stringplanmodifier

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// TrimTrailingSlash returns a plan modifier that suppresses diffs caused by APIM
// normalizing a URL by removing its trailing slash
// (e.g. "https://idp.example.com/realms/acme/" → "https://idp.example.com/realms/acme").
func TrimTrailingSlash() planmodifier.String {
	return trimTrailingSlash{}
}

type trimTrailingSlash struct{}

func (m trimTrailingSlash) Description(_ context.Context) string {
	return "Suppresses diff when the only difference is a trailing slash removed by APIM."
}

func (m trimTrailingSlash) MarkdownDescription(_ context.Context) string {
	return "Suppresses diff when the only difference is a trailing slash removed by APIM."
}

func (m trimTrailingSlash) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	planVal := req.PlanValue.ValueString()
	stateVal := req.StateValue.ValueString()

	if !strings.HasSuffix(stateVal, "/") && planVal == stateVal+"/" {
		resp.PlanValue = req.StateValue
	}
}
