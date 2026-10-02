// Package listkeys reads the string attributes that identify the elements of a
// list of objects, for the validators and plan modifiers that compare a declared
// list with the one the platform reports.
package listkeys

import (
	"slices"
	"unicode/utf16"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Of returns, for each element of the list, the values of the named attributes.
// ok is false when the list, an element or one of the values is null or unknown:
// there is nothing to compare yet.
func Of(list basetypes.ListValue, names ...string) (keys [][]string, ok bool) {
	if list.IsNull() || list.IsUnknown() {
		return nil, false
	}
	keys = make([][]string, 0, len(list.Elements()))
	for _, element := range list.Elements() {
		object, isObject := element.(basetypes.ObjectValue)
		if !isObject || object.IsNull() || object.IsUnknown() {
			return nil, false
		}
		key := make([]string, 0, len(names))
		for _, name := range names {
			value, found := lookup(object, name)
			if !found {
				return nil, false
			}
			key = append(key, value)
		}
		keys = append(keys, key)
	}
	return keys, true
}

// lookup reads the attribute on the object, or one level down in a nested object
// that is set: a union holds its attributes in the block of the variant declared.
func lookup(object basetypes.ObjectValue, name string) (string, bool) {
	attributes := object.Attributes()
	if value, declared := attributes[name]; declared {
		return stringOf(value)
	}
	for _, attribute := range attributes {
		nested, isObject := attribute.(basetypes.ObjectValue)
		if !isObject || nested.IsNull() {
			continue
		}
		if nested.IsUnknown() {
			return "", false
		}
		if value, declared := nested.Attributes()[name]; declared {
			return stringOf(value)
		}
	}
	return "", false
}

func stringOf(value any) (string, bool) {
	s, isString := value.(basetypes.StringValue)
	if !isString || s.IsNull() || s.IsUnknown() {
		return "", false
	}
	return s.ValueString(), true
}

// Compare orders two keys attribute by attribute, each by UTF-16 code unit: the
// order of Java's String.compareTo, which the platform sorts with.
func Compare(a, b []string) int {
	return slices.CompareFunc(a, b, func(x, y string) int {
		return slices.Compare(utf16.Encode([]rune(x)), utf16.Encode([]rune(y)))
	})
}

// Sorted returns the keys in the order the platform reports them.
func Sorted(keys [][]string) [][]string {
	sorted := slices.Clone(keys)
	slices.SortStableFunc(sorted, Compare)
	return sorted
}

// IsSorted tells whether the keys are already in that order.
func IsSorted(keys [][]string) bool {
	return slices.IsSortedFunc(keys, Compare)
}
