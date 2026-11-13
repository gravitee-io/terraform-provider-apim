package stringplanmodifier

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// IgnoreUndeclaredFields returns a plan modifier for a JSON configuration APIM stores with the default of every field
// left out (`{"http":{"readTimeout":5000}}` is read back with every field of `http`).
//
// When the state holds every configured field with the same value, the plan keeps the state and no change is planned.
// Fields only APIM has are not compared, so removing a field from the configuration plans nothing either.
func IgnoreUndeclaredFields() planmodifier.String {
	return ignoreUndeclaredFields{}
}

type ignoreUndeclaredFields struct{}

func (m ignoreUndeclaredFields) Description(_ context.Context) string {
	return "Plans no change when every configured JSON field has the same value in the state, whatever other fields the state holds."
}

func (m ignoreUndeclaredFields) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m ignoreUndeclaredFields) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if !isKnown(req.ConfigValue) || !isKnown(req.PlanValue) || !isKnown(req.StateValue) {
		return
	}
	if !sameListLengths(ctx, req) || !sameIdentity(ctx, req) {
		return
	}

	configured, err := decodeJSON(req.ConfigValue.ValueString())
	if err != nil {
		return
	}
	stored, err := decodeJSON(req.StateValue.ValueString())
	if err != nil {
		return
	}
	if holds(stored, configured) {
		resp.PlanValue = req.StateValue
	}
}

func isKnown(value types.String) bool {
	return !value.IsNull() && !value.IsUnknown()
}

// The framework pairs list items with the state by position. sameListLengths and sameIdentity tell whether the state
// at this path still belongs to the configured item, and not to an item inserted, removed or replaced before it.

// sameListLengths reports whether every list on the path has as many items in the configuration as in the state.
func sameListLengths(ctx context.Context, req planmodifier.StringRequest) bool {
	at := path.Empty()
	for _, step := range req.Path.Steps() {
		switch step := step.(type) {
		case path.PathStepAttributeName:
			at = at.AtName(string(step))
		case path.PathStepElementKeyString:
			at = at.AtMapKey(string(step))
		case path.PathStepElementKeyValue:
			at = at.AtSetValue(step.Value)
		case path.PathStepElementKeyInt:
			configured, configKnown := listLength(ctx, req.Config, at)
			stored, stateKnown := listLength(ctx, req.State, at)
			if !configKnown || !stateKnown || configured != stored {
				return false
			}
			at = at.AtListIndex(int(step))
		default:
			return false
		}
	}
	return true
}

// identityAttributes name the siblings that tell which plugin a configuration belongs to: `policy` for a flow step,
// `type` for an endpoint, an entrypoint, a resource, a service, a plan security or a page source. Other siblings can
// hold defaults in the state and would always differ.
var identityAttributes = []string{"policy", "type"}

// sameIdentity reports whether the siblings naming the plugin are the same in the configuration and in the state.
func sameIdentity(ctx context.Context, req planmodifier.StringRequest) bool {
	parent := req.Path.ParentPath()
	for _, name := range identityAttributes {
		sibling := parent.AtName(name)
		if _, diags := req.Config.Schema.AttributeAtPath(ctx, sibling); diags.HasError() {
			continue
		}
		configured, configFound := valueAt(ctx, req.Config, sibling)
		stored, stateFound := valueAt(ctx, req.State, sibling)
		if !configFound || !stateFound || !configured.Equal(stored) {
			return false
		}
	}
	return true
}

type attributeGetter interface {
	GetAttribute(ctx context.Context, path path.Path, target any) diag.Diagnostics
}

func valueAt(ctx context.Context, data attributeGetter, at path.Path) (attr.Value, bool) {
	var value attr.Value
	if diags := data.GetAttribute(ctx, at, &value); diags.HasError() || value == nil {
		return nil, false
	}
	return value, true
}

func listLength(ctx context.Context, data attributeGetter, at path.Path) (int, bool) {
	value, found := valueAt(ctx, data, at)
	if !found || value.IsNull() || value.IsUnknown() {
		return 0, false
	}
	list, ok := value.(basetypes.ListValuable)
	if !ok {
		return 0, false
	}
	listValue, diags := list.ToListValue(ctx)
	if diags.HasError() {
		return 0, false
	}
	return len(listValue.Elements()), true
}

func decodeJSON(document string) (any, error) {
	decoder := json.NewDecoder(strings.NewReader(document))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	if _, err := decoder.Token(); err != io.EOF {
		return nil, errors.New("unexpected data after the JSON value")
	}
	return value, nil
}

// holds reports whether stored holds every field of configured with the same value. Lists keep their length and
// order: the items are compared one by one.
func holds(stored, configured any) bool {
	switch configured := configured.(type) {
	case map[string]any:
		storedObject, ok := stored.(map[string]any)
		if !ok {
			return false
		}
		for key, configuredValue := range configured {
			storedValue, found := storedObject[key]
			if !found || !holds(storedValue, configuredValue) {
				return false
			}
		}
		return true
	case []any:
		storedList, ok := stored.([]any)
		if !ok || len(storedList) != len(configured) {
			return false
		}
		for i := range configured {
			if !holds(storedList[i], configured[i]) {
				return false
			}
		}
		return true
	case json.Number:
		storedNumber, ok := stored.(json.Number)
		return ok && sameNumber(configured, storedNumber)
	default:
		return configured == stored
	}
}

// sameNumber compares two JSON numbers by value: APIM may return 30000 for a configured 30000.0.
func sameNumber(a, b json.Number) bool {
	if a == b {
		return true
	}
	x, errX := a.Float64()
	y, errY := b.Float64()
	return errX == nil && errY == nil && x == y
}
