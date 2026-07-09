package provider

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var cloudScopePathPatternOrg = regexp.MustCompile(`/organizations/([^/]+)`)
var cloudScopePathPatternOrgEnv = regexp.MustCompile(`/organizations/([^/]+)/environments/([^/]+)`)

type cloudScopeTransport struct {
	transport http.RoundTripper
	claims    *CloudTokenClaimsData
}

func NewCloudScopeTransport(transport http.RoundTripper, claims *CloudTokenClaimsData) http.RoundTripper {
	if claims == nil {
		return transport
	}
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &cloudScopeTransport{
		transport: transport,
		claims:    claims,
	}
}

func (t *cloudScopeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	orgEnv, ok := scopeFromRequestPath(req.URL.Path)
	if ok {
		if err := validateCloudScope(*t.claims, orgEnv); err != nil {
			defer func() {
				if req != nil && req.Body != nil {
					_ = req.Body.Close()
				}
			}()
			tflog.Debug(req.Context(), "cloud token scope violation", map[string]interface{}{
				"error": err,
			})
			return nil, fmt.Errorf("cloud token scope violation: %w", err)
		}
	}
	return t.transport.RoundTrip(req)
}

func scopeFromRequestPath(path string) (orgEnv, bool) {
	matchesOrgEnv := cloudScopePathPatternOrgEnv.FindStringSubmatch(path)
	if len(matchesOrgEnv) == 3 {
		return orgEnv{Org: matchesOrgEnv[1], Env: matchesOrgEnv[2]}, true
	}
	matchesOrgOnly := cloudScopePathPatternOrg.FindStringSubmatch(path)
	if len(matchesOrgOnly) == 2 {
		return orgEnv{Org: matchesOrgOnly[1]}, true
	}
	return orgEnv{}, false
}
