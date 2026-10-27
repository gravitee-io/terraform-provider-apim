package provider

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/internal/planmodifiers/objectplanmodifier"
	tfTypes "github.com/gravitee-io/terraform-provider-apim/internal/provider/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
		Optional: true, Computed: true,
		Attributes: map[string]schema.Attribute{
			"top_navbar": schema.ListNestedAttribute{
				Optional: true,
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
		{"explicitly clear structure", emptyStructure, emptyStructure, nil, `{"hrid":"test-portal","name":"Test Portal","structure":{}}`, false},
		{"omitted structure preserves state", nullStructure, types.ObjectUnknown(portalTestStructureTypes), &tfTypes.PortalNavigationStructure{TopNavbar: []tfTypes.PortalNavigationPath{{Path: types.StringValue("/old"), DisplayName: types.StringNull()}}}, `{"hrid":"test-portal","name":"Test Portal","structure":{"topNavbar":[{"path":"/old"}]}}`, false},
		{"create without structure", nullStructure, unknownStructure, nil, `{"hrid":"test-portal","name":"Test Portal"}`, false},
		{"empty structure rejected", emptyStructureConfig, emptyStructureConfig, nil, "", true},
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

func TestPortalStructurePlanModifier(t *testing.T) {
	t.Parallel()
	oldStructure := portalTestStructure(portalTestNavigation("/old"))
	newStructure := portalTestStructure(portalTestNavigation("/new"))
	emptyStructure := portalTestStructure(portalTestNavigation())
	emptyStructureConfig := portalTestStructure(types.ListNull(portalTestPathType))
	nullStructure := types.ObjectNull(portalTestStructureTypes)
	unknownStructure := types.ObjectUnknown(portalTestStructureTypes)

	tests := []struct {
		name      string
		config    types.Object
		plan      types.Object
		state     types.Object
		creating  bool
		want      types.Object
		wantError bool
	}{
		{"omitted structure preserves state", nullStructure, unknownStructure, oldStructure, false, oldStructure, false},
		{"configured structure remains in plan", newStructure, newStructure, oldStructure, false, newStructure, false},
		{"empty navbar is valid", emptyStructure, emptyStructure, oldStructure, false, emptyStructure, false},
		{"missing navbar rejected", emptyStructureConfig, emptyStructureConfig, oldStructure, false, emptyStructureConfig, true},
		{"creation without structure stays unknown", nullStructure, unknownStructure, oldStructure, true, unknownStructure, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			stateRaw := portalTestRaw(t, tc.state)
			if tc.creating {
				stateRaw = tftypes.NewValue(stateRaw.Type(), nil)
			}
			config := tfsdk.Config{Schema: portalTestSchema, Raw: portalTestRaw(t, tc.config)}
			plan := tfsdk.Plan{Schema: portalTestSchema, Raw: portalTestRaw(t, tc.plan)}
			state := tfsdk.State{Schema: portalTestSchema, Raw: stateRaw}
			response := &planmodifier.ObjectResponse{PlanValue: tc.plan}
			objectplanmodifier.PortalStructure().PlanModifyObject(context.Background(), planmodifier.ObjectRequest{
				Path: path.Root("structure"), Config: config, Plan: plan, State: state,
				ConfigValue: tc.config, PlanValue: tc.plan, StateValue: tc.state,
			}, response)
			if response.Diagnostics.HasError() != tc.wantError {
				t.Fatalf("unexpected diagnostics: %v", response.Diagnostics)
			}
			if !response.PlanValue.Equal(tc.want) {
				t.Fatalf("expected plan %s, got %s", tc.want, response.PlanValue)
			}
		})
	}
}
