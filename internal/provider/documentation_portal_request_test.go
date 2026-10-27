package provider

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestConfiguredDocumentationPortalRequest(t *testing.T) {
	t.Parallel()
	model := DocumentationPortalResourceModel{
		Hrid: types.StringValue("test-item"), Name: types.StringValue("Docs"),
		Type: types.StringValue("GRAVITEE_MARKDOWN"), Content: types.StringValue("# Docs"),
		PortalHrid: types.StringValue("test-portal"), Visibility: types.StringValue("PUBLIC"),
	}
	original := model
	testConfiguredVisibilityRequest(t, func(ctx context.Context, config tfsdk.Config, plan tfsdk.Plan) (any, diag.Diagnostics) {
		request, diags := model.toConfiguredDocumentationPortalRequest(ctx, config, plan)
		if !reflect.DeepEqual(model, original) {
			t.Fatal("request construction mutated state")
		}
		if request == nil {
			return nil, diags
		}
		if request.PortalHrid != "test-portal" {
			t.Fatal("request lost its parent portal")
		}
		return request.DocumentationPortalSpec, diags
	})
}
