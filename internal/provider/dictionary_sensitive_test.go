package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var dictionaryHTTPProvider = path.Root("dynamic").AtName("provider").AtName("http")

var dictionarySensitivePaths = []path.Path{
	dictionaryHTTPProvider.AtName("url"),
	dictionaryHTTPProvider.AtName("headers").AtListIndex(0).AtName("value"),
}

func TestDictionaryResourceHidesHTTPProviderSecrets(t *testing.T) {
	ctx := context.Background()
	resp := &resource.SchemaResponse{}
	NewDictionaryResource().Schema(ctx, resource.SchemaRequest{}, resp)
	require.False(t, resp.Diagnostics.HasError())

	for _, p := range dictionarySensitivePaths {
		attr, diags := resp.Schema.AttributeAtPath(ctx, p)
		require.False(t, diags.HasError(), p.String())
		assert.True(t, attr.IsSensitive(), p.String())
	}
}

func TestDictionaryDataSourceHidesHTTPProviderSecrets(t *testing.T) {
	ctx := context.Background()
	resp := &datasource.SchemaResponse{}
	NewDictionaryDataSource().Schema(ctx, datasource.SchemaRequest{}, resp)
	require.False(t, resp.Diagnostics.HasError())

	for _, p := range dictionarySensitivePaths {
		attr, diags := resp.Schema.AttributeAtPath(ctx, p)
		require.False(t, diags.HasError(), p.String())
		assert.True(t, attr.IsSensitive(), p.String())
	}
}
