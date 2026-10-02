package listvalidators_test

import (
	"context"
	"strings"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/internal/validators/listvalidators"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var planType = types.ObjectType{AttrTypes: map[string]attr.Type{"name": types.StringType}}

func plans(names ...attr.Value) basetypes.ListValue {
	elements := make([]attr.Value, 0, len(names))
	for _, name := range names {
		elements = append(elements, types.ObjectValueMust(planType.AttrTypes, map[string]attr.Value{"name": name}))
	}
	return types.ListValueMust(planType, elements)
}

var toolType = types.ObjectType{AttrTypes: map[string]attr.Type{"server": types.StringType, "tool": types.StringType}}

func tools(pairs ...[2]string) basetypes.ListValue {
	elements := make([]attr.Value, 0, len(pairs))
	for _, pair := range pairs {
		elements = append(elements, types.ObjectValueMust(toolType.AttrTypes, map[string]attr.Value{
			"server": types.StringValue(pair[0]),
			"tool":   types.StringValue(pair[1]),
		}))
	}
	return types.ListValueMust(toolType, elements)
}

// An identity provider is a union: its name sits in the block of the variant declared.
var (
	variantType  = types.ObjectType{AttrTypes: map[string]attr.Type{"name": types.StringType}}
	providerType = types.ObjectType{AttrTypes: map[string]attr.Type{"gravitee_am": variantType, "auth0": variantType}}
)

func providers(declared ...[2]string) basetypes.ListValue {
	elements := make([]attr.Value, 0, len(declared))
	for _, provider := range declared {
		variants := map[string]attr.Value{
			"gravitee_am": types.ObjectNull(variantType.AttrTypes),
			"auth0":       types.ObjectNull(variantType.AttrTypes),
		}
		variants[provider[0]] = types.ObjectValueMust(variantType.AttrTypes, map[string]attr.Value{"name": types.StringValue(provider[1])})
		elements = append(elements, types.ObjectValueMust(providerType.AttrTypes, variants))
	}
	return types.ListValueMust(providerType, elements)
}

func TestSortedBy(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		validator validator.List
		config    basetypes.ListValue
		expected  string // text the error must carry, empty when the list is accepted
	}{
		"sorted by name": {
			validator: listvalidators.SortedByName(),
			config:    plans(types.StringValue("bronze"), types.StringValue("gold")),
		},
		"a single element": {
			validator: listvalidators.SortedByName(),
			config:    plans(types.StringValue("gold")),
		},
		"not sorted by name": {
			validator: listvalidators.SortedByName(),
			config:    plans(types.StringValue("gold"), types.StringValue("bronze")),
			expected:  "must be sorted by name. Declare it in this order: bronze, gold.",
		},
		"upper case sorts before lower case, as the platform does": {
			validator: listvalidators.SortedByName(),
			config:    plans(types.StringValue("alpha"), types.StringValue("Beta")),
			expected:  "Declare it in this order: Beta, alpha.",
		},
		"an unknown name is left to the apply": {
			validator: listvalidators.SortedByName(),
			config:    plans(types.StringValue("gold"), types.StringUnknown()),
		},
		"a null list": {
			validator: listvalidators.SortedByName(),
			config:    types.ListNull(planType),
		},
		"an unknown list": {
			validator: listvalidators.SortedByName(),
			config:    types.ListUnknown(planType),
		},
		"names read in the variant of a union, sorted": {
			validator: listvalidators.SortedByName(),
			config:    providers([2]string{"gravitee_am", "am"}, [2]string{"auth0", "tenant"}),
		},
		"names read in the variant of a union, not sorted": {
			validator: listvalidators.SortedByName(),
			config:    providers([2]string{"auth0", "tenant"}, [2]string{"gravitee_am", "am"}),
			expected:  "Declare it in this order: am, tenant.",
		},
		"sorted by server": {
			validator: listvalidators.SortedByServer(),
			config:    tools([2]string{"deepwiki", "b"}, [2]string{"github", "a"}),
		},
		"not sorted by server": {
			validator: listvalidators.SortedByServer(),
			config:    tools([2]string{"github", "a"}, [2]string{"deepwiki", "b"}),
			expected:  "must be sorted by server. Declare it in this order: deepwiki, github.",
		},
		"sorted by server then tool": {
			validator: listvalidators.SortedByServerAndTool(),
			config:    tools([2]string{"deepwiki", "ask"}, [2]string{"deepwiki", "read"}, [2]string{"github", "ask"}),
		},
		"tools of one server not sorted": {
			validator: listvalidators.SortedByServerAndTool(),
			config:    tools([2]string{"deepwiki", "read"}, [2]string{"deepwiki", "ask"}),
			expected:  "must be sorted by server then tool. Declare it in this order: deepwiki/ask, deepwiki/read.",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &validator.ListResponse{}
			tc.validator.ValidateList(context.Background(), validator.ListRequest{
				Path:        path.Root("plans"),
				ConfigValue: tc.config,
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
