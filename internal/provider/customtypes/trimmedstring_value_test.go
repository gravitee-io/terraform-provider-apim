// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package customtypes_test

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/internal/provider/customtypes"
)

func TestTrimmedString_StringSemanticEquals(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		currentTrimmedStringtime customtypes.TrimmedString
		givenTrimmedStringtime   basetypes.StringValuable
		expectedMatch            bool
		expectedDiags            diag.Diagnostics
	}{
		"not equal - different string": {
			currentTrimmedStringtime: customtypes.NewTrimmedStringValue("foo"),
			givenTrimmedStringtime:   customtypes.NewTrimmedStringValue("bar"),
			expectedMatch:            false,
		},
		"not equal - null vs string": {
			currentTrimmedStringtime: customtypes.NewTrimmedStringValue("foo"),
			givenTrimmedStringtime:   customtypes.NewTrimmedStringNull(),
			expectedMatch:            false,
		},
		"not equal - unknown vs string": {
			currentTrimmedStringtime: customtypes.NewTrimmedStringValue("foo"),
			givenTrimmedStringtime:   customtypes.NewTrimmedStringUnknown(),
			expectedMatch:            false,
		},
		"semantically equal - whitespace left": {
			currentTrimmedStringtime: customtypes.NewTrimmedStringValue(" \t  foo\r\n"),
			givenTrimmedStringtime:   customtypes.NewTrimmedStringValue("foo"),
			expectedMatch:            true,
		},
		"semantically equal - whitespace right": {
			currentTrimmedStringtime: customtypes.NewTrimmedStringValue("foo"),
			givenTrimmedStringtime:   customtypes.NewTrimmedStringValue(" \t  foo\r\n"),
			expectedMatch:            true,
		},
		"semantically equal - all both": {
			currentTrimmedStringtime: customtypes.NewTrimmedStringValue(" \t  foo\f\v"),
			givenTrimmedStringtime:   customtypes.NewTrimmedStringValue("\t foo\n\n\n"),
			expectedMatch:            true,
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			match, diags := testCase.currentTrimmedStringtime.StringSemanticEquals(context.Background(), testCase.givenTrimmedStringtime)

			if testCase.expectedMatch != match {
				t.Errorf("Expected StringSemanticEquals to return: %t, but got: %t", testCase.expectedMatch, match)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestTrimmedStringValidateAttribute(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		TrimmedString customtypes.TrimmedString
		expectedDiags diag.Diagnostics
	}{
		"empty-struct": {
			TrimmedString: customtypes.TrimmedString{},
		},
		"null": {
			TrimmedString: customtypes.NewTrimmedStringNull(),
		},
		"unknown": {
			TrimmedString: customtypes.NewTrimmedStringUnknown(),
		},
		"valid TrimmedString": {
			TrimmedString: customtypes.NewTrimmedStringValue("foo"),
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := xattr.ValidateAttributeResponse{}

			testCase.TrimmedString.ValidateAttribute(
				context.Background(),
				xattr.ValidateAttributeRequest{
					Path: path.Root("test"),
				},
				&resp,
			)

			if diff := cmp.Diff(resp.Diagnostics, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestTrimmedStringValidateParameter(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		TrimmedString   customtypes.TrimmedString
		expectedFuncErr *function.FuncError
	}{
		"empty-struct": {
			TrimmedString: customtypes.TrimmedString{},
		},
		"null": {
			TrimmedString: customtypes.NewTrimmedStringNull(),
		},
		"unknown": {
			TrimmedString: customtypes.NewTrimmedStringUnknown(),
		},
		"valid TrimmedString": {
			TrimmedString: customtypes.NewTrimmedStringValue("foo"),
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := function.ValidateParameterResponse{}

			testCase.TrimmedString.ValidateParameter(
				context.Background(),
				function.ValidateParameterRequest{
					Position: int64(0),
				},
				&resp,
			)

			if diff := cmp.Diff(resp.Error, testCase.expectedFuncErr); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}
