package listplanmodifier_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/gravitee-io/terraform-provider-apim/internal/planmodifiers/listplanmodifier"
)

var parentType = types.ObjectType{AttrTypes: map[string]attr.Type{"entity_type": types.StringType, "entity_id": types.StringType}}

func targets(values ...attr.Value) basetypes.ListValue {
	return types.ListValueMust(types.StringType, values)
}

func parents(keys ...[2]string) basetypes.ListValue {
	elements := make([]attr.Value, 0, len(keys))
	for _, key := range keys {
		elements = append(elements, types.ObjectValueMust(parentType.AttrTypes, map[string]attr.Value{
			"entity_type": types.StringValue(key[0]),
			"entity_id":   types.StringValue(key[1]),
		}))
	}
	return types.ListValueMust(parentType, elements)
}

func ptr(value basetypes.ListValue) *basetypes.ListValue {
	return &value
}

func TestKeepStateOrder_PlanModifyList(t *testing.T) {
	t.Parallel()

	a, b, c := types.StringValue("a@eu"), types.StringValue("b"), types.StringValue("default")

	testCases := map[string]struct {
		config   *basetypes.ListValue // the plan when nil
		plan     basetypes.ListValue
		state    basetypes.ListValue
		expected basetypes.ListValue
	}{
		"a list left out of the configuration plans an empty one": {
			config:   ptr(types.ListNull(parentType)),
			plan:     parents([2]string{"Group", "agents"}),
			state:    parents([2]string{"Group", "agents"}),
			expected: types.ListValueMust(parentType, []attr.Value{}),
		},
		"a list left out on create plans an empty one": {
			config:   ptr(types.ListNull(types.StringType)),
			plan:     types.ListUnknown(types.StringType),
			state:    types.ListNull(types.StringType),
			expected: targets(),
		},
		"reordered strings plan the state's order": {
			plan:     targets(c, a, b),
			state:    targets(a, b, c),
			expected: targets(a, b, c),
		},
		"reordered objects plan the state's order": {
			plan:     parents([2]string{"Group", "travel"}, [2]string{"Group", "agents"}),
			state:    parents([2]string{"Group", "agents"}, [2]string{"Group", "travel"}),
			expected: parents([2]string{"Group", "agents"}, [2]string{"Group", "travel"}),
		},
		"same order is left as planned": {
			plan:     targets(a, b),
			state:    targets(a, b),
			expected: targets(a, b),
		},
		"an added element is a change": {
			plan:     targets(b, a, c),
			state:    targets(a, b),
			expected: targets(b, a, c),
		},
		"a removed element is a change": {
			plan:     targets(b),
			state:    targets(a, b),
			expected: targets(b),
		},
		"another element is a change": {
			plan:     targets(c, a),
			state:    targets(a, b),
			expected: targets(c, a),
		},
		"repeated elements are counted": {
			plan:     targets(a, a, b),
			state:    targets(a, b, b),
			expected: targets(a, a, b),
		},
		"an element unknown at plan time is left as planned": {
			plan:     targets(types.StringUnknown(), a),
			state:    targets(a, b),
			expected: targets(types.StringUnknown(), a),
		},
		"on create the declared order is planned": {
			plan:     targets(b, a),
			state:    types.ListNull(types.StringType),
			expected: targets(b, a),
		},
		"an emptied list is a change": {
			plan:     targets(),
			state:    targets(a),
			expected: targets(),
		},
		"an unknown list is left unknown": {
			plan:     types.ListUnknown(types.StringType),
			state:    targets(a, b),
			expected: types.ListUnknown(types.StringType),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			config := testCase.plan
			if testCase.config != nil {
				config = *testCase.config
			}
			req := planmodifier.ListRequest{
				Path:        path.Root("targets"),
				ConfigValue: config,
				PlanValue:   testCase.plan,
				StateValue:  testCase.state,
			}
			resp := &planmodifier.ListResponse{PlanValue: req.PlanValue}

			listplanmodifier.KeepStateOrder().PlanModifyList(context.Background(), req, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
			if !resp.PlanValue.Equal(testCase.expected) {
				t.Errorf("planned %v, expected %v", resp.PlanValue, testCase.expected)
			}
		})
	}
}
