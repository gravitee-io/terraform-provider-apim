package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestConfiguredDocumentationPlacement(t *testing.T) {
	t.Parallel()
	s := schema.Schema{Attributes: map[string]schema.Attribute{
		"visibility": schema.StringAttribute{Optional: true, Computed: true},
		"location":   schema.StringAttribute{Optional: true, Computed: true},
		"order":      schema.Int64Attribute{Optional: true, Computed: true},
	}}
	raw := func(location types.String, order types.Int64) types.Object {
		return types.ObjectValueMust(
			map[string]attr.Type{"visibility": types.StringType, "location": types.StringType, "order": types.Int64Type},
			map[string]attr.Value{"visibility": types.StringNull(), "location": location, "order": order},
		)
	}
	builders := map[string]func(context.Context, tfsdk.Config, tfsdk.Plan) (any, diag.Diagnostics){
		"portal": func(ctx context.Context, config tfsdk.Config, plan tfsdk.Plan) (any, diag.Diagnostics) {
			model := &DocumentationPortalResourceModel{Hrid: types.StringValue("item"), Location: types.StringValue("/old"), Order: types.Int64Value(1)}
			request, diags := model.toConfiguredDocumentationPortalRequest(ctx, config, plan)
			if model.Location.ValueString() != "/old" || model.Order.ValueInt64() != 1 {
				t.Fatal("request construction mutated state")
			}
			if request == nil {
				return nil, diags
			}
			return request.DocumentationPortalSpec, diags
		},
		"api": func(ctx context.Context, config tfsdk.Config, plan tfsdk.Plan) (any, diag.Diagnostics) {
			model := &DocumentationAPIResourceModel{Hrid: types.StringValue("item"), Location: types.StringValue("/old"), Order: types.Int64Value(1)}
			request, diags := model.toConfiguredDocumentationAPIRequest(ctx, config, plan)
			if model.Location.ValueString() != "/old" || model.Order.ValueInt64() != 1 {
				t.Fatal("request construction mutated state")
			}
			if request == nil {
				return nil, diags
			}
			return request.DocumentationPortalSpec, diags
		},
	}
	for name, build := range builders {
		for _, tc := range []struct {
			name                         string
			configLocation, planLocation types.String
			configOrder, planOrder       types.Int64
			wantLocation                 string
			wantOrder                    int64
			wantError                    bool
		}{
			{"omitted placement", types.StringNull(), types.StringUnknown(), types.Int64Null(), types.Int64Unknown(), "", 0, false},
			{"explicit placement", types.StringValue("/new"), types.StringValue("/new"), types.Int64Value(2), types.Int64Value(2), "/new", 2, false},
			{"resolved expressions", types.StringUnknown(), types.StringValue("/new"), types.Int64Unknown(), types.Int64Value(2), "/new", 2, false},
			{"unresolved location", types.StringUnknown(), types.StringUnknown(), types.Int64Null(), types.Int64Unknown(), "", 0, true},
			{"unresolved order", types.StringNull(), types.StringUnknown(), types.Int64Unknown(), types.Int64Unknown(), "", 0, true},
		} {
			t.Run(name+"/"+tc.name, func(t *testing.T) {
				ctx := context.Background()
				configRaw, err := raw(tc.configLocation, tc.configOrder).ToTerraformValue(ctx)
				if err != nil {
					t.Fatal(err)
				}
				planRaw, err := raw(tc.planLocation, tc.planOrder).ToTerraformValue(ctx)
				if err != nil {
					t.Fatal(err)
				}
				body, diags := build(ctx, tfsdk.Config{Schema: s, Raw: configRaw}, tfsdk.Plan{Schema: s, Raw: planRaw})
				if diags.HasError() != tc.wantError {
					t.Fatalf("unexpected diagnostics: %v", diags)
				}
				if tc.wantError {
					if body != nil {
						t.Fatal("unresolved placement produced a request")
					}
					return
				}
				payload, err := json.Marshal(body)
				if err != nil {
					t.Fatal(err)
				}
				var fields map[string]json.RawMessage
				if err := json.Unmarshal(payload, &fields); err != nil {
					t.Fatal(err)
				}
				if tc.wantLocation == "" {
					if _, ok := fields["location"]; ok {
						t.Fatalf("unexpected location: %s", payload)
					}
					if _, ok := fields["order"]; ok {
						t.Fatalf("unexpected order: %s", payload)
					}
				} else {
					var location string
					var order int64
					if err := json.Unmarshal(fields["location"], &location); err != nil {
						t.Fatal(err)
					}
					if err := json.Unmarshal(fields["order"], &order); err != nil {
						t.Fatal(err)
					}
					if location != tc.wantLocation || order != tc.wantOrder {
						t.Fatalf("incorrect placement: %s", payload)
					}
				}
			})
		}
	}
}
