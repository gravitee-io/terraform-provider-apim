package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validateEncryptedValue(t *testing.T, value string) bool {
	t.Helper()
	ctx := context.Background()
	resp := &resource.SchemaResponse{}
	NewDictionaryResource().Schema(ctx, resource.SchemaRequest{}, resp)
	require.False(t, resp.Diagnostics.HasError())

	p := dictionaryEncryptedProperties.AtMapKey("api-key").AtName("value")
	attr, diags := resp.Schema.AttributeAtPath(ctx, p)
	require.False(t, diags.HasError(), p.String())
	stringAttr, ok := attr.(schema.StringAttribute)
	require.True(t, ok, p.String())

	request := validator.StringRequest{Path: p, ConfigValue: types.StringValue(value)}
	valid := true
	for _, v := range stringAttr.Validators {
		response := &validator.StringResponse{}
		v.ValidateString(ctx, request, response)
		valid = valid && !response.Diagnostics.HasError()
	}
	return valid
}

func TestDictionaryResourceRejectsAnEmptyEncryptedValue(t *testing.T) {
	assert.False(t, validateEncryptedValue(t, ""))
	assert.True(t, validateEncryptedValue(t, "s3cr3t"))
}
