package objectplanmodifier_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/gravitee-io/terraform-provider-apim/internal/planmodifiers/objectplanmodifier"
)

var (
	policyAttributeTypes = map[string]attr.Type{"encrypt_on_fetch": types.BoolType}
	policyType           = tftypes.Object{AttributeTypes: map[string]tftypes.Type{"encrypt_on_fetch": tftypes.Bool}}
	dynamicType          = tftypes.Object{AttributeTypes: map[string]tftypes.Type{"encryption": policyType}}
	rootType             = tftypes.Object{AttributeTypes: map[string]tftypes.Type{"dynamic": dynamicType}}
	encryptionPath       = path.Root("dynamic").AtName("encryption")
)

var dictionarySchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"dynamic": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{
				"encryption": schema.SingleNestedAttribute{
					Optional: true,
					Computed: true,
					Attributes: map[string]schema.Attribute{
						"encrypt_on_fetch": schema.BoolAttribute{Optional: true, Computed: true},
					},
				},
			},
		},
	},
}

func policy(encrypt bool) types.Object {
	return types.ObjectValueMust(policyAttributeTypes, map[string]attr.Value{"encrypt_on_fetch": types.BoolValue(encrypt)})
}

func rawPolicy(encrypt bool) tftypes.Value {
	return tftypes.NewValue(policyType, map[string]tftypes.Value{"encrypt_on_fetch": tftypes.NewValue(tftypes.Bool, encrypt)})
}

func rawState(encryption *tftypes.Value, dynamicSet bool) tftypes.Value {
	if !dynamicSet {
		return tftypes.NewValue(rootType, map[string]tftypes.Value{"dynamic": tftypes.NewValue(dynamicType, nil)})
	}
	policyValue := tftypes.NewValue(policyType, nil)
	if encryption != nil {
		policyValue = *encryption
	}
	return tftypes.NewValue(rootType, map[string]tftypes.Value{
		"dynamic": tftypes.NewValue(dynamicType, map[string]tftypes.Value{"encryption": policyValue}),
	})
}

func TestDisableEncryptionWhenOmitted_PlanModifyObject(t *testing.T) {
	t.Parallel()

	encrypted := rawPolicy(true)
	notEncrypted := rawPolicy(false)
	nullPolicy := types.ObjectNull(policyAttributeTypes)
	unknownPolicy := types.ObjectUnknown(policyAttributeTypes)

	testCases := map[string]struct {
		configValue   types.Object
		stateRaw      tftypes.Value
		stateValue    types.Object
		expectedValue types.Object
	}{
		"policy declared - plan unchanged": {
			configValue:   policy(true),
			stateRaw:      rawState(&notEncrypted, true),
			stateValue:    policy(false),
			expectedValue: unknownPolicy,
		},
		"policy unknown in configuration - plan unchanged": {
			configValue:   unknownPolicy,
			stateRaw:      rawState(&encrypted, true),
			stateValue:    policy(true),
			expectedValue: unknownPolicy,
		},
		"create - plan unchanged, APIM applies its default": {
			configValue:   nullPolicy,
			stateRaw:      tftypes.NewValue(rootType, nil),
			stateValue:    nullPolicy,
			expectedValue: unknownPolicy,
		},
		"dictionary becomes dynamic - plan unchanged, APIM applies its default": {
			configValue:   nullPolicy,
			stateRaw:      rawState(nil, false),
			stateValue:    nullPolicy,
			expectedValue: unknownPolicy,
		},
		"APIM never returned a policy - plan stays null": {
			configValue:   nullPolicy,
			stateRaw:      rawState(nil, true),
			stateValue:    nullPolicy,
			expectedValue: nullPolicy,
		},
		"policy omitted while APIM encrypts - plan turns encryption off": {
			configValue:   nullPolicy,
			stateRaw:      rawState(&encrypted, true),
			stateValue:    policy(true),
			expectedValue: policy(false),
		},
		"policy omitted while APIM does not encrypt - plan keeps it off": {
			configValue:   nullPolicy,
			stateRaw:      rawState(&notEncrypted, true),
			stateValue:    policy(false),
			expectedValue: policy(false),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := planmodifier.ObjectRequest{
				Path:        encryptionPath,
				ConfigValue: tc.configValue,
				StateValue:  tc.stateValue,
				PlanValue:   unknownPolicy,
				State:       tfsdk.State{Schema: dictionarySchema, Raw: tc.stateRaw},
			}
			resp := &planmodifier.ObjectResponse{
				PlanValue: unknownPolicy,
			}

			objectplanmodifier.DisableEncryptionWhenOmitted().PlanModifyObject(context.Background(), req, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("Unexpected diagnostics: %v", resp.Diagnostics)
			}
			if !resp.PlanValue.Equal(tc.expectedValue) {
				t.Errorf("Expected plan value %s, got %s", tc.expectedValue, resp.PlanValue)
			}
		})
	}
}
