package acceptance_test

import (
	"encoding/json"

	"github.com/gravitee-io/terraform-provider-apim/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Returns a mapping of provider type names to provider server implementations,
// suitable for acceptance testing via the
// resource.TestCase.ProtoV6ProtocolFactories field.
func testProviders() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"apim": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
}

// importStateIDFunc builds an ImportStateIdFunc that JSON-encodes a composite import ID
// from the resource's state attributes. The keys slice lists the attribute names to include.
// Values are read from state, but can be overridden via the overrides map.
func importStateIDFunc(resourceAddress string, keys []string, overrides map[string]string) func(*terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		attrs := s.RootModule().Resources[resourceAddress].Primary.Attributes
		id := make(map[string]string, len(keys))
		for _, k := range keys {
			if v, ok := overrides[k]; ok {
				id[k] = v
			} else {
				id[k] = attrs[k]
			}
		}
		b, err := json.Marshal(id)
		return string(b), err
	}
}
