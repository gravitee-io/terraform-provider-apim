package stringplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// EmptyObjectWhenOmitted returns a plan modifier that plans `{}` for a JSON
// object left out of the configuration, such as the attributes of an
// authorization entity. The platform replaces the object whole on every apply
// and reports `{}` when it holds nothing.
//
// The attribute has to be computed: Terraform refuses a planned value that
// differs from the configuration of an attribute that is not. Computed, an
// object removed from the configuration would otherwise keep the state's value
// with no change planned.
func EmptyObjectWhenOmitted() planmodifier.String {
	return emptyObjectWhenOmitted{}
}

type emptyObjectWhenOmitted struct{}

func (m emptyObjectWhenOmitted) Description(_ context.Context) string {
	return "Left out of the configuration, the value is an empty object."
}

func (m emptyObjectWhenOmitted) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m emptyObjectWhenOmitted) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.ConfigValue.IsNull() {
		resp.PlanValue = types.StringValue("{}")
	}
}
