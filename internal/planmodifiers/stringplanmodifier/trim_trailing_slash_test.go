package stringplanmodifier_test

import (
	"context"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/internal/planmodifiers/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestTrimTrailingSlash(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		plan     types.String
		state    types.String
		expected types.String
	}{
		"the platform removed the trailing slash - plan gets state value": {
			plan:     types.StringValue("https://idp.example.com/realms/acme/"),
			state:    types.StringValue("https://idp.example.com/realms/acme"),
			expected: types.StringValue("https://idp.example.com/realms/acme"),
		},
		"same value - plan unchanged": {
			plan:     types.StringValue("https://idp.example.com/realms/acme"),
			state:    types.StringValue("https://idp.example.com/realms/acme"),
			expected: types.StringValue("https://idp.example.com/realms/acme"),
		},
		"another URL - plan unchanged": {
			plan:     types.StringValue("https://idp.example.com/realms/other/"),
			state:    types.StringValue("https://idp.example.com/realms/acme"),
			expected: types.StringValue("https://idp.example.com/realms/other/"),
		},
		"slash added rather than removed - plan unchanged": {
			plan:     types.StringValue("https://idp.example.com/realms/acme"),
			state:    types.StringValue("https://idp.example.com/realms/acme/"),
			expected: types.StringValue("https://idp.example.com/realms/acme"),
		},
		"state null (creation) - plan unchanged": {
			plan:     types.StringValue("https://idp.example.com/realms/acme/"),
			state:    types.StringNull(),
			expected: types.StringValue("https://idp.example.com/realms/acme/"),
		},
		"plan unknown - plan unchanged": {
			plan:     types.StringUnknown(),
			state:    types.StringValue("https://idp.example.com/realms/acme"),
			expected: types.StringUnknown(),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &planmodifier.StringResponse{PlanValue: tc.plan}
			stringplanmodifier.TrimTrailingSlash().PlanModifyString(context.Background(), planmodifier.StringRequest{
				PlanValue:  tc.plan,
				StateValue: tc.state,
			}, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
			}
			if !resp.PlanValue.Equal(tc.expected) {
				t.Fatalf("expected plan value %s, got %s", tc.expected, resp.PlanValue)
			}
		})
	}
}
