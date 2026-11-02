package provider

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestConfiguredAPILinkRequest(t *testing.T) {
	t.Parallel()
	model := APILinkResourceModel{
		Hrid: types.StringValue("test-item"), Name: types.StringValue("Docs"),
		Href: types.StringValue("https://example.com/docs"), APIHrid: types.StringValue("test-api"),
		Visibility: types.StringValue("PUBLIC"),
	}
	original := model
	testConfiguredVisibilityRequest(t, func(ctx context.Context, config tfsdk.Config, plan tfsdk.Plan) (any, diag.Diagnostics) {
		request, diags := model.toConfiguredAPILinkRequest(ctx, config, plan)
		if !reflect.DeepEqual(model, original) {
			t.Fatal("request construction mutated state")
		}
		if request == nil {
			return nil, diags
		}
		if request.APIHrid != "test-api" {
			t.Fatal("request lost its parent API")
		}
		return request.PortalAttachedLinkSpec, diags
	})
}
