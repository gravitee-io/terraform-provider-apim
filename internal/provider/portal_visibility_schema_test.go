package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestPortalDeclaredVisibilitySchema(t *testing.T) {
	t.Parallel()
	var portal resource.SchemaResponse
	NewPortalResource().(*PortalResource).Schema(context.Background(), resource.SchemaRequest{}, &portal)
	structure := portal.Schema.Attributes["structure"].(schema.SingleNestedAttribute)
	visibility := structure.Attributes["top_navbar"].(schema.ListNestedAttribute).NestedObject.Attributes["visibility"]
	if !visibility.IsOptional() || visibility.IsComputed() {
		t.Fatal("declared navigation visibility must remain optional and non-computed")
	}
}

func TestRelatedVisibilitySchemas(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		resource resource.Resource
		nested   string
		computed bool
	}{
		{name: "portal listing", resource: NewPortalListingResource(), nested: "apis"},
		{name: "API navigation", resource: NewApiv4Resource(), nested: "portal_navigation"},
		{name: "portal link", resource: NewPortalLinkResource(), computed: true},
		{name: "API link", resource: NewAPILinkResource(), computed: true},
		{name: "portal documentation", resource: NewDocumentationPortalResource(), computed: true},
		{name: "API documentation", resource: NewDocumentationAPIResource(), computed: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var response resource.SchemaResponse
			test.resource.Schema(context.Background(), resource.SchemaRequest{}, &response)
			attributes := response.Schema.Attributes
			if test.nested != "" {
				attributes = attributes[test.nested].(schema.ListNestedAttribute).NestedObject.Attributes
			}
			visibility := attributes["visibility"]
			if !visibility.IsOptional() || visibility.IsComputed() != test.computed {
				t.Fatalf("visibility should be optional with computed=%t", test.computed)
			}
		})
	}
}
