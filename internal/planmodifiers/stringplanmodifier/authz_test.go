package stringplanmodifier_test

import (
	"context"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/internal/planmodifiers/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestEmptyObjectWhenOmitted(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		config   basetypes.StringValue
		plan     basetypes.StringValue
		state    basetypes.StringValue
		expected basetypes.StringValue
	}{
		"left out on create plans an empty object": {
			config:   types.StringNull(),
			plan:     types.StringUnknown(),
			state:    types.StringNull(),
			expected: types.StringValue("{}"),
		},
		"removed from the configuration plans an empty object": {
			config:   types.StringNull(),
			plan:     types.StringValue(`{"department":"travel"}`),
			state:    types.StringValue(`{"department":"travel"}`),
			expected: types.StringValue("{}"),
		},
		"declared is left as planned": {
			config:   types.StringValue(`{"department":"sales"}`),
			plan:     types.StringValue(`{"department":"sales"}`),
			state:    types.StringValue(`{"department":"travel"}`),
			expected: types.StringValue(`{"department":"sales"}`),
		},
		"unknown at plan time is left unknown": {
			config:   types.StringUnknown(),
			plan:     types.StringUnknown(),
			state:    types.StringValue("{}"),
			expected: types.StringUnknown(),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := planmodifier.StringRequest{
				Path:        path.Root("attributes"),
				ConfigValue: testCase.config,
				PlanValue:   testCase.plan,
				StateValue:  testCase.state,
			}
			resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}

			stringplanmodifier.EmptyObjectWhenOmitted().PlanModifyString(context.Background(), req, resp)

			if !resp.PlanValue.Equal(testCase.expected) {
				t.Errorf("planned %v, expected %v", resp.PlanValue, testCase.expected)
			}
		})
	}
}

func TestDeployedOnApply(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		plan     basetypes.StringValue
		state    basetypes.StringValue
		expected basetypes.StringValue
	}{
		"a change to a deployed policy plans it deployed": {
			plan:     types.StringUnknown(),
			state:    types.StringValue("DEPLOYED"),
			expected: types.StringValue("DEPLOYED"),
		},
		"a change to a disabled policy plans it deployed": {
			plan:     types.StringUnknown(),
			state:    types.StringValue("DISABLED"),
			expected: types.StringValue("DEPLOYED"),
		},
		"no change keeps the state's status": {
			plan:     types.StringValue("DISABLED"),
			state:    types.StringValue("DISABLED"),
			expected: types.StringValue("DISABLED"),
		},
		"on create the status stays unknown": {
			plan:     types.StringUnknown(),
			state:    types.StringNull(),
			expected: types.StringUnknown(),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := planmodifier.StringRequest{
				Path:        path.Root("policy_status"),
				ConfigValue: types.StringNull(),
				PlanValue:   testCase.plan,
				StateValue:  testCase.state,
			}
			resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}

			stringplanmodifier.DeployedOnApply().PlanModifyString(context.Background(), req, resp)

			if !resp.PlanValue.Equal(testCase.expected) {
				t.Errorf("planned %v, expected %v", resp.PlanValue, testCase.expected)
			}
		})
	}
}
