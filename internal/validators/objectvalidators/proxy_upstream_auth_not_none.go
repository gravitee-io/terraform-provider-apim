package objectvalidators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ validator.Object = proxyUpstreamAuthNotNone{}

// ProxyUpstreamAuthNotNone returns a validator for the proxy block of an MCP
// proxy. The platform reads a proxy declared with no upstream authentication
// back without upstream_auth, so declaring the none variant shows as a change on
// every plan. Leaving upstream_auth out says the same and reads back as declared.
func ProxyUpstreamAuthNotNone() validator.Object {
	return proxyUpstreamAuthNotNone{}
}

type proxyUpstreamAuthNotNone struct{}

func (v proxyUpstreamAuthNotNone) Description(_ context.Context) string {
	return "upstream_auth must be left out when the upstream server needs no authentication"
}

func (v proxyUpstreamAuthNotNone) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v proxyUpstreamAuthNotNone) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if !declared(req.ConfigValue) {
		return
	}
	upstreamAuth, isObject := req.ConfigValue.Attributes()["upstream_auth"].(basetypes.ObjectValue)
	if !isObject || !declared(upstreamAuth) {
		return
	}
	none, isObject := upstreamAuth.Attributes()["none"].(basetypes.ObjectValue)
	if !isObject || !declared(none) {
		return
	}

	resp.Diagnostics.AddAttributeError(
		req.Path.AtName("upstream_auth"),
		"No upstream authentication is declared by leaving upstream_auth out",
		"The platform reports a proxy declared with upstream_auth.none without upstream_auth. "+
			"Remove upstream_auth: the gateway passes the caller's credentials through either way.",
	)
}

func declared(object basetypes.ObjectValue) bool {
	return !object.IsNull() && !object.IsUnknown()
}
