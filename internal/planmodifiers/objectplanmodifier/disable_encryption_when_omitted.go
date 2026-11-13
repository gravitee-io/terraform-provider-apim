package objectplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const encryptOnFetch = "encrypt_on_fetch"

// DisableEncryptionWhenOmitted plans encryption off when the configuration omits the policy. It
// leaves the plan alone on a create, when the dictionary becomes dynamic, and when APIM returned
// no policy: APIM 4.13.1 and earlier ignore it, so a planned one would not match their answer.
func DisableEncryptionWhenOmitted() planmodifier.Object {
	return disableEncryptionWhenOmitted{}
}

// disableEncryptionWhenOmitted implements the plan modifier.
type disableEncryptionWhenOmitted struct {
}

// Description returns a human-readable description of the plan modifier.
func (m disableEncryptionWhenOmitted) Description(_ context.Context) string {
	return "Turns encryption off when the configuration omits the policy."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m disableEncryptionWhenOmitted) MarkdownDescription(_ context.Context) string {
	return "Turns encryption off when the configuration omits the policy."
}

// PlanModifyObject implements the plan modification logic.
func (m disableEncryptionWhenOmitted) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if !req.ConfigValue.IsNull() || req.State.Raw.IsNull() {
		return
	}

	var parent types.Object
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, req.Path.ParentPath(), &parent)...)
	if resp.Diagnostics.HasError() || parent.IsNull() {
		return
	}

	if req.StateValue.IsNull() {
		resp.PlanValue = req.StateValue
		return
	}

	attributes := make(map[string]attr.Value, len(req.StateValue.Attributes()))
	for name, value := range req.StateValue.Attributes() {
		attributes[name] = value
	}
	attributes[encryptOnFetch] = types.BoolValue(false)

	disabled, diags := types.ObjectValue(req.StateValue.AttributeTypes(ctx), attributes)
	resp.Diagnostics.Append(diags...)
	if !diags.HasError() {
		resp.PlanValue = disabled
	}
}
