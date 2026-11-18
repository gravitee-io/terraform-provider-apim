package listplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// declaredPropertyFields are the fields of an API property a configuration sets.
// `encrypted` is left out: APIM never returns it.
var declaredPropertyFields = []string{"key", "value", "dynamic", "encryptable"}

// KeepDynamicProperties returns a plan modifier that plans the stored properties
// of an API when they differ from the planned ones only by properties the dynamic
// properties service added, or by their order. The service adds its properties,
// marked dynamic, on a schedule and sorts the list by key: planning their removal
// would drop them from the API until its next run.
//
// Terraform accepts a list longer than the configuration only when it is the
// stored list, so when a declared property changes, the plan removes the
// service's properties and the service adds them back on its next run.
func KeepDynamicProperties() planmodifier.List {
	return keepDynamicProperties{}
}

type keepDynamicProperties struct{}

func (m keepDynamicProperties) Description(_ context.Context) string {
	return "Properties the dynamic properties service adds are kept while the declared ones do not change."
}

func (m keepDynamicProperties) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m keepDynamicProperties) PlanModifyList(_ context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	held, known := propertiesByKey(req.StateValue)
	if !known {
		return
	}
	planned, known := propertiesByKey(req.PlanValue)
	if !known {
		return
	}
	for key, property := range planned {
		stored, found := held[key]
		if !found || !sameDeclaredFields(property, stored) {
			return
		}
	}
	for key, property := range held {
		if _, declared := planned[key]; !declared && !isDynamic(property) {
			return
		}
	}
	resp.PlanValue = req.StateValue
}

// propertiesByKey indexes the properties by key. known is false when the list or
// a key is null or unknown, or a key is repeated: there is nothing to match.
func propertiesByKey(list basetypes.ListValue) (properties map[string]basetypes.ObjectValue, known bool) {
	if list.IsNull() || list.IsUnknown() {
		return nil, false
	}
	properties = make(map[string]basetypes.ObjectValue, len(list.Elements()))
	for _, element := range list.Elements() {
		property, isObject := element.(basetypes.ObjectValue)
		if !isObject || property.IsNull() || property.IsUnknown() {
			return nil, false
		}
		key, isString := property.Attributes()["key"].(basetypes.StringValue)
		if !isString || key.IsNull() || key.IsUnknown() {
			return nil, false
		}
		if _, repeated := properties[key.ValueString()]; repeated {
			return nil, false
		}
		properties[key.ValueString()] = property
	}
	return properties, true
}

func sameDeclaredFields(planned, stored basetypes.ObjectValue) bool {
	for _, name := range declaredPropertyFields {
		if !planned.Attributes()[name].Equal(stored.Attributes()[name]) {
			return false
		}
	}
	return true
}

func isDynamic(property basetypes.ObjectValue) bool {
	dynamic, isBool := property.Attributes()["dynamic"].(basetypes.BoolValue)
	return isBool && dynamic.ValueBool()
}
