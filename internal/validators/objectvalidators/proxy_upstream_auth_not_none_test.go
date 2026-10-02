package objectvalidators_test

import (
	"context"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/internal/validators/objectvalidators"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	noneType   = map[string]attr.Type{"type": types.StringType}
	bearerType = map[string]attr.Type{"token": types.StringType}
	authType   = map[string]attr.Type{
		"none":   types.ObjectType{AttrTypes: noneType},
		"bearer": types.ObjectType{AttrTypes: bearerType},
	}
	proxyType = map[string]attr.Type{
		"server_url":    types.StringType,
		"upstream_auth": types.ObjectType{AttrTypes: authType},
	}
)

func proxy(upstreamAuth basetypes.ObjectValue) basetypes.ObjectValue {
	return types.ObjectValueMust(proxyType, map[string]attr.Value{
		"server_url":    types.StringValue("https://mcp.example.com/mcp"),
		"upstream_auth": upstreamAuth,
	})
}

func TestProxyUpstreamAuthNotNone(t *testing.T) {
	t.Parallel()

	none := types.ObjectValueMust(authType, map[string]attr.Value{
		"none":   types.ObjectValueMust(noneType, map[string]attr.Value{"type": types.StringValue("NONE")}),
		"bearer": types.ObjectNull(bearerType),
	})
	bearer := types.ObjectValueMust(authType, map[string]attr.Value{
		"none":   types.ObjectNull(noneType),
		"bearer": types.ObjectValueMust(bearerType, map[string]attr.Value{"token": types.StringValue("t")}),
	})

	testCases := map[string]struct {
		config  basetypes.ObjectValue
		refused bool
	}{
		"none declared":             {config: proxy(none), refused: true},
		"another variant declared":  {config: proxy(bearer)},
		"upstream_auth left out":    {config: proxy(types.ObjectNull(authType))},
		"upstream_auth unknown":     {config: proxy(types.ObjectUnknown(authType))},
		"proxy left out (a studio)": {config: types.ObjectNull(proxyType)},
		"proxy unknown":             {config: types.ObjectUnknown(proxyType)},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &validator.ObjectResponse{}
			objectvalidators.ProxyUpstreamAuthNotNone().ValidateObject(context.Background(), validator.ObjectRequest{
				Path:        path.Root("proxy"),
				ConfigValue: tc.config,
			}, resp)

			if resp.Diagnostics.HasError() != tc.refused {
				t.Fatalf("refused = %t, expected %t: %s", resp.Diagnostics.HasError(), tc.refused, resp.Diagnostics.Errors())
			}
		})
	}
}
