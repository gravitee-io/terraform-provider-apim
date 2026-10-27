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

func TestConfiguredPortalVisibility(t *testing.T) {
	t.Parallel()
	s := schema.Schema{Attributes: map[string]schema.Attribute{
		"visibility": schema.StringAttribute{Optional: true, Computed: true},
	}}
	raw := func(value types.String) types.Object {
		return types.ObjectValueMust(map[string]attr.Type{"visibility": types.StringType}, map[string]attr.Value{"visibility": value})
	}
	for _, tc := range []struct {
		name               string
		config, plan, want types.String
		wantError          bool
	}{
		{"omission does not reuse public state", types.StringNull(), types.StringValue("PUBLIC"), types.StringNull(), false},
		{"omission stays omitted with unknown plan", types.StringNull(), types.StringUnknown(), types.StringNull(), false},
		{"explicit private", types.StringValue("PRIVATE"), types.StringValue("PRIVATE"), types.StringValue("PRIVATE"), false},
		{"resolved expression", types.StringUnknown(), types.StringValue("PRIVATE"), types.StringValue("PRIVATE"), false},
		{"unresolved expression", types.StringUnknown(), types.StringUnknown(), types.StringUnknown(), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			configRaw, err := raw(tc.config).ToTerraformValue(ctx)
			if err != nil {
				t.Fatal(err)
			}
			planRaw, err := raw(tc.plan).ToTerraformValue(ctx)
			if err != nil {
				t.Fatal(err)
			}
			config := tfsdk.Config{Schema: s, Raw: configRaw}
			plan := tfsdk.Plan{Schema: s, Raw: planRaw}
			value, diags := configuredPortalVisibility(ctx, config, plan)
			if diags.HasError() != tc.wantError || !value.Equal(tc.want) {
				t.Fatalf("visibility %s, diagnostics %v", value, diags)
			}
			if tc.wantError {
				return
			}
			model := &PortalLinkResourceModel{Visibility: types.StringValue("PUBLIC")}
			request, diags := model.toConfiguredPortalLinkRequest(ctx, config, plan)
			if diags.HasError() {
				t.Fatal(diags)
			}
			got := request.PortalAttachedLinkSpec.Visibility
			if tc.want.IsNull() && got != nil || !tc.want.IsNull() && (got == nil || string(*got) != tc.want.ValueString()) {
				t.Fatalf("request visibility %v, want %s", got, tc.want)
			}
			if model.Visibility.ValueString() != "PUBLIC" {
				t.Fatal("request construction mutated state")
			}
		})
	}
}

func testConfiguredVisibilityRequest(t *testing.T, build func(context.Context, tfsdk.Config, tfsdk.Plan) (any, diag.Diagnostics)) {
	t.Helper()
	s := schema.Schema{Attributes: map[string]schema.Attribute{
		"visibility": schema.StringAttribute{Optional: true, Computed: true},
		"location":   schema.StringAttribute{Optional: true, Computed: true},
		"order":      schema.Int64Attribute{Optional: true, Computed: true},
	}}
	raw := func(t *testing.T, value types.String) tfsdk.Config {
		t.Helper()
		object := types.ObjectValueMust(
			map[string]attr.Type{"visibility": types.StringType, "location": types.StringType, "order": types.Int64Type},
			map[string]attr.Value{"visibility": value, "location": types.StringNull(), "order": types.Int64Null()},
		)
		valueRaw, err := object.ToTerraformValue(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		return tfsdk.Config{Schema: s, Raw: valueRaw}
	}
	for _, tc := range []struct {
		name         string
		config, plan types.String
		want         string
		wantError    bool
	}{
		{"omitted despite prior public visibility", types.StringNull(), types.StringValue("PUBLIC"), "", false},
		{"omitted with unknown plan", types.StringNull(), types.StringUnknown(), "", false},
		{"explicit private", types.StringValue("PRIVATE"), types.StringValue("PRIVATE"), "PRIVATE", false},
		{"explicit public", types.StringValue("PUBLIC"), types.StringValue("PUBLIC"), "PUBLIC", false},
		{"resolved expression", types.StringUnknown(), types.StringValue("PRIVATE"), "PRIVATE", false},
		{"unresolved expression", types.StringUnknown(), types.StringUnknown(), "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config, plan := raw(t, tc.config), raw(t, tc.plan)
			body, diags := build(context.Background(), config, tfsdk.Plan{Schema: s, Raw: plan.Raw})
			if diags.HasError() != tc.wantError {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			if tc.wantError {
				if body != nil {
					t.Fatal("unresolved visibility must not produce a request")
				}
				return
			}
			payload, err := json.Marshal(body)
			if err != nil {
				t.Fatal(err)
			}
			var attributes map[string]json.RawMessage
			if err := json.Unmarshal(payload, &attributes); err != nil {
				t.Fatal(err)
			}
			value, present := attributes["visibility"]
			if tc.want == "" {
				if present {
					t.Fatalf("omitted visibility must not appear in JSON: %s", payload)
				}
			} else if string(value) != `"`+tc.want+`"` {
				t.Fatalf("request visibility %s, want %q", value, tc.want)
			}
			if string(attributes["hrid"]) != `"test-item"` {
				t.Fatalf("request lost the resource HRID: %s", payload)
			}
		})
	}
}
