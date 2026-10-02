package listplanmodifier_test

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/gravitee-io/terraform-provider-apim/internal/planmodifiers/listplanmodifier"
)

// An identity provider is a union: its name sits in the block of the variant declared.
var (
	variantType  = types.ObjectType{AttrTypes: map[string]attr.Type{"name": types.StringType}}
	providerType = types.ObjectType{AttrTypes: map[string]attr.Type{"gravitee_am": variantType, "auth0": variantType}}
)

func providers(names ...attr.Value) basetypes.ListValue {
	elements := make([]attr.Value, 0, len(names))
	for _, name := range names {
		elements = append(elements, types.ObjectValueMust(providerType.AttrTypes, map[string]attr.Value{
			"gravitee_am": types.ObjectNull(variantType.AttrTypes),
			"auth0":       types.ObjectValueMust(variantType.AttrTypes, map[string]attr.Value{"name": name}),
		}))
	}
	return types.ListValueMust(providerType, elements)
}

func TestNoRemovalByName_PlanModifyList(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		config   basetypes.ListValue
		state    basetypes.ListValue
		expected string // text the error must carry, empty when the plan is accepted
	}{
		"state null (creation)": {
			config: providers(types.StringValue("keycloak")),
			state:  types.ListNull(providerType),
		},
		"same providers": {
			config: providers(types.StringValue("am"), types.StringValue("keycloak")),
			state:  providers(types.StringValue("am"), types.StringValue("keycloak")),
		},
		"a provider added": {
			config: providers(types.StringValue("am"), types.StringValue("keycloak")),
			state:  providers(types.StringValue("am")),
		},
		"a provider removed": {
			config:   providers(types.StringValue("am")),
			state:    providers(types.StringValue("am"), types.StringValue("keycloak")),
			expected: "no longer declares keycloak",
		},
		"a provider renamed": {
			config:   providers(types.StringValue("am"), types.StringValue("okta")),
			state:    providers(types.StringValue("am"), types.StringValue("keycloak")),
			expected: "no longer declares keycloak",
		},
		"every provider removed": {
			config:   types.ListNull(providerType),
			state:    providers(types.StringValue("am"), types.StringValue("keycloak")),
			expected: "no longer declares am, keycloak",
		},
		"nothing held, nothing declared": {
			config: types.ListNull(providerType),
			state:  providers(),
		},
		"config unknown": {
			config: types.ListUnknown(providerType),
			state:  providers(types.StringValue("am")),
		},
		"a declared name unknown": {
			config: providers(types.StringUnknown()),
			state:  providers(types.StringValue("am")),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &planmodifier.ListResponse{PlanValue: tc.config}
			listplanmodifier.NoRemovalByName().PlanModifyList(context.Background(), planmodifier.ListRequest{
				Path:        path.Root("identity_providers"),
				ConfigValue: tc.config,
				StateValue:  tc.state,
				PlanValue:   tc.config,
			}, resp)

			if tc.expected == "" {
				if resp.Diagnostics.HasError() {
					t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
				}
				return
			}
			if !resp.Diagnostics.HasError() {
				t.Fatal("expected an error, got none")
			}
			if detail := resp.Diagnostics.Errors()[0].Detail(); !strings.Contains(detail, tc.expected) {
				t.Fatalf("expected the error to carry %q, got %q", tc.expected, detail)
			}
		})
	}
}
