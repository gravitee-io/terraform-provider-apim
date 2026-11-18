package listplanmodifier_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/gravitee-io/terraform-provider-apim/internal/planmodifiers/listplanmodifier"
)

var propertyType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"key":         types.StringType,
	"value":       types.StringType,
	"dynamic":     types.BoolType,
	"encryptable": types.BoolType,
	"encrypted":   types.BoolType,
}}

func property(key, value string, dynamic bool) attr.Value {
	return propertyWith(key, value, dynamic, types.BoolNull())
}

func propertyWith(key, value string, dynamic bool, encrypted basetypes.BoolValue) attr.Value {
	return types.ObjectValueMust(propertyType.AttrTypes, map[string]attr.Value{
		"key":         types.StringValue(key),
		"value":       types.StringValue(value),
		"dynamic":     types.BoolValue(dynamic),
		"encryptable": types.BoolNull(),
		"encrypted":   encrypted,
	})
}

func properties(elements ...attr.Value) basetypes.ListValue {
	return types.ListValueMust(propertyType, elements)
}

func TestKeepDynamicProperties_PlanModifyList(t *testing.T) {
	t.Parallel()

	manual := property("manual", "m", false)
	other := property("other", "o", false)
	fromService := property("X-Service", "s", true)

	testCases := map[string]struct {
		plan     basetypes.ListValue
		state    basetypes.ListValue
		expected basetypes.ListValue
	}{
		"a property the service added is kept": {
			plan:     properties(manual),
			state:    properties(fromService, manual),
			expected: properties(fromService, manual),
		},
		"the state's order is kept": {
			plan:     properties(other, manual),
			state:    properties(manual, other),
			expected: properties(manual, other),
		},
		"an unknown encrypted flag is planned from the state": {
			plan:     properties(propertyWith("manual", "m", false, types.BoolUnknown())),
			state:    properties(fromService, manual),
			expected: properties(fromService, manual),
		},
		"a property declared by hand and removed is a change": {
			plan:     properties(manual),
			state:    properties(manual, other),
			expected: properties(manual),
		},
		"an added property is a change": {
			plan:     properties(manual, other),
			state:    properties(fromService, manual),
			expected: properties(manual, other),
		},
		"a changed value is a change": {
			plan:     properties(property("manual", "changed", false)),
			state:    properties(fromService, manual),
			expected: properties(property("manual", "changed", false)),
		},
		"a property the service holds, declared by hand, is a change": {
			plan:     properties(property("X-Service", "s", false)),
			state:    properties(fromService),
			expected: properties(property("X-Service", "s", false)),
		},
		"an empty list keeps the service's properties": {
			plan:     properties(),
			state:    properties(fromService),
			expected: properties(fromService),
		},
		"a repeated key is left as planned": {
			plan:     properties(manual, manual),
			state:    properties(fromService, manual),
			expected: properties(manual, manual),
		},
		"a create is left as planned": {
			plan:     properties(manual),
			state:    types.ListNull(propertyType),
			expected: properties(manual),
		},
		"an unknown list is left as planned": {
			plan:     types.ListUnknown(propertyType),
			state:    properties(fromService, manual),
			expected: types.ListUnknown(propertyType),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := planmodifier.ListRequest{ConfigValue: testCase.plan, PlanValue: testCase.plan, StateValue: testCase.state}
			resp := &planmodifier.ListResponse{PlanValue: req.PlanValue}

			listplanmodifier.KeepDynamicProperties().PlanModifyList(context.Background(), req, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
			if !resp.PlanValue.Equal(testCase.expected) {
				t.Errorf("planned %v, want %v", resp.PlanValue, testCase.expected)
			}
		})
	}
}
