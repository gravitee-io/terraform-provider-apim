// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package customtypes

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"strings"
)

var (
	_ basetypes.StringValuable       = (*TrimmedString)(nil)
	_ xattr.ValidateableAttribute    = (*TrimmedString)(nil)
	_ function.ValidateableParameter = (*TrimmedString)(nil)
)

// TrimmedString represents a string that are equal regardless of leading and trailing spaces, tab, and new line.
type TrimmedString struct {
	basetypes.StringValue
}

// Type returns an RFC3339Type.
func (v TrimmedString) Type(_ context.Context) attr.Type {
	return TrimmedStringType{}
}

// Equal returns true if the given value is equivalent.
func (v TrimmedString) Equal(o attr.Value) bool {
	other, ok := o.(TrimmedString)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// StringSemanticEquals
func (v TrimmedString) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(TrimmedString)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	currentValue := v.ValueString()

	if v.Equal(newValue) {
		return true, diags
	}

	return strings.TrimSpace(currentValue) == strings.TrimSpace(newValue.ValueString()), diags
}

// ValidateAttribute does nothing as it accepts any string
func (v TrimmedString) ValidateAttribute(ctx context.Context, req xattr.ValidateAttributeRequest, resp *xattr.ValidateAttributeResponse) {
	// no op
}

// ValidateParameter
func (v TrimmedString) ValidateParameter(ctx context.Context, req function.ValidateParameterRequest, resp *function.ValidateParameterResponse) {
	// no op
}

// NewTrimmedStringNull creates a TrimmedString with a null value. Determine whether the value is null via IsNull method.
func NewTrimmedStringNull() TrimmedString {
	return TrimmedString{
		StringValue: basetypes.NewStringNull(),
	}
}

// NewTrimmedStringUnknown creates a TrimmedString with an unknown value. Determine whether the value is unknown via IsUnknown method.
func NewTrimmedStringUnknown() TrimmedString {
	return TrimmedString{
		StringValue: basetypes.NewStringUnknown(),
	}
}

// NewTrimmedStringValue creates a TrimmedString with a known value.
func NewTrimmedStringValue(value string) TrimmedString {
	return TrimmedString{
		StringValue: basetypes.NewStringValue(value),
	}
}

// NewTrimmedStringPointerValue creates a TrimmedString from a pointer or a null value if the pointer is nil.
func NewTrimmedStringPointerValue(value *string) TrimmedString {
	if value == nil {
		return NewTrimmedStringNull()
	}
	return NewTrimmedStringValue(*value)
}
