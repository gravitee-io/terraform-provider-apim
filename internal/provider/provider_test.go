package provider_test

import (
	"github.com/gravitee-io/terraform-provider-apim/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// Returns a mapping of provider type names to provider server implementations,
// suitable for acceptance testing via the
// resource.TestCase.ProtoV6ProtocolFactories field.
func testProviders() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"apim": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
}
