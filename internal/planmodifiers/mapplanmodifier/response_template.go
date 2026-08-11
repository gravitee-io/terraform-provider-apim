package mapplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ResponseTemplate returns a plan modifier for `apim_apiv4.response_templates`.
//
// It mirrors IgnoreEmptyList for empty maps, and also suppresses the GKO-3100
// drift where config sets `headers = {}` but Read maps it back to `null`.
func ResponseTemplate() planmodifier.Map {
	return responseTemplate{}
}

type responseTemplate struct{}

// Description returns a human-readable description of the plan modifier.
func (m responseTemplate) Description(_ context.Context) string {
	return "Suppresses empty response_templates.headers map drift"
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m responseTemplate) MarkdownDescription(_ context.Context) string {
	return "Suppresses perpetual plan diffs when `response_templates.*.*.headers` is set to `{}` but APIM returns `null`."
}

// PlanModifyMap implements the plan modification logic.
func (m responseTemplate) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	// nothing set and plan is {} then use state value
	// (same semantics as IgnoreEmptyList).
	if req.ConfigValue.IsNull() && !req.StateValue.IsNull() && len(req.StateValue.Elements()) == 0 {
		resp.PlanValue = req.StateValue
		return
	}

	// Special-case suppression for the APIV4 `response_templates.*.*.headers`
	// drift described in GKO-3100.
	//
	// When config sets `headers = {}` (empty map), APIM can persist it but the
	// provider Read maps it back to `null`. This produces a perpetual plan diff.
	//
	// We can’t attach a plan modifier directly to nested `headers` with the
	// current generated schema, so we apply this modifier on the whole
	// `response_templates` map and detect the specific “only headers differs”
	// pattern.
	if req.ConfigValue.IsNull() || req.StateValue.IsNull() || req.ConfigValue.IsUnknown() || req.StateValue.IsUnknown() {
		return
	}

	if ignoreEmptyResponseTemplatesHeadersDiff(req.ConfigValue, req.StateValue) {
		resp.PlanValue = req.StateValue
	}
}

func ignoreEmptyResponseTemplatesHeadersDiff(configValue, stateValue attr.Value) bool {
	configMap, ok := configValue.(basetypes.MapValue)
	if !ok {
		return false
	}
	stateMap, ok := stateValue.(basetypes.MapValue)
	if !ok {
		return false
	}

	// Outer keys: template code (e.g. INVALID_HTTP_METHOD)
	configOuterElements := configMap.Elements()
	stateOuterElements := stateMap.Elements()

	// Require identical key sets so we never hide additions/removals.
	if len(configOuterElements) != len(stateOuterElements) {
		return false
	}

	for templateCode, configInnerAttr := range configOuterElements {
		stateInnerAttr, ok := stateOuterElements[templateCode]
		if !ok {
			return false
		}

		configInnerMap, ok := configInnerAttr.(basetypes.MapValue) // content-type -> ResponseTemplate object
		if !ok {
			return false
		}
		stateInnerMap, ok := stateInnerAttr.(basetypes.MapValue)
		if !ok {
			return false
		}

		configInnerElements := configInnerMap.Elements()
		stateInnerElements := stateInnerMap.Elements()

		if len(configInnerElements) != len(stateInnerElements) {
			return false
		}

		for contentType, configTemplateObjAttr := range configInnerElements {
			stateTemplateObjAttr, ok := stateInnerElements[contentType]
			if !ok {
				return false
			}

			configTemplateObj, ok := configTemplateObjAttr.(basetypes.ObjectValue)
			if !ok {
				return false
			}
			stateTemplateObj, ok := stateTemplateObjAttr.(basetypes.ObjectValue)
			if !ok {
				return false
			}

			configObjAttrs := configTemplateObj.Attributes()
			stateObjAttrs := stateTemplateObj.Attributes()

			// Validate other fields are equal; otherwise we could hide real diffs.
			if !configObjAttrs["status"].Equal(stateObjAttrs["status"]) {
				return false
			}
			if !configObjAttrs["body"].Equal(stateObjAttrs["body"]) {
				return false
			}
			if !configObjAttrs["propagate_error_key_to_logs"].Equal(stateObjAttrs["propagate_error_key_to_logs"]) {
				return false
			}

			// headers: config has empty map {}, state has null.
			configHeadersAttr := configObjAttrs["headers"]
			stateHeadersAttr := stateObjAttrs["headers"]

			configHeadersMap, ok := configHeadersAttr.(basetypes.MapValue)
			if !ok {
				return false
			}
			stateHeadersMap, ok := stateHeadersAttr.(basetypes.MapValue)
			if !ok {
				return false
			}

			if configHeadersMap.IsNull() {
				return false
			}

			if len(configHeadersMap.Elements()) == 0 && stateHeadersMap.IsNull() {
				continue
			}

			// If headers aren’t the exact drift pattern, don’t suppress.
			return false
		}
	}

	return true
}
