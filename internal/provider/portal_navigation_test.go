package provider

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	tfTypes "github.com/gravitee-io/terraform-provider-apim/internal/provider/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var portalTestPathType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"path": types.StringType, "display_name": types.StringType,
	"visibility": types.StringType,
}}

var portalTestStructureTypes = map[string]attr.Type{
	"top_navbar": types.ListType{ElemType: portalTestPathType},
}

var portalTestSchema = schema.Schema{Attributes: map[string]schema.Attribute{
	"structure": schema.SingleNestedAttribute{
		Required: true,
		Attributes: map[string]schema.Attribute{
			"top_navbar": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"path": schema.StringAttribute{Required: true}, "display_name": schema.StringAttribute{Optional: true},
					"visibility": schema.StringAttribute{Optional: true, Computed: true},
				}},
			},
		},
	},
}}

func portalTestNavigation(paths ...string) types.List {
	values := make([]attr.Value, 0, len(paths))
	for _, p := range paths {
		values = append(values, types.ObjectValueMust(portalTestPathType.AttrTypes, map[string]attr.Value{
			"path": types.StringValue(p), "display_name": types.StringNull(),
			"visibility": types.StringNull(),
		}))
	}
	return types.ListValueMust(portalTestPathType, values)
}

func portalTestStructure(navigation types.List) types.Object {
	return types.ObjectValueMust(portalTestStructureTypes, map[string]attr.Value{"top_navbar": navigation})
}

func portalTestRaw(t *testing.T, structure types.Object) tftypes.Value {
	t.Helper()
	value := types.ObjectValueMust(map[string]attr.Type{
		"structure": types.ObjectType{AttrTypes: portalTestStructureTypes},
	}, map[string]attr.Value{"structure": structure})
	raw, err := value.ToTerraformValue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func TestPortalResourceSchema_usesStructureOnly(t *testing.T) {
	var response resource.SchemaResponse
	NewPortalResource().(*PortalResource).Schema(context.Background(), resource.SchemaRequest{}, &response)
	if _, ok := response.Schema.Attributes["navigation"]; ok {
		t.Fatal("portal resource must not expose the deprecated navigation field")
	}
	if _, ok := response.Schema.Attributes["structure"]; !ok {
		t.Fatal("portal resource must expose structure")
	}
	if structure, ok := response.Schema.Attributes["structure"].(schema.SingleNestedAttribute); !ok || !structure.Required {
		t.Fatal("portal resource structure must be required")
	} else if navbar, ok := structure.Attributes["top_navbar"].(schema.ListNestedAttribute); !ok || !navbar.Required {
		t.Fatal("portal resource top_navbar must be required")
	}

	var dataSourceResponse datasource.SchemaResponse
	NewPortalDataSource().(*PortalDataSource).Schema(context.Background(), datasource.SchemaRequest{}, &dataSourceResponse)
	if _, ok := dataSourceResponse.Schema.Attributes["navigation"]; ok {
		t.Fatal("portal data source must not expose the deprecated navigation field")
	}
	if _, ok := dataSourceResponse.Schema.Attributes["structure"]; !ok {
		t.Fatal("portal data source must expose structure")
	}
}

func TestConfiguredPortalRequest(t *testing.T) {
	t.Parallel()
	nullStructure := types.ObjectNull(portalTestStructureTypes)
	unknownStructure := types.ObjectUnknown(portalTestStructureTypes)
	newStructure := portalTestStructure(portalTestNavigation("/new"))
	emptyStructure := portalTestStructure(portalTestNavigation())
	emptyStructureConfig := portalTestStructure(types.ListNull(portalTestPathType))
	unknownPath := types.ListValueMust(portalTestPathType, []attr.Value{
		types.ObjectValueMust(portalTestPathType.AttrTypes, map[string]attr.Value{
			"path": types.StringUnknown(), "display_name": types.StringNull(),
			"visibility": types.StringNull(),
		}),
	})
	unknownPathStructure := portalTestStructure(unknownPath)
	unknownVisibility := portalTestStructure(types.ListValueMust(portalTestPathType, []attr.Value{
		types.ObjectValueMust(portalTestPathType.AttrTypes, map[string]attr.Value{
			"path": types.StringValue("/new"), "display_name": types.StringNull(),
			"visibility": types.StringUnknown(),
		}),
	}))

	tests := []struct {
		name      string
		config    types.Object
		plan      types.Object
		state     *tfTypes.PortalNavigationStructure
		wantJSON  string
		wantError bool
	}{
		{"configured structure", newStructure, newStructure, nil, `{"hrid":"test-portal","name":"Test Portal","structure":{"topNavbar":[{"path":"/new"}]}}`, false},
		{"omitted computed visibility", newStructure, unknownVisibility, nil, `{"hrid":"test-portal","name":"Test Portal","structure":{"topNavbar":[{"path":"/new"}]}}`, false},
		{"configured unknown visibility rejected", unknownVisibility, unknownVisibility, nil, "", true},
		{"explicitly clear top navbar", emptyStructure, emptyStructure, nil, `{"hrid":"test-portal","name":"Test Portal","structure":{"topNavbar":[]}}`, false},
		{"null top navbar rejected", emptyStructureConfig, unknownStructure, nil, "", true},
		{"missing required structure", nullStructure, types.ObjectUnknown(portalTestStructureTypes), nil, "", true},
		{"unknown structure rejected", unknownStructure, unknownStructure, nil, "", true},
		{"unknown path rejected", unknownPathStructure, unknownPathStructure, nil, "", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			model := &PortalResourceModel{
				Hrid: types.StringValue("test-portal"), Name: types.StringValue("Test Portal"),
				EnvironmentID: types.StringValue("DEFAULT"), OrganizationID: types.StringValue("DEFAULT"),
				Structure: tc.state,
			}
			original := *model
			if tc.state != nil {
				structure := *tc.state
				structure.TopNavbar = append([]tfTypes.PortalNavigationPath(nil), tc.state.TopNavbar...)
				original.Structure = &structure
			}
			request, diags := model.toConfiguredPortalRequest(context.Background(),
				tfsdk.Config{Schema: portalTestSchema, Raw: portalTestRaw(t, tc.config)},
				tfsdk.Plan{Schema: portalTestSchema, Raw: portalTestRaw(t, tc.plan)},
			)
			if diags.HasError() != tc.wantError {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			if !reflect.DeepEqual(original, *model) {
				t.Fatal("request selection mutated the state model")
			}
			if tc.wantError {
				if request != nil {
					t.Fatal("expected no request on invalid input")
				}
				return
			}
			body, err := json.Marshal(request.PortalSpec)
			if err != nil {
				t.Fatal(err)
			}
			var actual, expected any
			if err := json.Unmarshal(body, &actual); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal([]byte(tc.wantJSON), &expected); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("expected %s, got %s", tc.wantJSON, body)
			}
		})
	}
}
