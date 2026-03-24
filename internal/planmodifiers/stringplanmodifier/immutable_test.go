package stringplanmodifier_test

import (
	"context"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/internal/planmodifiers/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestImmutable_Creation(t *testing.T) {
	t.Parallel()

	req := planmodifier.StringRequest{
		Path:       path.Root("test"),
		StateValue: types.StringNull(),
		PlanValue:  types.StringValue("new-value"),
	}
	resp := &planmodifier.StringResponse{}

	stringplanmodifier.Immutable().PlanModifyString(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error on creation, got: %s", resp.Diagnostics.Errors())
	}
}

func TestImmutable_NoChange(t *testing.T) {
	t.Parallel()

	req := planmodifier.StringRequest{
		Path:       path.Root("test"),
		StateValue: types.StringValue("same"),
		PlanValue:  types.StringValue("same"),
	}
	resp := &planmodifier.StringResponse{}

	stringplanmodifier.Immutable().PlanModifyString(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error when value unchanged, got: %s", resp.Diagnostics.Errors())
	}
}

func TestImmutable_Changed(t *testing.T) {
	t.Parallel()

	req := planmodifier.StringRequest{
		Path:       path.Root("test"),
		StateValue: types.StringValue("original"),
		PlanValue:  types.StringValue("changed"),
	}
	resp := &planmodifier.StringResponse{}

	stringplanmodifier.Immutable().PlanModifyString(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error when changing immutable attribute, got none")
	}

	if resp.Diagnostics.ErrorsCount() != 1 {
		t.Fatalf("expected 1 error, got %d", resp.Diagnostics.ErrorsCount())
	}

	got := resp.Diagnostics.Errors()[0].Summary()
	if got != "Immutable attribute" {
		t.Fatalf("expected summary 'Immutable attribute', got '%s'", got)
	}
}

func TestImmutable_UnknownPlanValue(t *testing.T) {
	t.Parallel()

	req := planmodifier.StringRequest{
		Path:       path.Root("test"),
		StateValue: types.StringValue("existing-api"),
		PlanValue:  types.StringUnknown(), // computed, not yet resolved
	}
	resp := &planmodifier.StringResponse{}

	stringplanmodifier.Immutable().PlanModifyString(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error for unknown plan value on existing resource, got: %s", resp.Diagnostics.Errors())
	}
}
