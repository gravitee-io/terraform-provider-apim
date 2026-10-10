package listplanmodifier_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/gravitee-io/terraform-provider-apim/internal/planmodifiers/listplanmodifier"
)

func TestKeepWhenGroupsIgnored_PlanModifyList(t *testing.T) {
	t.Parallel()

	console, _ := basetypes.NewListValue(types.StringType, []attr.Value{basetypes.NewStringValue("console-group")})
	declared, _ := basetypes.NewListValue(types.StringType, []attr.Value{basetypes.NewStringValue("declared")})
	empty, _ := basetypes.NewListValue(types.StringType, []attr.Value{})
	null := basetypes.NewListNull(types.StringType)

	testCases := map[string]struct {
		ignoreGroups  *bool
		config        basetypes.ListValue
		state         basetypes.ListValue
		plan          basetypes.ListValue
		expected      basetypes.ListValue
		expectWarning bool
	}{
		"ignore_groups true, groups unset: plan keeps the state's groups": {
			ignoreGroups: boolPtr(true), config: null, state: console, plan: empty, expected: console,
		},
		"ignore_groups true, groups declared: plan keeps the state's groups and warns": {
			ignoreGroups: boolPtr(true), config: declared, state: console, plan: declared, expected: console, expectWarning: true,
		},
		"ignore_groups false: plan unchanged": {
			ignoreGroups: boolPtr(false), config: null, state: console, plan: empty, expected: empty,
		},
		"ignore_groups unset: plan unchanged": {
			ignoreGroups: nil, config: null, state: console, plan: empty, expected: empty,
		},
		"no state yet (create): plan unchanged": {
			ignoreGroups: boolPtr(true), config: null, state: null, plan: empty, expected: empty,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := schema.Schema{Attributes: map[string]schema.Attribute{
				"groups":        schema.ListAttribute{ElementType: types.StringType, Optional: true},
				"ignore_groups": schema.BoolAttribute{Optional: true},
			}}
			var ignore tftypes.Value
			if tc.ignoreGroups == nil {
				ignore = tftypes.NewValue(tftypes.Bool, nil)
			} else {
				ignore = tftypes.NewValue(tftypes.Bool, *tc.ignoreGroups)
			}
			raw := tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{
				"groups":        tftypes.List{ElementType: tftypes.String},
				"ignore_groups": tftypes.Bool,
			}}, map[string]tftypes.Value{
				"groups":        tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil),
				"ignore_groups": ignore,
			})

			req := planmodifier.ListRequest{
				Path:        path.Root("groups"),
				ConfigValue: tc.config,
				StateValue:  tc.state,
				PlanValue:   tc.plan,
				Plan:        tfsdk.Plan{Schema: s, Raw: raw},
			}
			resp := &planmodifier.ListResponse{PlanValue: req.PlanValue}

			listplanmodifier.KeepWhenGroupsIgnored().PlanModifyList(context.Background(), req, resp)

			if !resp.PlanValue.Equal(tc.expected) {
				t.Errorf("plan: want %v, got %v", tc.expected, resp.PlanValue)
			}
			if got := resp.Diagnostics.WarningsCount() > 0; got != tc.expectWarning {
				t.Errorf("warning: want %v, got %v", tc.expectWarning, got)
			}
		})
	}
}

func boolPtr(b bool) *bool { return &b }
