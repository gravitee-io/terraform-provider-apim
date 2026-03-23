package listplanmodifier_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/gravitee-io/terraform-provider-apim/internal/planmodifiers/listplanmodifier"
)

func TestIgnoreEmptyList_PlanModifyList(t *testing.T) {
	t.Parallel()

	emptyList, _ := basetypes.NewListValue(types.StringType, []attr.Value{})
	nonEmptyList, _ := basetypes.NewListValue(types.StringType, []attr.Value{basetypes.NewStringValue("a")})
	nullList := basetypes.NewListNull(types.StringType)
	unknownList := basetypes.NewListUnknown(types.StringType)

	testCases := map[string]struct {
		configValue   basetypes.ListValue
		stateValue    basetypes.ListValue
		planValue     basetypes.ListValue
		expectedValue basetypes.ListValue
	}{
		"config null + state empty list - plan gets state value": {
			configValue:   nullList,
			stateValue:    emptyList,
			planValue:     nullList,
			expectedValue: emptyList,
		},
		"config null + state non-empty list - plan unchanged": {
			configValue:   nullList,
			stateValue:    nonEmptyList,
			planValue:     nullList,
			expectedValue: nullList,
		},
		"config null + state null - plan unchanged": {
			configValue:   nullList,
			stateValue:    nullList,
			planValue:     nullList,
			expectedValue: nullList,
		},
		"config set - plan unchanged": {
			configValue:   nonEmptyList,
			stateValue:    emptyList,
			planValue:     nonEmptyList,
			expectedValue: nonEmptyList,
		},
		"config null + state unknown - plan gets state value": {
			configValue:   nullList,
			stateValue:    unknownList,
			planValue:     nullList,
			expectedValue: unknownList,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := planmodifier.ListRequest{
				ConfigValue: tc.configValue,
				StateValue:  tc.stateValue,
				PlanValue:   tc.planValue,
			}
			resp := &planmodifier.ListResponse{
				PlanValue: tc.planValue,
			}

			listplanmodifier.IgnoreEmptyList().PlanModifyList(context.Background(), req, resp)

			if !resp.PlanValue.Equal(tc.expectedValue) {
				t.Errorf("Expected plan value %s, got %s", tc.expectedValue, resp.PlanValue)
			}
		})
	}
}

func TestIgnoreEmptyList_Description(t *testing.T) {
	t.Parallel()

	desc := listplanmodifier.IgnoreEmptyList().Description(context.Background())
	if desc == "" {
		t.Error("Expected non-empty description")
	}
}

func TestIgnoreEmptyList_MarkdownDescription(t *testing.T) {
	t.Parallel()

	desc := listplanmodifier.IgnoreEmptyList().MarkdownDescription(context.Background())
	if desc == "" {
		t.Error("Expected non-empty markdown description")
	}
}
