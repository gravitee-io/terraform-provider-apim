package mapplanmodifier_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/gravitee-io/terraform-provider-apim/internal/planmodifiers/mapplanmodifier"
)

var responseTemplateAttrTypes = map[string]attr.Type{
	"status":                      types.Int64Type,
	"headers":                     types.MapType{ElemType: types.StringType},
	"body":                        types.StringType,
	"propagate_error_key_to_logs": types.BoolType,
}

var responseTemplateType = types.ObjectType{AttrTypes: responseTemplateAttrTypes}

func newResponseTemplateObject(t *testing.T, status int64, body string, propagate bool, headers basetypes.MapValue) basetypes.ObjectValue {
	t.Helper()

	obj, diags := basetypes.NewObjectValue(responseTemplateAttrTypes, map[string]attr.Value{
		"status":                      basetypes.NewInt64Value(status),
		"headers":                     headers,
		"body":                        basetypes.NewStringValue(body),
		"propagate_error_key_to_logs": basetypes.NewBoolValue(propagate),
	})
	if diags.HasError() {
		t.Fatalf("failed to build response template object: %v", diags)
	}
	return obj
}

func newResponseTemplatesMap(t *testing.T, templates map[string]map[string]basetypes.ObjectValue) basetypes.MapValue {
	t.Helper()

	outerElementType := types.MapType{ElemType: responseTemplateType}
	outer := make(map[string]attr.Value, len(templates))
	for code, contentTypes := range templates {
		inner := make(map[string]attr.Value, len(contentTypes))
		for contentType, obj := range contentTypes {
			inner[contentType] = obj
		}
		innerMap, diags := basetypes.NewMapValue(responseTemplateType, inner)
		if diags.HasError() {
			t.Fatalf("failed to build content-type map for %q: %v", code, diags)
		}
		outer[code] = innerMap
	}

	outerMap, diags := basetypes.NewMapValue(outerElementType, outer)
	if diags.HasError() {
		t.Fatalf("failed to build response_templates map: %v", diags)
	}
	return outerMap
}

func emptyHeaders(t *testing.T) basetypes.MapValue {
	t.Helper()
	m, diags := basetypes.NewMapValue(types.StringType, map[string]attr.Value{})
	if diags.HasError() {
		t.Fatalf("failed to build empty headers map: %v", diags)
	}
	return m
}

func nullHeaders() basetypes.MapValue {
	return basetypes.NewMapNull(types.StringType)
}

func nonEmptyHeaders(t *testing.T) basetypes.MapValue {
	t.Helper()
	m, diags := basetypes.NewMapValue(types.StringType, map[string]attr.Value{
		"X-Error": basetypes.NewStringValue("boom"),
	})
	if diags.HasError() {
		t.Fatalf("failed to build non-empty headers map: %v", diags)
	}
	return m
}

func runResponseTemplatePlanModify(config, state, plan basetypes.MapValue) basetypes.MapValue {
	req := planmodifier.MapRequest{
		ConfigValue: config,
		StateValue:  state,
		PlanValue:   plan,
	}
	resp := &planmodifier.MapResponse{PlanValue: plan}
	mapplanmodifier.ResponseTemplate().PlanModifyMap(context.Background(), req, resp)
	return resp.PlanValue
}

func TestResponseTemplate_PlanModifyMap(t *testing.T) {
	t.Parallel()

	emptyMap, _ := basetypes.NewMapValue(types.StringType, map[string]attr.Value{})
	nonEmptyMap, _ := basetypes.NewMapValue(types.StringType, map[string]attr.Value{
		"key": basetypes.NewStringValue("value"),
	})
	nullMap := basetypes.NewMapNull(types.StringType)
	unknownMap := basetypes.NewMapUnknown(types.StringType)

	testCases := map[string]struct {
		configValue   basetypes.MapValue
		stateValue    basetypes.MapValue
		planValue     basetypes.MapValue
		expectedValue basetypes.MapValue
	}{
		"config null + state empty map - plan gets state value": {
			configValue:   nullMap,
			stateValue:    emptyMap,
			planValue:     nullMap,
			expectedValue: emptyMap,
		},
		"config null + state non-empty map - plan unchanged": {
			configValue:   nullMap,
			stateValue:    nonEmptyMap,
			planValue:     nullMap,
			expectedValue: nullMap,
		},
		"config null + state null - plan unchanged": {
			configValue:   nullMap,
			stateValue:    nullMap,
			planValue:     nullMap,
			expectedValue: nullMap,
		},
		"config set - plan unchanged": {
			configValue:   nonEmptyMap,
			stateValue:    emptyMap,
			planValue:     nonEmptyMap,
			expectedValue: nonEmptyMap,
		},
		"config null + state unknown - plan gets state value": {
			configValue:   nullMap,
			stateValue:    unknownMap,
			planValue:     nullMap,
			expectedValue: unknownMap,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := runResponseTemplatePlanModify(tc.configValue, tc.stateValue, tc.planValue)
			if !got.Equal(tc.expectedValue) {
				t.Errorf("Expected plan value %s, got %s", tc.expectedValue, got)
			}
		})
	}
}

func TestResponseTemplate_SuppressesResponseTemplateHeadersEmptyMapDrift(t *testing.T) {
	t.Parallel()

	configObj := newResponseTemplateObject(t, 400, "http method override denied", false, emptyHeaders(t))
	stateObj := newResponseTemplateObject(t, 400, "http method override denied", false, nullHeaders())

	outerConfig := newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
		"INVALID_HTTP_METHOD": {"*/*": configObj},
	})
	outerState := newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
		"INVALID_HTTP_METHOD": {"*/*": stateObj},
	})

	got := runResponseTemplatePlanModify(outerConfig, outerState, outerConfig)
	if !got.Equal(outerState) {
		t.Errorf("Expected plan to be state value; got plan %s, state %s", got, outerState)
	}
}

func TestResponseTemplate_DoesNotSuppressRealDiffs(t *testing.T) {
	t.Parallel()

	baseConfigObj := newResponseTemplateObject(t, 400, "empty headers", false, emptyHeaders(t))
	baseStateObj := newResponseTemplateObject(t, 400, "empty headers", false, nullHeaders())

	testCases := map[string]struct {
		config basetypes.MapValue
		state  basetypes.MapValue
	}{
		"status differs": {
			config: newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
				"EMPTY_HEADERS": {"*/*": newResponseTemplateObject(t, 400, "empty headers", false, emptyHeaders(t))},
			}),
			state: newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
				"EMPTY_HEADERS": {"*/*": newResponseTemplateObject(t, 500, "empty headers", false, nullHeaders())},
			}),
		},
		"body differs": {
			config: newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
				"EMPTY_HEADERS": {"*/*": newResponseTemplateObject(t, 400, "empty headers", false, emptyHeaders(t))},
			}),
			state: newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
				"EMPTY_HEADERS": {"*/*": newResponseTemplateObject(t, 400, "different body", false, nullHeaders())},
			}),
		},
		"propagate_error_key_to_logs differs": {
			config: newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
				"EMPTY_HEADERS": {"*/*": newResponseTemplateObject(t, 400, "empty headers", true, emptyHeaders(t))},
			}),
			state: newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
				"EMPTY_HEADERS": {"*/*": newResponseTemplateObject(t, 400, "empty headers", false, nullHeaders())},
			}),
		},
		"headers are non-empty in config": {
			config: newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
				"EMPTY_HEADERS": {"*/*": newResponseTemplateObject(t, 400, "empty headers", false, nonEmptyHeaders(t))},
			}),
			state: newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
				"EMPTY_HEADERS": {"*/*": newResponseTemplateObject(t, 400, "empty headers", false, nullHeaders())},
			}),
		},
		"headers are non-empty in state": {
			config: newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
				"EMPTY_HEADERS": {"*/*": newResponseTemplateObject(t, 400, "empty headers", false, emptyHeaders(t))},
			}),
			state: newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
				"EMPTY_HEADERS": {"*/*": newResponseTemplateObject(t, 400, "empty headers", false, nonEmptyHeaders(t))},
			}),
		},
		"config template missing from state": {
			config: newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
				"EMPTY_HEADERS": {"*/*": baseConfigObj},
			}),
			state: newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
				"OTHER_TEMPLATE": {"*/*": baseStateObj},
			}),
		},
		"state has extra template": {
			config: newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
				"EMPTY_HEADERS": {"*/*": baseConfigObj},
			}),
			state: newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
				"EMPTY_HEADERS":  {"*/*": baseStateObj},
				"OTHER_TEMPLATE": {"*/*": baseStateObj},
			}),
		},
		"state has extra content-type": {
			config: newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
				"EMPTY_HEADERS": {"*/*": baseConfigObj},
			}),
			state: newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
				"EMPTY_HEADERS": {
					"*/*":              baseStateObj,
					"application/json": baseStateObj,
				},
			}),
		},
		"config content-type missing from state": {
			config: newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
				"EMPTY_HEADERS": {
					"*/*":              baseConfigObj,
					"application/json": baseConfigObj,
				},
			}),
			state: newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
				"EMPTY_HEADERS": {"*/*": baseStateObj},
			}),
		},
		"empty config map with non-empty state": {
			config: newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{}),
			state: newResponseTemplatesMap(t, map[string]map[string]basetypes.ObjectValue{
				"EMPTY_HEADERS": {"*/*": baseStateObj},
			}),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := runResponseTemplatePlanModify(tc.config, tc.state, tc.config)
			if !got.Equal(tc.config) {
				t.Errorf("Expected plan to remain config (no suppression); got %s", got)
			}
		})
	}
}

func TestResponseTemplate_Description(t *testing.T) {
	t.Parallel()

	desc := mapplanmodifier.ResponseTemplate().Description(context.Background())
	if desc == "" {
		t.Error("Expected non-empty description")
	}
}

func TestResponseTemplate_MarkdownDescription(t *testing.T) {
	t.Parallel()

	desc := mapplanmodifier.ResponseTemplate().MarkdownDescription(context.Background())
	if desc == "" {
		t.Error("Expected non-empty markdown description")
	}
}
