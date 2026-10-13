package provider

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/internal/planmodifiers/listplanmodifier"
	"github.com/gravitee-io/terraform-provider-apim/internal/planmodifiers/objectplanmodifier"
	tfTypes "github.com/gravitee-io/terraform-provider-apim/internal/provider/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var portalTestPathType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"path": types.StringType, "display_name": types.StringType,
}}

var portalTestStructureTypes = map[string]attr.Type{
	"top_navbar": types.ListType{ElemType: portalTestPathType},
}

var portalTestSchema = schema.Schema{Attributes: map[string]schema.Attribute{
	"navigation": schema.ListNestedAttribute{
		Optional: true, Computed: true,
		NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{Required: true}, "display_name": schema.StringAttribute{Optional: true},
		}},
	},
	"structure": schema.SingleNestedAttribute{
		Optional: true, Computed: true,
		Attributes: map[string]schema.Attribute{
			"top_navbar": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"path": schema.StringAttribute{Required: true}, "display_name": schema.StringAttribute{Optional: true},
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
		}))
	}
	return types.ListValueMust(portalTestPathType, values)
}

func portalTestStructure(navigation types.List) types.Object {
	return types.ObjectValueMust(portalTestStructureTypes, map[string]attr.Value{"top_navbar": navigation})
}

func portalTestRaw(t *testing.T, navigation types.List, structure types.Object) tftypes.Value {
	t.Helper()
	value := types.ObjectValueMust(map[string]attr.Type{
		"navigation": types.ListType{ElemType: portalTestPathType},
		"structure":  types.ObjectType{AttrTypes: portalTestStructureTypes},
	}, map[string]attr.Value{"navigation": navigation, "structure": structure})
	raw, err := value.ToTerraformValue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func TestConfiguredPortalRequest(t *testing.T) {
	t.Parallel()
	nullNavigation := types.ListNull(portalTestPathType)
	unknownNavigation := types.ListUnknown(portalTestPathType)
	nullStructure := types.ObjectNull(portalTestStructureTypes)
	unknownStructure := types.ObjectUnknown(portalTestStructureTypes)
	newNavigation := portalTestNavigation("/new")
	newStructure := portalTestStructure(newNavigation)
	emptyNavigation := portalTestNavigation()
	emptyStructure := portalTestStructure(emptyNavigation)
	emptyStructureConfig := portalTestStructure(nullNavigation)
	unknownPath := types.ListValueMust(portalTestPathType, []attr.Value{
		types.ObjectValueMust(portalTestPathType.AttrTypes, map[string]attr.Value{
			"path": types.StringUnknown(), "display_name": types.StringNull(),
		}),
	})

	tests := []struct {
		name              string
		configNavigation  types.List
		configStructure   types.Object
		planNavigation    types.List
		planStructure     types.Object
		withoutStructure  bool
		withoutNavigation bool
		wantJSON          string
		wantError         bool
	}{
		{"legacy uses resolved plan", newNavigation, nullStructure, newNavigation, unknownStructure, false, false, `{"hrid":"test-portal","name":"Test Portal","navigation":[{"path":"/new"}]}`, false},
		{"structure uses resolved plan", nullNavigation, newStructure, unknownNavigation, newStructure, false, false, `{"hrid":"test-portal","name":"Test Portal","structure":{"topNavbar":[{"path":"/new"}]}}`, false},
		{"clear legacy", emptyNavigation, nullStructure, emptyNavigation, unknownStructure, false, false, `{"hrid":"test-portal","name":"Test Portal"}`, false},
		{"clear structure", nullNavigation, emptyStructure, unknownNavigation, emptyStructure, false, false, `{"hrid":"test-portal","name":"Test Portal","structure":{}}`, false},
		{"empty structure rejected", nullNavigation, emptyStructureConfig, unknownNavigation, emptyStructureConfig, false, false, "", true},
		{"unknown structure resolves to missing navbar", nullNavigation, unknownStructure, unknownNavigation, emptyStructureConfig, false, false, "", true},
		{"omitted preserves structure", nullNavigation, nullStructure, unknownNavigation, unknownStructure, false, false, `{"hrid":"test-portal","name":"Test Portal","structure":{"topNavbar":[{"path":"/old"}]}}`, false},
		{"omitted preserves legacy state", nullNavigation, nullStructure, unknownNavigation, unknownStructure, true, false, `{"hrid":"test-portal","name":"Test Portal","navigation":[{"path":"/old"}]}`, false},
		{"create without navigation", nullNavigation, nullStructure, unknownNavigation, unknownStructure, true, true, `{"hrid":"test-portal","name":"Test Portal"}`, false},
		{"conflicting inputs", newNavigation, newStructure, newNavigation, newStructure, false, false, "", true},
		{"unresolved input", unknownNavigation, nullStructure, unknownNavigation, unknownStructure, false, false, "", true},
		{"unresolved path", unknownPath, nullStructure, unknownPath, unknownStructure, false, false, "", true},
		{"unknown configuration resolves from plan", unknownNavigation, nullStructure, newNavigation, unknownStructure, false, false, `{"hrid":"test-portal","name":"Test Portal","navigation":[{"path":"/new"}]}`, false},
		{"unknown configuration resolves to null", newNavigation, unknownStructure, newNavigation, nullStructure, false, false, `{"hrid":"test-portal","name":"Test Portal","navigation":[{"path":"/new"}]}`, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			model := &PortalResourceModel{
				Hrid: types.StringValue("test-portal"), Name: types.StringValue("Test Portal"),
				EnvironmentID: types.StringValue("DEFAULT"), OrganizationID: types.StringValue("DEFAULT"),
				Navigation: []tfTypes.PortalNavigationPath{{Path: types.StringValue("/old"), DisplayName: types.StringNull()}},
				Structure: &tfTypes.PortalNavigationStructure{TopNavbar: []tfTypes.PortalNavigationPath{
					{Path: types.StringValue("/old"), DisplayName: types.StringNull()},
				}},
			}
			if tc.withoutStructure {
				model.Structure = nil
			}
			if tc.withoutNavigation {
				model.Navigation = nil
			}
			original := *model
			original.Navigation = append([]tfTypes.PortalNavigationPath(nil), model.Navigation...)
			if model.Structure != nil {
				structure := *model.Structure
				structure.TopNavbar = append([]tfTypes.PortalNavigationPath(nil), model.Structure.TopNavbar...)
				original.Structure = &structure
			}
			request, diags := model.toConfiguredPortalRequest(context.Background(),
				tfsdk.Config{Schema: portalTestSchema, Raw: portalTestRaw(t, tc.configNavigation, tc.configStructure)},
				tfsdk.Plan{Schema: portalTestSchema, Raw: portalTestRaw(t, tc.planNavigation, tc.planStructure)},
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

func TestPortalAliasPlanModifiers(t *testing.T) {
	t.Parallel()
	nullNavigation := types.ListNull(portalTestPathType)
	unknownNavigation := types.ListUnknown(portalTestPathType)
	nullStructure := types.ObjectNull(portalTestStructureTypes)
	unknownStructure := types.ObjectUnknown(portalTestStructureTypes)
	oldNavigation := portalTestNavigation("/old")
	oldStructure := portalTestStructure(oldNavigation)
	newNavigation := portalTestNavigation("/new")
	newStructure := portalTestStructure(newNavigation)
	emptyNavigation := portalTestNavigation()
	emptyStructure := portalTestStructure(emptyNavigation)
	emptyStructureConfig := portalTestStructure(nullNavigation)

	tests := []struct {
		name             string
		configNavigation types.List
		configStructure  types.Object
		planNavigation   types.List
		planStructure    types.Object
		wantNavigation   types.List
		wantStructure    types.Object
		creating         bool
		wantError        bool
	}{
		{"unmanaged aliases remain stable", nullNavigation, nullStructure, unknownNavigation, unknownStructure, oldNavigation, oldStructure, false, false},
		{"legacy change recomputes structure", newNavigation, nullStructure, newNavigation, unknownStructure, newNavigation, unknownStructure, false, false},
		{"unchanged legacy preserves structure", oldNavigation, nullStructure, oldNavigation, unknownStructure, oldNavigation, oldStructure, false, false},
		{"structure change recomputes navigation", nullNavigation, newStructure, unknownNavigation, newStructure, unknownNavigation, newStructure, false, false},
		{"unchanged structure preserves navigation", nullNavigation, oldStructure, unknownNavigation, oldStructure, oldNavigation, oldStructure, false, false},
		{"clear legacy recomputes structure", emptyNavigation, nullStructure, emptyNavigation, unknownStructure, emptyNavigation, unknownStructure, false, false},
		{"clear structure recomputes navigation", nullNavigation, emptyStructure, unknownNavigation, emptyStructure, unknownNavigation, emptyStructure, false, false},
		{"unknown legacy keeps structure unknown", unknownNavigation, nullStructure, unknownNavigation, unknownStructure, unknownNavigation, unknownStructure, false, false},
		{"unknown structure keeps navigation unknown", nullNavigation, unknownStructure, unknownNavigation, unknownStructure, unknownNavigation, unknownStructure, false, false},
		{"creation has no prior structure", newNavigation, nullStructure, newNavigation, unknownStructure, newNavigation, unknownStructure, true, false},
		{"conflicting inputs rejected", newNavigation, newStructure, newNavigation, newStructure, newNavigation, newStructure, false, true},
		{"empty structure rejected on creation", nullNavigation, emptyStructureConfig, unknownNavigation, emptyStructureConfig, unknownNavigation, emptyStructureConfig, true, true},
		{"empty structure rejected on update", nullNavigation, emptyStructureConfig, unknownNavigation, emptyStructureConfig, unknownNavigation, emptyStructureConfig, false, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			config := tfsdk.Config{Schema: portalTestSchema, Raw: portalTestRaw(t, tc.configNavigation, tc.configStructure)}
			plan := tfsdk.Plan{Schema: portalTestSchema, Raw: portalTestRaw(t, tc.planNavigation, tc.planStructure)}
			state := tfsdk.State{Schema: portalTestSchema, Raw: portalTestRaw(t, oldNavigation, oldStructure)}
			if tc.creating {
				state.Raw = tftypes.NewValue(state.Raw.Type(), nil)
			}
			objectResponse := &planmodifier.ObjectResponse{PlanValue: tc.planStructure}
			objectplanmodifier.PortalStructure().PlanModifyObject(context.Background(), planmodifier.ObjectRequest{
				Path: path.Root("structure"), Config: config, Plan: plan, State: state,
				ConfigValue: tc.configStructure, PlanValue: tc.planStructure, StateValue: oldStructure,
			}, objectResponse)
			listResponse := &planmodifier.ListResponse{PlanValue: tc.planNavigation}
			listplanmodifier.PortalNavigation().PlanModifyList(context.Background(), planmodifier.ListRequest{
				Path: path.Root("navigation"), Config: config, Plan: plan, State: state,
				ConfigValue: tc.configNavigation, PlanValue: tc.planNavigation, StateValue: oldNavigation,
			}, listResponse)
			if objectResponse.Diagnostics.HasError() != tc.wantError || listResponse.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v, %v", objectResponse.Diagnostics, listResponse.Diagnostics)
			}
			if !objectResponse.PlanValue.Equal(tc.wantStructure) || !listResponse.PlanValue.Equal(tc.wantNavigation) {
				t.Fatalf("unexpected plan: structure=%s, navigation=%s", objectResponse.PlanValue, listResponse.PlanValue)
			}
		})
	}
}
